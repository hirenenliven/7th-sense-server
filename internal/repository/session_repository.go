package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// sessionColumns are the columns of a model.Session. The timestamps are timestamp without time
// zone written with NOW(), so they are read back in the database's time zone.
const sessionColumns = `id, user_id, token,
	COALESCE(device_info, '') AS device_info,
	COALESCE(ip_address, '') AS ip_address,
	created_at::timestamptz AS created_at,
	COALESCE(last_used_at, created_at)::timestamptz AS last_used_at,
	expires_at::timestamptz AS expires_at,
	impersonated_by`

// activeSession is the condition of a session that still authenticates.
const activeSession = `token_type = 'session' AND status = 'active' AND is_deleted = false AND expires_at > NOW()`

// SessionRepository runs the SQL queries of session management on auth_tokens: listing a user's
// active sessions, revoking them and cleaning up ended ones. Ownership is always part of the
// WHERE clause. Creating and resolving a session stay in AuthRepository.
type SessionRepository struct {
	database db.Querier
}

// NewSessionRepository creates a SessionRepository that runs its queries on database (a
// connection pool or a transaction).
func NewSessionRepository(database db.Querier) *SessionRepository {
	return &SessionRepository{database: database}
}

// LockUser locks the users row of userID until the transaction ends, so two logins of the same
// user run one after the other and cannot both pass the session limit. It must run in a transaction.
func (sessionRepository *SessionRepository) LockUser(ctx context.Context, userID int64) error {
	rows, err := db.Execute(ctx, sessionRepository.database, `SELECT id FROM users WHERE id = $1 FOR UPDATE`, userID)
	if err != nil {
		return fmt.Errorf("lock user: %w", err)
	}
	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return mapAuthError("lock user", err)
}

// FindAllActive returns userID's active sessions, "log in as" ones included, most recently used first.
func (sessionRepository *SessionRepository) FindAllActive(ctx context.Context, userID int64) ([]model.Session, error) {
	query := `SELECT ` + sessionColumns + `
		FROM auth_tokens
		WHERE user_id = $1 AND ` + activeSession + `
		ORDER BY COALESCE(last_used_at, created_at) DESC, id DESC`

	rows, err := db.Execute(ctx, sessionRepository.database, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find active sessions: %w", err)
	}
	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Session])
	if err != nil {
		return nil, fmt.Errorf("find active sessions: %w", err)
	}
	return sessions, nil
}

// RevokeByIDs revokes the active sessions among ids for reason, by revokedBy, and returns the
// token hashes of the sessions it revoked.
func (sessionRepository *SessionRepository) RevokeByIDs(ctx context.Context, ids []int64, revokedBy int64, reason string) ([]string, error) {
	return sessionRepository.revoke(ctx, `id = ANY($3)`, revokedBy, reason, ids)
}

// RevokeOwn revokes session id when it is an active password-login session of userID and returns
// its token hash; ErrNotFound for any other session ("log in as" sessions are not the user's to end).
func (sessionRepository *SessionRepository) RevokeOwn(ctx context.Context, id, userID int64, reason string) (string, error) {
	hashes, err := sessionRepository.revoke(ctx, `id = $3 AND user_id = $4 AND impersonated_by IS NULL`, userID, reason, id, userID)
	if err != nil {
		return "", err
	}
	if len(hashes) == 0 {
		return "", ErrNotFound
	}
	return hashes[0], nil
}

// RevokeOthers revokes userID's active password-login sessions except keepID and returns their
// token hashes.
func (sessionRepository *SessionRepository) RevokeOthers(ctx context.Context, userID, keepID int64, reason string) ([]string, error) {
	return sessionRepository.revoke(ctx, `user_id = $3 AND id <> $4 AND impersonated_by IS NULL`, userID, reason, userID, keepID)
}

// RevokeByID revokes active session id of any user (superadmin) and returns its token hash;
// ErrNotFound when it is not an active session.
func (sessionRepository *SessionRepository) RevokeByID(ctx context.Context, id, revokedBy int64, reason string) (string, error) {
	hashes, err := sessionRepository.revoke(ctx, `id = $3`, revokedBy, reason, id)
	if err != nil {
		return "", err
	}
	if len(hashes) == 0 {
		return "", ErrNotFound
	}
	return hashes[0], nil
}

// RevokeAllOfUser revokes every active session of userID ("log in as" ones included) and returns
// their token hashes.
func (sessionRepository *SessionRepository) RevokeAllOfUser(ctx context.Context, userID, revokedBy int64, reason string) ([]string, error) {
	return sessionRepository.revoke(ctx, `user_id = $3`, revokedBy, reason, userID)
}

// revoke marks the active sessions matching condition revoked (and soft-deleted, so queries on
// is_deleted keep working) and returns their token hashes. condition uses $3 onwards for args.
func (sessionRepository *SessionRepository) revoke(ctx context.Context, condition string, revokedBy int64, reason string, args ...any) ([]string, error) {
	query := `UPDATE auth_tokens
		SET status = 'revoked', is_deleted = true, deleted_at = NOW(),
			revoked_at = NOW(), revoked_by = $1, revoke_reason = $2
		WHERE ` + activeSession + ` AND ` + condition + `
		RETURNING token`

	rows, err := db.Execute(ctx, sessionRepository.database, query, append([]any{revokedBy, reason}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("revoke sessions: %w", err)
	}
	hashes, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("revoke sessions: %w", err)
	}
	return hashes, nil
}

// ExpireOverdue marks active sessions past their lifetime expired and returns how many.
func (sessionRepository *SessionRepository) ExpireOverdue(ctx context.Context) (int64, error) {
	query := `UPDATE auth_tokens SET status = $1
		WHERE token_type = $2 AND status = $3 AND expires_at <= NOW()`

	rows, err := db.Execute(ctx, sessionRepository.database, query,
		model.SessionStatusExpired, sessionTokenType, model.SessionStatusActive)
	if err != nil {
		return 0, fmt.Errorf("expire sessions: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("expire sessions: %w", err)
	}
	return rows.CommandTag().RowsAffected(), nil
}

// PurgeEnded deletes sessions that ended (logged out, revoked or expired) more than retention ago
// and returns how many. Active sessions are never deleted.
func (sessionRepository *SessionRepository) PurgeEnded(ctx context.Context, retention time.Duration) (int64, error) {
	query := `DELETE FROM auth_tokens
		WHERE token_type = $1 AND status <> $2
			AND COALESCE(revoked_at, deleted_at, expires_at::timestamptz) < NOW() - $3 * INTERVAL '1 second'`

	rows, err := db.Execute(ctx, sessionRepository.database, query,
		sessionTokenType, model.SessionStatusActive, int64(retention.Seconds()))
	if err != nil {
		return 0, fmt.Errorf("purge sessions: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("purge sessions: %w", err)
	}
	return rows.CommandTag().RowsAffected(), nil
}
