package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// UserRepository runs the SQL queries for the users table. Soft-deleted
// users (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.User, so nullable
// columns are read with COALESCE and given their column name with AS.
type UserRepository struct {
	database db.Querier
}

// NewUserRepository creates a UserRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserRepository(database db.Querier) *UserRepository {
	return &UserRepository{database: database}
}

// FindAll returns the users. A non-empty role keeps only users with that role
// name; a non-nil companyID only the members of that customer.
func (userRepository *UserRepository) FindAll(ctx context.Context, role string, companyID *int64) ([]model.User, error) {
	query := `SELECT id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, COALESCE((SELECT name FROM roles WHERE roles.id = users.role_id), '') AS role,
			branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at
		FROM users
		WHERE is_deleted = false
			AND (NULLIF($1, '') IS NULL OR role_id = (SELECT id FROM roles WHERE name = $1))
			AND ($2::integer IS NULL OR company_id = $2)
		ORDER BY id`

	rows, err := db.Execute(ctx, userRepository.database, query, role, companyID)
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	return users, nil
}

func (userRepository *UserRepository) FindByID(ctx context.Context, id int64) (model.User, error) {
	query := `SELECT id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, COALESCE((SELECT name FROM roles WHERE roles.id = users.role_id), '') AS role,
			branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userRepository.database, query, id)
	if err != nil {
		return model.User{}, fmt.Errorf("find user: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return user, mapUserError("find user", err)
}

func (userRepository *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (username, email, password, first_name, last_name, mobile, role_id, branch_id, company_id)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), $7, $8, $9)
		RETURNING id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, COALESCE((SELECT name FROM roles WHERE roles.id = users.role_id), '') AS role,
			branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at`

	rows, err := db.Execute(ctx, userRepository.database, query,
		user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Mobile, user.RoleID, user.BranchID, user.CompanyID,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return created, mapUserError("create user", err)
}

// Update replaces the profile fields of user.ID. An empty Status keeps the
// current status. The password is not changed.
func (userRepository *UserRepository) Update(ctx context.Context, user model.User) (model.User, error) {
	query := `UPDATE users
		SET username = $2, email = $3,
			first_name = NULLIF($4, ''), last_name = NULLIF($5, ''), mobile = NULLIF($6, ''),
			role_id = $7, branch_id = $8, company_id = $9,
			status = COALESCE(NULLIF($10, ''), status)
		WHERE id = $1 AND is_deleted = false
		RETURNING id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, COALESCE((SELECT name FROM roles WHERE roles.id = users.role_id), '') AS role,
			branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at`

	rows, err := db.Execute(ctx, userRepository.database, query,
		user.ID, user.Username, user.Email, user.FirstName, user.LastName, user.Mobile, user.RoleID, user.BranchID, user.CompanyID, user.Status,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return updated, mapUserError("update user", err)
}

// Delete soft-deletes the user: the row is kept with is_deleted = true.
func (userRepository *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapUserError turns driver errors into repository errors. A nil err stays nil.
func mapUserError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "users_email_key": // unique_violation
			return ErrDuplicateEmail
		case pgError.Code == "23505" && pgError.ConstraintName == "users_username_key":
			return ErrDuplicateUsername
		case pgError.Code == "23503" && pgError.ConstraintName == "users_role_id_fkey": // foreign_key_violation
			return ErrInvalidRole
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// IsCustomer reports whether id is an existing user with the customer role.
func (userRepository *UserRepository) IsCustomer(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM users u
			JOIN roles r ON r.id = u.role_id
			WHERE u.id = $1 AND u.is_deleted = false AND r.name = $2
		)`

	rows, err := db.Execute(ctx, userRepository.database, query, id, model.RoleCustomer)
	if err != nil {
		return false, fmt.Errorf("check customer: %w", err)
	}
	isCustomer, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check customer: %w", err)
	}
	return isCustomer, nil
}
