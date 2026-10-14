// Command add_superadmin_raw creates or updates the one superadmin user. It reads
// the database settings from .env like the server does.
//
// There is only ever one superadmin:
//   - none yet: it is created from the flags (-email and -password are required)
//   - one: only the flags you pass are changed, the rest is kept
//   - more than one: nothing is changed, fix the data by hand first
//
// The superadmin also gets a user_profiles row (where the UI theme is saved); a
// missing or soft-deleted one is (re)created on every run.
//
// Changing the password logs the superadmin out everywhere (their sessions are ended).
// The password may be given with SUPERADMIN_PASSWORD instead of -password so it
// does not end up in the shell history.
//
//	go run ./scripts/add_superadmin_raw.go -email admin@company.com -password 'S3cret!pass'
//	go run ./scripts/add_superadmin_raw.go -password 'N3w!password'
//	go run ./scripts/add_superadmin_raw.go -username root -first-name Super -last-name Admin
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
)

const (
	defaultUsername  = "superadmin"
	sessionTokenType = "session"
	// Profile defaults, the same as a user created through the API (service.UserService).
	defaultTheme    = dto.ThemeDark
	defaultLanguage = "en"
)

// usernamePattern matches the rule of the user API (dto.CreateUserRequest).
var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,50}$`)

// options are the command line values. An empty value means "not given".
type options struct {
	username  string
	email     string
	password  string
	firstName string
	lastName  string
}

func main() {
	opts := parseFlags()
	if err := opts.validate(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		flag.Usage()
		os.Exit(2)
	}
	if err := run(opts); err != nil {
		slog.Error("superadmin not saved", "error", err)
		os.Exit(1)
	}
}

func parseFlags() options {
	var opts options
	flag.StringVar(&opts.username, "username", "", "login username (default "+defaultUsername+" when creating)")
	flag.StringVar(&opts.email, "email", "", "email (required when creating)")
	flag.StringVar(&opts.password, "password", "", "new password, 8-72 characters (or set SUPERADMIN_PASSWORD; required when creating)")
	flag.StringVar(&opts.firstName, "first-name", "", "first name")
	flag.StringVar(&opts.lastName, "last-name", "", "last name")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Creates the superadmin, or updates the existing one with the flags given.")
		fmt.Fprintln(os.Stderr, "\nUsage: go run ./scripts/add_superadmin_raw.go [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if opts.password == "" {
		opts.password = os.Getenv("SUPERADMIN_PASSWORD")
	}
	opts.username = strings.TrimSpace(opts.username)
	opts.email = strings.ToLower(strings.TrimSpace(opts.email))
	opts.firstName = strings.TrimSpace(opts.firstName)
	opts.lastName = strings.TrimSpace(opts.lastName)
	return opts
}

// validate checks the values that were given; which ones are required depends on
// whether the superadmin exists, and is checked in run.
func (opts options) validate() error {
	if opts == (options{}) {
		return errors.New("nothing to do: pass at least one flag")
	}
	if opts.username != "" && !usernamePattern.MatchString(opts.username) {
		return errors.New("username must be 3-50 letters, digits, _ or .")
	}
	if opts.email != "" {
		if address, err := mail.ParseAddress(opts.email); err != nil || address.Address != opts.email {
			return errors.New("email is not valid")
		}
	}
	// bcrypt only uses the first 72 bytes of a password.
	if opts.password != "" && (len(opts.password) < 8 || len(opts.password) > 72) {
		return errors.New("password must be 8-72 characters")
	}
	if len(opts.firstName) > 100 || len(opts.lastName) > 100 {
		return errors.New("first and last name must be at most 100 characters")
	}
	return nil
}

func run(opts options) error {
	ctx := context.Background()
	appConfig := config.Load()

	dbPool, err := db.Connect(ctx, appConfig.DB)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	hashedPassword := ""
	if opts.password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(opts.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		hashedPassword = string(hash)
	}

	return db.WithTransaction(ctx, dbPool, func(tx pgx.Tx) error {
		// Two runs at the same time must not both create a superadmin.
		if err := execute(ctx, tx, `SELECT pg_advisory_xact_lock(hashtext('add_superadmin'))`); err != nil {
			return fmt.Errorf("lock: %w", err)
		}

		roleID, err := superadminRoleID(ctx, tx)
		if err != nil {
			return err
		}

		rows, err := db.Execute(ctx, tx, `SELECT id FROM users
			WHERE role_id = $1 AND is_deleted = false
			ORDER BY id
			FOR UPDATE`, roleID)
		if err != nil {
			return fmt.Errorf("find superadmin: %w", err)
		}
		superadminIDs, err := pgx.CollectRows(rows, pgx.RowTo[int64])
		if err != nil {
			return fmt.Errorf("find superadmin: %w", err)
		}

		switch len(superadminIDs) {
		case 0:
			return createSuperadmin(ctx, tx, roleID, opts, hashedPassword)
		case 1:
			return updateSuperadmin(ctx, tx, superadminIDs[0], opts, hashedPassword)
		default:
			return fmt.Errorf("found %d superadmins (user ids %v), there must be only one: soft-delete or change the role of the extra ones first",
				len(superadminIDs), superadminIDs)
		}
	})
}

// superadminRoleID returns the id of the superadmin role, adding the role if the
// database does not have it. Role ids differ between databases, so it is looked
// up by name.
func superadminRoleID(ctx context.Context, tx pgx.Tx) (int64, error) {
	rows, err := db.Execute(ctx, tx, `INSERT INTO roles (name, description)
		VALUES ($1, 'Super administrator with full access')
		ON CONFLICT (name) DO NOTHING
		RETURNING id`, model.RoleSuperAdmin)
	if err != nil {
		return 0, fmt.Errorf("add superadmin role: %w", err)
	}
	roleID, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err == nil {
		slog.Info("superadmin role added", "roleId", roleID)
		return roleID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("add superadmin role: %w", err)
	}

	rows, err = db.Execute(ctx, tx, `SELECT id FROM roles WHERE name = $1`, model.RoleSuperAdmin)
	if err != nil {
		return 0, fmt.Errorf("find superadmin role: %w", err)
	}
	roleID, err = pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return 0, fmt.Errorf("find superadmin role: %w", err)
	}
	return roleID, nil
}

func createSuperadmin(ctx context.Context, tx pgx.Tx, roleID int64, opts options, hashedPassword string) error {
	if opts.email == "" || hashedPassword == "" {
		return errors.New("no superadmin yet: -email and -password are required to create one")
	}
	if opts.username == "" {
		opts.username = defaultUsername
	}
	if err := checkNotTaken(ctx, tx, 0, opts.username, opts.email); err != nil {
		return err
	}

	rows, err := db.Execute(ctx, tx, `INSERT INTO users (username, email, password, first_name, last_name, role_id, status)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6, 'active')
		RETURNING id`,
		opts.username, opts.email, hashedPassword, opts.firstName, opts.lastName, roleID,
	)
	if err != nil {
		return fmt.Errorf("create superadmin: %w", err)
	}
	userID, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return fmt.Errorf("create superadmin: %w", err)
	}
	if err := ensureProfile(ctx, tx, userID); err != nil {
		return err
	}

	slog.Info("superadmin created", "id", userID, "username", opts.username, "email", opts.email)
	return nil
}

// updateSuperadmin changes only the given fields of user id. The superadmin is
// also made active again, so the script can be used to get back in.
func updateSuperadmin(ctx context.Context, tx pgx.Tx, id int64, opts options, hashedPassword string) error {
	if err := checkNotTaken(ctx, tx, id, opts.username, opts.email); err != nil {
		return err
	}

	rows, err := db.Execute(ctx, tx, `UPDATE users
		SET username = COALESCE(NULLIF($2, ''), username),
			email = COALESCE(NULLIF($3, ''), email),
			password = COALESCE(NULLIF($4, ''), password),
			first_name = COALESCE(NULLIF($5, ''), first_name),
			last_name = COALESCE(NULLIF($6, ''), last_name),
			status = 'active'
		WHERE id = $1
		RETURNING username, email`,
		id, opts.username, opts.email, hashedPassword, opts.firstName, opts.lastName,
	)
	if err != nil {
		return fmt.Errorf("update superadmin: %w", err)
	}
	type savedUser struct {
		Username string `db:"username"`
		Email    string `db:"email"`
	}
	saved, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[savedUser])
	if err != nil {
		return fmt.Errorf("update superadmin: %w", err)
	}

	// A superadmin added before profiles existed has none, so give them one.
	if err := ensureProfile(ctx, tx, id); err != nil {
		return err
	}

	if hashedPassword != "" {
		// A new password ends every session, so an old login cannot be reused.
		err := execute(ctx, tx, `UPDATE auth_tokens
			SET is_deleted = true, deleted_at = NOW()
			WHERE user_id = $1 AND token_type = $2 AND is_deleted = false`, id, sessionTokenType)
		if err != nil {
			return fmt.Errorf("end superadmin sessions: %w", err)
		}
	}

	slog.Info("superadmin updated", "id", id, "username", saved.Username, "email", saved.Email,
		"passwordChanged", hashedPassword != "")
	return nil
}

// ensureProfile gives userID an active user_profiles row: it is added with the
// default theme and language, or restored if it was soft-deleted. An active
// profile is left as it is, so saved preferences are kept.
func ensureProfile(ctx context.Context, tx pgx.Tx, userID int64) error {
	err := execute(ctx, tx, `INSERT INTO user_profiles (user_id, theme, language)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
			SET is_deleted = false, deleted_at = NULL
			WHERE user_profiles.is_deleted`, userID, defaultTheme, defaultLanguage)
	if err != nil {
		return fmt.Errorf("add superadmin profile: %w", err)
	}
	return nil
}

// checkNotTaken fails when another active user (not exceptID) already has the
// username or email. Empty values are not checked.
func checkNotTaken(ctx context.Context, tx pgx.Tx, exceptID int64, username, email string) error {
	rows, err := db.Execute(ctx, tx, `SELECT id, username, email FROM users
		WHERE id <> $1 AND is_deleted = false
			AND ((NULLIF($2, '') IS NOT NULL AND lower(username) = lower($2))
				OR (NULLIF($3, '') IS NOT NULL AND lower(email) = lower($3)))
		LIMIT 1`, exceptID, username, email)
	if err != nil {
		return fmt.Errorf("check username and email: %w", err)
	}
	type otherUser struct {
		ID       int64  `db:"id"`
		Username string `db:"username"`
		Email    string `db:"email"`
	}
	other, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[otherUser])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check username and email: %w", err)
	}
	if username != "" && strings.EqualFold(other.Username, username) {
		return fmt.Errorf("username %q is already used by user %d", username, other.ID)
	}
	return fmt.Errorf("email %q is already used by user %d", email, other.ID)
}

// execute runs a query whose rows are not needed.
func execute(ctx context.Context, tx pgx.Tx, query string, args ...any) error {
	rows, err := db.Execute(ctx, tx, query, args...)
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}
