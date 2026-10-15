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

// shiftSelect reads model.Shift: the shifts row (s) plus its customer's display name.
// TIME columns are read as "HH:MM".
const shiftSelect = `SELECT s.id, s.user_id, s.name, s.shift_type,
		to_char(s.start_time, 'HH24:MI') AS start_time, to_char(s.end_time, 'HH24:MI') AS end_time,
		s.working_days, s.status, s.is_deleted, s.deleted_at, s.created_at, s.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name
	FROM shifts s
	LEFT JOIN users u ON u.id = s.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = s.user_id AND cp.is_deleted = false`

// ShiftRepository runs the SQL queries for the shifts table (a shift's breaks are rows
// of breaks with its shift_id, see BreakRepository). Soft-deleted shifts
// (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the shifts of that customer
// (shifts.user_id) are found.
type ShiftRepository struct {
	database db.Querier
}

// NewShiftRepository creates a ShiftRepository that runs its queries on database
// (a connection pool or a transaction).
func NewShiftRepository(database db.Querier) *ShiftRepository {
	return &ShiftRepository{database: database}
}

// FindAll returns the shifts of ownerID (every shift when nil), by customer and start time.
func (shiftRepository *ShiftRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.Shift, error) {
	query := shiftSelect + `
		WHERE s.is_deleted = false AND ($1::integer IS NULL OR s.user_id = $1)
		ORDER BY s.user_id, s.start_time, s.id`

	rows, err := db.Execute(ctx, shiftRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all shifts: %w", err)
	}
	shifts, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Shift])
	if err != nil {
		return nil, fmt.Errorf("find all shifts: %w", err)
	}
	return shifts, nil
}

// FindByID returns the shift, or ErrNotFound when it is not a shift of ownerID.
func (shiftRepository *ShiftRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.Shift, error) {
	query := shiftSelect + `
		WHERE s.id = $1 AND s.is_deleted = false AND ($2::integer IS NULL OR s.user_id = $2)`

	rows, err := db.Execute(ctx, shiftRepository.database, query, id, ownerID)
	if err != nil {
		return model.Shift{}, fmt.Errorf("find shift: %w", err)
	}
	shift, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Shift])
	return shift, mapShiftError("find shift", err)
}

// Create inserts shift and returns the new id.
func (shiftRepository *ShiftRepository) Create(ctx context.Context, shift model.Shift) (int64, error) {
	query := `INSERT INTO shifts (user_id, name, shift_type, start_time, end_time, working_days, status)
		VALUES ($1, $2, $3, $4::time, $5::time, $6::integer[], $7)
		RETURNING id`

	rows, err := db.Execute(ctx, shiftRepository.database, query,
		shift.UserID, shift.Name, shift.ShiftType, shift.StartTime, shift.EndTime, shift.WorkingDays, shift.Status)
	if err != nil {
		return 0, fmt.Errorf("create shift: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapShiftError("create shift", err)
}

// Update replaces the fields of shift.ID. A blank Status keeps the saved status.
func (shiftRepository *ShiftRepository) Update(ctx context.Context, shift model.Shift) error {
	query := `UPDATE shifts
		SET user_id = $2, name = $3, shift_type = $4, start_time = $5::time, end_time = $6::time,
			working_days = $7::integer[], status = COALESCE(NULLIF($8, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, shiftRepository.database, query,
		shift.ID, shift.UserID, shift.Name, shift.ShiftType, shift.StartTime, shift.EndTime, shift.WorkingDays, shift.Status)
	if err != nil {
		return mapShiftError("update shift", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the shift of ownerID (any shift when nil). The service releases
// its breaks (BreakRepository.ReleaseShift).
func (shiftRepository *ShiftRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE shifts
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, shiftRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete shift: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapShiftError turns driver errors into repository errors. A nil err stays nil.
func mapShiftError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "shifts_user_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "shifts_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
