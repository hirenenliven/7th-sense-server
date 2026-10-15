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

// breakSelect reads model.Break: the breaks row (b) plus the display names of its customer
// and shift. TIME columns are read as "HH:MM".
const breakSelect = `SELECT b.id, b.user_id, b.shift_id, b.name,
		to_char(b.start_time, 'HH24:MI') AS start_time, to_char(b.end_time, 'HH24:MI') AS end_time,
		b.status, b.is_deleted, b.deleted_at, b.created_at, b.updated_at,
		COALESCE(NULLIF(cp.company_name, ''), NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(s.name, '') AS shift_name
	FROM breaks b
	LEFT JOIN users u ON u.id = b.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = b.user_id AND cp.is_deleted = false
	LEFT JOIN shifts s ON s.id = b.shift_id`

// BreakRepository runs the SQL queries for the breaks table. A break is taken by the
// shift breaks.shift_id (NULL while no shift takes it). Soft-deleted breaks
// (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the breaks of that customer
// (breaks.user_id) are found.
type BreakRepository struct {
	database db.Querier
}

// NewBreakRepository creates a BreakRepository that runs its queries on database
// (a connection pool or a transaction).
func NewBreakRepository(database db.Querier) *BreakRepository {
	return &BreakRepository{database: database}
}

// FindAll returns the breaks of ownerID (every break when nil), by customer and start time.
func (breakRepository *BreakRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.Break, error) {
	query := breakSelect + `
		WHERE b.is_deleted = false AND ($1::integer IS NULL OR b.user_id = $1)
		ORDER BY b.user_id, b.start_time, b.id`

	rows, err := db.Execute(ctx, breakRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all breaks: %w", err)
	}
	breaks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Break])
	if err != nil {
		return nil, fmt.Errorf("find all breaks: %w", err)
	}
	return breaks, nil
}

// FindByID returns the break, or ErrNotFound when it is not a break of ownerID.
func (breakRepository *BreakRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.Break, error) {
	query := breakSelect + `
		WHERE b.id = $1 AND b.is_deleted = false AND ($2::integer IS NULL OR b.user_id = $2)`

	rows, err := db.Execute(ctx, breakRepository.database, query, id, ownerID)
	if err != nil {
		return model.Break{}, fmt.Errorf("find break: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Break])
	return found, mapBreakError("find break", err)
}

// FindByIDs returns the breaks of customer ownerID among ids; ids of other customers'
// or deleted breaks are left out.
func (breakRepository *BreakRepository) FindByIDs(ctx context.Context, ownerID int64, ids []int64) ([]model.Break, error) {
	query := breakSelect + `
		WHERE b.id = ANY($1::bigint[]) AND b.user_id = $2 AND b.is_deleted = false
		ORDER BY b.start_time, b.id`

	rows, err := db.Execute(ctx, breakRepository.database, query, ids, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find breaks: %w", err)
	}
	breaks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Break])
	if err != nil {
		return nil, fmt.Errorf("find breaks: %w", err)
	}
	return breaks, nil
}

// FindByShifts returns the breaks the shifts shiftIDs take, by shift and start time.
func (breakRepository *BreakRepository) FindByShifts(ctx context.Context, shiftIDs []int64) ([]model.ShiftBreak, error) {
	if len(shiftIDs) == 0 {
		return []model.ShiftBreak{}, nil
	}
	query := `SELECT id, shift_id, name,
			to_char(start_time, 'HH24:MI') AS start_time, to_char(end_time, 'HH24:MI') AS end_time
		FROM breaks
		WHERE shift_id = ANY($1::bigint[]) AND is_deleted = false
		ORDER BY shift_id, start_time, id`

	rows, err := db.Execute(ctx, breakRepository.database, query, shiftIDs)
	if err != nil {
		return nil, fmt.Errorf("find shift breaks: %w", err)
	}
	breaks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ShiftBreak])
	if err != nil {
		return nil, fmt.Errorf("find shift breaks: %w", err)
	}
	return breaks, nil
}

// AssignShift makes breakIDs the breaks shiftID takes: its other breaks are released
// (shift_id NULL) and breakIDs get shift_id = shiftID.
func (breakRepository *BreakRepository) AssignShift(ctx context.Context, shiftID int64, breakIDs []int64) error {
	if breakIDs == nil {
		breakIDs = []int64{}
	}
	releaseQuery := `UPDATE breaks SET shift_id = NULL
		WHERE shift_id = $1 AND id <> ALL($2::bigint[])`
	if _, err := executeAffected(ctx, breakRepository.database, releaseQuery, shiftID, breakIDs); err != nil {
		return fmt.Errorf("release shift breaks: %w", err)
	}

	assignQuery := `UPDATE breaks SET shift_id = $1
		WHERE id = ANY($2::bigint[]) AND shift_id IS DISTINCT FROM $1`
	if _, err := executeAffected(ctx, breakRepository.database, assignQuery, shiftID, breakIDs); err != nil {
		return fmt.Errorf("assign shift breaks: %w", err)
	}
	return nil
}

// ReleaseShift sets shift_id back to NULL on every break shiftID takes.
func (breakRepository *BreakRepository) ReleaseShift(ctx context.Context, shiftID int64) error {
	query := `UPDATE breaks SET shift_id = NULL WHERE shift_id = $1`
	if _, err := executeAffected(ctx, breakRepository.database, query, shiftID); err != nil {
		return fmt.Errorf("release shift breaks: %w", err)
	}
	return nil
}

// Create inserts shiftBreak and returns the new id.
func (breakRepository *BreakRepository) Create(ctx context.Context, shiftBreak model.Break) (int64, error) {
	query := `INSERT INTO breaks (user_id, shift_id, name, start_time, end_time, status)
		VALUES ($1, $2, $3, $4::time, $5::time, $6)
		RETURNING id`

	rows, err := db.Execute(ctx, breakRepository.database, query,
		shiftBreak.UserID, shiftBreak.ShiftID, shiftBreak.Name, shiftBreak.StartTime, shiftBreak.EndTime, shiftBreak.Status)
	if err != nil {
		return 0, fmt.Errorf("create break: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapBreakError("create break", err)
}

// Update replaces the fields of shiftBreak.ID. A blank Status keeps the saved status.
func (breakRepository *BreakRepository) Update(ctx context.Context, shiftBreak model.Break) error {
	query := `UPDATE breaks
		SET user_id = $2, shift_id = $3, name = $4, start_time = $5::time, end_time = $6::time,
			status = COALESCE(NULLIF($7, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, breakRepository.database, query,
		shiftBreak.ID, shiftBreak.UserID, shiftBreak.ShiftID, shiftBreak.Name, shiftBreak.StartTime, shiftBreak.EndTime, shiftBreak.Status)
	if err != nil {
		return mapBreakError("update break", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the break of ownerID (any break when nil).
func (breakRepository *BreakRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE breaks
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, breakRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete break: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapBreakError turns driver errors into repository errors. A nil err stays nil.
func mapBreakError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "breaks_user_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "breaks_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		case pgError.Code == "23503" && pgError.ConstraintName == "breaks_shift_id_fkey":
			return ErrInvalidParent
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
