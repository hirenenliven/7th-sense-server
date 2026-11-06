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

// machineModeReasonSelect reads model.MachineModeReason: the machine_mode_reasons row (r) with its
// reason from whichever master list it points at. The masters are joined without is_deleted, so
// a reason deleted from its list still shows on the periods it explains.
const machineModeReasonSelect = `SELECT r.id, r.mode_log_id AS period_id,
		CASE WHEN r.downtime_reason_id IS NOT NULL THEN 'downtime' ELSE 'idle' END AS reason_type,
		COALESCE(r.downtime_reason_id, r.idle_reason_id)::bigint AS reason_id,
		COALESCE(d.reason_code, i.reason_code, '') AS reason_code,
		COALESCE(d.reason, i.reason, '') AS reason,
		r.started_at, r.ended_at, r.is_deleted, r.deleted_at, r.created_at, r.updated_at
	FROM machine_mode_reasons r
	LEFT JOIN downtime_reasons d ON d.id = r.downtime_reason_id
	LEFT JOIN idle_reasons i ON i.id = r.idle_reason_id`

// MachineModeReasonRepository runs the SQL queries for the machine_mode_reasons table. Soft-deleted
// rows (is_deleted = true) are invisible to every method.
type MachineModeReasonRepository struct {
	database db.Querier
}

// NewMachineModeReasonRepository creates a MachineModeReasonRepository that runs its queries on
// database (a connection pool or a transaction).
func NewMachineModeReasonRepository(database db.Querier) *MachineModeReasonRepository {
	return &MachineModeReasonRepository{database: database}
}

// FindPeriodByID returns the whole machine_mode_log period, or ErrNotFound.
func (reasonRepository *MachineModeReasonRepository) FindPeriodByID(ctx context.Context, id int64) (model.ReasonPeriod, error) {
	query := `SELECT l.id, l.machine_id::bigint AS machine_id, l.mode_id, COALESCE(mm.name, '') AS mode_name,
			l.started_at, l.ended_at
		FROM machine_mode_log l
		LEFT JOIN machine_modes mm ON mm.id = l.mode_id
		WHERE l.id = $1`

	rows, err := db.Execute(ctx, reasonRepository.database, query, id)
	if err != nil {
		return model.ReasonPeriod{}, fmt.Errorf("find mode period: %w", err)
	}
	period, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ReasonPeriod])
	return period, mapMachineModeReasonError("find mode period", err)
}

// LockPeriod locks the period row until the transaction ends, so two saves on one period cannot
// both pass the overlap check. Only finished periods take reasons and the collector never writes
// those again, so this does not wait on it.
func (reasonRepository *MachineModeReasonRepository) LockPeriod(ctx context.Context, id int64) error {
	query := `SELECT id FROM machine_mode_log WHERE id = $1 FOR UPDATE`

	rows, err := db.Execute(ctx, reasonRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("lock mode period: %w", err)
	}
	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return mapMachineModeReasonError("lock mode period", err)
}

// FindAllByPeriod returns the reasons of the period, in time order.
func (reasonRepository *MachineModeReasonRepository) FindAllByPeriod(ctx context.Context, periodID int64) ([]model.MachineModeReason, error) {
	query := machineModeReasonSelect + `
		WHERE r.mode_log_id = $1 AND r.is_deleted = false
		ORDER BY r.started_at, r.id`

	rows, err := db.Execute(ctx, reasonRepository.database, query, periodID)
	if err != nil {
		return nil, fmt.Errorf("find mode reasons: %w", err)
	}
	reasons, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineModeReason])
	if err != nil {
		return nil, fmt.Errorf("find mode reasons: %w", err)
	}
	return reasons, nil
}

// FindByID returns the reason row, or ErrNotFound.
func (reasonRepository *MachineModeReasonRepository) FindByID(ctx context.Context, id int64) (model.MachineModeReason, error) {
	query := machineModeReasonSelect + `
		WHERE r.id = $1 AND r.is_deleted = false`

	rows, err := db.Execute(ctx, reasonRepository.database, query, id)
	if err != nil {
		return model.MachineModeReason{}, fmt.Errorf("find mode reason: %w", err)
	}
	reason, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineModeReason])
	return reason, mapMachineModeReasonError("find mode reason", err)
}

// HasOverlap reports whether another reason of the period covers any time between startedAt and
// endedAt. Reasons that only touch do not overlap. exceptID (0 = none) is left out, for an update.
func (reasonRepository *MachineModeReasonRepository) HasOverlap(ctx context.Context, periodID int64, startedAt, endedAt time.Time, exceptID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM machine_mode_reasons
		WHERE mode_log_id = $1 AND is_deleted = false AND id <> $4
		  AND started_at < $3 AND ended_at > $2)`

	rows, err := db.Execute(ctx, reasonRepository.database, query, periodID, startedAt, endedAt, exceptID)
	if err != nil {
		return false, fmt.Errorf("check mode reason overlap: %w", err)
	}
	overlaps, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check mode reason overlap: %w", err)
	}
	return overlaps, nil
}

// Create inserts a reason of type reasonType for the period and returns the new id.
func (reasonRepository *MachineModeReasonRepository) Create(ctx context.Context, periodID int64, reasonType string, reasonID int64, startedAt, endedAt time.Time) (int64, error) {
	query := `INSERT INTO machine_mode_reasons (mode_log_id, downtime_reason_id, idle_reason_id, started_at, ended_at)
		VALUES ($1, CASE WHEN $2 = 'downtime' THEN $3::integer END, CASE WHEN $2 = 'idle' THEN $3::integer END, $4, $5)
		RETURNING id`

	rows, err := db.Execute(ctx, reasonRepository.database, query, periodID, reasonType, reasonID, startedAt, endedAt)
	if err != nil {
		return 0, fmt.Errorf("create mode reason: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapMachineModeReasonError("create mode reason", err)
}

// Update replaces the reason and the time of row id. A period's mode never changes, so the row
// stays of reasonType.
func (reasonRepository *MachineModeReasonRepository) Update(ctx context.Context, id int64, reasonType string, reasonID int64, startedAt, endedAt time.Time) error {
	query := `UPDATE machine_mode_reasons
		SET downtime_reason_id = CASE WHEN $2 = 'downtime' THEN $3::integer END,
			idle_reason_id = CASE WHEN $2 = 'idle' THEN $3::integer END,
			started_at = $4, ended_at = $5
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, reasonRepository.database, query, id, reasonType, reasonID, startedAt, endedAt)
	if err != nil {
		return fmt.Errorf("update mode reason: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the reason row.
func (reasonRepository *MachineModeReasonRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE machine_mode_reasons
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, reasonRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete mode reason: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapMachineModeReasonError turns driver errors into repository errors. A nil err stays nil.
func mapMachineModeReasonError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
