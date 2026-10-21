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

// downtimeReasonSelect reads model.DowntimeReason: the downtime_reasons row (r) plus its customer's display name.
const downtimeReasonSelect = `SELECT r.id, r.user_id, COALESCE(r.reason_code, '') AS reason_code, r.reason, r.status,
		r.is_deleted, r.deleted_at, r.created_at, r.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name
	FROM downtime_reasons r
	LEFT JOIN users u ON u.id = r.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = r.user_id AND cp.is_deleted = false`

// DowntimeReasonRepository runs the SQL queries for the downtime_reasons table. Soft-deleted
// reasons (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the reasons of that customer
// (downtime_reasons.user_id) are found.
type DowntimeReasonRepository struct {
	database db.Querier
}

// NewDowntimeReasonRepository creates a DowntimeReasonRepository that runs its queries on
// database (a connection pool or a transaction).
func NewDowntimeReasonRepository(database db.Querier) *DowntimeReasonRepository {
	return &DowntimeReasonRepository{database: database}
}

// FindAll returns the downtime reasons of ownerID (every reason when nil), by reason.
func (downtimeReasonRepository *DowntimeReasonRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.DowntimeReason, error) {
	query := downtimeReasonSelect + `
		WHERE r.is_deleted = false AND ($1::integer IS NULL OR r.user_id = $1)
		ORDER BY LOWER(r.reason), r.id`

	rows, err := db.Execute(ctx, downtimeReasonRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all downtime reasons: %w", err)
	}
	reasons, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.DowntimeReason])
	if err != nil {
		return nil, fmt.Errorf("find all downtime reasons: %w", err)
	}
	return reasons, nil
}

// FindByID returns the downtime reason, or ErrNotFound when it is not a reason of ownerID.
func (downtimeReasonRepository *DowntimeReasonRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.DowntimeReason, error) {
	query := downtimeReasonSelect + `
		WHERE r.id = $1 AND r.is_deleted = false AND ($2::integer IS NULL OR r.user_id = $2)`

	rows, err := db.Execute(ctx, downtimeReasonRepository.database, query, id, ownerID)
	if err != nil {
		return model.DowntimeReason{}, fmt.Errorf("find downtime reason: %w", err)
	}
	reason, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.DowntimeReason])
	return reason, mapDowntimeReasonError("find downtime reason", err)
}

// Create inserts reason and returns the new id.
func (downtimeReasonRepository *DowntimeReasonRepository) Create(ctx context.Context, reason model.DowntimeReason) (int64, error) {
	query := `INSERT INTO downtime_reasons (user_id, reason_code, reason, status)
		VALUES ($1, NULLIF($2, ''), $3, $4)
		RETURNING id`

	rows, err := db.Execute(ctx, downtimeReasonRepository.database, query,
		reason.UserID, reason.ReasonCode, reason.Reason, reason.Status)
	if err != nil {
		return 0, fmt.Errorf("create downtime reason: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapDowntimeReasonError("create downtime reason", err)
}

// Update replaces the fields of reason.ID. A blank Status keeps the saved status.
func (downtimeReasonRepository *DowntimeReasonRepository) Update(ctx context.Context, reason model.DowntimeReason) error {
	query := `UPDATE downtime_reasons
		SET user_id = $2, reason_code = NULLIF($3, ''), reason = $4,
			status = COALESCE(NULLIF($5, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, downtimeReasonRepository.database, query,
		reason.ID, reason.UserID, reason.ReasonCode, reason.Reason, reason.Status)
	if err != nil {
		return mapDowntimeReasonError("update downtime reason", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the downtime reason of ownerID (any reason when nil).
func (downtimeReasonRepository *DowntimeReasonRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE downtime_reasons
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, downtimeReasonRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete downtime reason: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapDowntimeReasonError turns driver errors into repository errors. A nil err stays nil.
func mapDowntimeReasonError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "downtime_reasons_user_reason_code_key": // unique_violation
			return ErrDuplicateCode
		case pgError.Code == "23505" && pgError.ConstraintName == "downtime_reasons_user_reason_key":
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "downtime_reasons_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
