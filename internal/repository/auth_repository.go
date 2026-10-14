package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// sessionTokenType is auth_tokens.token_type of a login session.
const sessionTokenType = "session"

// AuthRepository runs the SQL queries for login: reading users with their role
// and storing sessions in auth_tokens. Soft-deleted users and sessions are
// invisible to every method. Tokens are stored as hashes; this repository
// never sees a raw session token.
type AuthRepository struct {
	database db.Querier
}

// NewAuthRepository creates an AuthRepository that runs its queries on database
// (a connection pool or a transaction).
func NewAuthRepository(database db.Querier) *AuthRepository {
	return &AuthRepository{database: database}
}

// FindUserByUsername returns the user with username, including the password
// hash and role name.
func (authRepository *AuthRepository) FindUserByUsername(ctx context.Context, username string) (model.AuthUser, error) {
	query := `SELECT users.id, users.username, users.email, users.password,
			COALESCE(users.first_name, '') AS first_name,
			COALESCE(users.last_name, '') AS last_name,
			COALESCE(users.status, '') AS status,
			users.company_id,
			roles.name AS role
		FROM users
		JOIN roles ON roles.id = users.role_id
		WHERE users.username = $1 AND users.is_deleted = false`

	rows, err := db.Execute(ctx, authRepository.database, query, username)
	if err != nil {
		return model.AuthUser{}, fmt.Errorf("find login user: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.AuthUser])
	return user, mapAuthError("find login user", err)
}

// CreateSession stores a session of userID that expires after lifetime.
func (authRepository *AuthRepository) CreateSession(ctx context.Context, userID int64, tokenHash, deviceInfo string, lifetime time.Duration) error {
	query := `INSERT INTO auth_tokens (user_id, token, token_type, device_type, device_info, expires_at)
		VALUES ($1, $2, $3, 'web', NULLIF($4, ''), NOW() + $5 * INTERVAL '1 second')`

	rows, err := db.Execute(ctx, authRepository.database, query,
		userID, tokenHash, sessionTokenType, deviceInfo, int64(lifetime.Seconds()),
	)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// FindUserBySession returns the user of a live session and marks the session
// as used. A session is live when it is not deleted or expired and its user is
// active and not deleted, so deactivating a user ends their sessions at once.
func (authRepository *AuthRepository) FindUserBySession(ctx context.Context, tokenHash string) (model.AuthUser, error) {
	query := `UPDATE auth_tokens
		SET last_used_at = NOW()
		FROM users
		JOIN roles ON roles.id = users.role_id
		WHERE auth_tokens.token = $1
			AND auth_tokens.token_type = $2
			AND auth_tokens.is_deleted = false
			AND auth_tokens.expires_at > NOW()
			AND users.id = auth_tokens.user_id
			AND users.is_deleted = false
			AND COALESCE(users.status, 'active') = 'active'
		RETURNING users.id, users.username, users.email,
			COALESCE(users.first_name, '') AS first_name,
			COALESCE(users.last_name, '') AS last_name,
			COALESCE(users.status, '') AS status,
			users.company_id,
			roles.name AS role`

	rows, err := db.Execute(ctx, authRepository.database, query, tokenHash, sessionTokenType)
	if err != nil {
		return model.AuthUser{}, fmt.Errorf("find session user: %w", err)
	}
	// Lax: the password hash is not read for a session.
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[model.AuthUser])
	return user, mapAuthError("find session user", err)
}

// DeleteSession soft-deletes the session. Deleting a missing or already deleted
// session succeeds, so logout is idempotent.
func (authRepository *AuthRepository) DeleteSession(ctx context.Context, tokenHash string) error {
	query := `UPDATE auth_tokens
		SET is_deleted = true, deleted_at = NOW()
		WHERE token = $1 AND token_type = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, authRepository.database, query, tokenHash, sessionTokenType)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// mapAuthError turns driver errors into repository errors. A nil err stays nil.
func mapAuthError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
