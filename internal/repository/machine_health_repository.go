package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// MachineHealthRepository runs the read-only queries of the Machine Health charts that read the
// collector's telemetry directly (today the X / Y / Z axis charts); the CNC collector writes those
// tables (no soft delete).
type MachineHealthRepository struct {
	database db.Querier
}

// NewMachineHealthRepository creates a MachineHealthRepository that runs its queries on database.
func NewMachineHealthRepository(database db.Querier) *MachineHealthRepository {
	return &MachineHealthRepository{database: database}
}

// FindAxisVariables returns the active cnc_variables with codes (any order; unknown codes are left out).
func (healthRepository *MachineHealthRepository) FindAxisVariables(ctx context.Context, codes []string) ([]model.AxisVariable, error) {
	query := `SELECT id, code, name, COALESCE(unit, '') AS unit
		FROM cnc_variables
		WHERE code = ANY($1) AND is_deleted = false`

	rows, err := db.Execute(ctx, healthRepository.database, query, codes)
	if err != nil {
		return nil, fmt.Errorf("find axis variables: %w", err)
	}
	variables, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AxisVariable])
	if err != nil {
		return nil, fmt.Errorf("find axis variables: %w", err)
	}
	return variables, nil
}

// FindAxisRows returns, per variable, its first limit stored values from from to to (oldest first),
// ordered by variable and time: one index range scan each.
func (healthRepository *MachineHealthRepository) FindAxisRows(ctx context.Context, machineID int64, variableIDs []int64, from, to time.Time, limit int) ([]model.AxisRow, error) {
	query := `SELECT ids.variable_id, axis_row.recorded_at, axis_row.value
		FROM unnest($2::bigint[]) AS ids(variable_id)
		CROSS JOIN LATERAL (
			SELECT t.recorded_at, t.value
			FROM machine_telementry t
			WHERE t.machine_id = $1 AND t.variable_id = ids.variable_id
				AND t.recorded_at >= $3 AND t.recorded_at < $4 AND t.value IS NOT NULL
			ORDER BY t.recorded_at
			LIMIT $5
		) AS axis_row
		ORDER BY ids.variable_id, axis_row.recorded_at`

	rows, err := db.Execute(ctx, healthRepository.database, query, machineID, variableIDs, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("find axis rows: %w", err)
	}
	axisRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AxisRow])
	if err != nil {
		return nil, fmt.Errorf("find axis rows: %w", err)
	}
	return axisRows, nil
}

// FindAxisBuckets cuts from–to into buckets of bucketSeconds (starting at from) and returns, per
// variable and bucket with rows, the last value, the lowest and highest value and the row count,
// ordered by variable and time. Uses the (machine_id, variable_id, recorded_at) index (migration 038).
func (healthRepository *MachineHealthRepository) FindAxisBuckets(ctx context.Context, machineID int64, variableIDs []int64, from, to time.Time, bucketSeconds int64) ([]model.AxisBucket, error) {
	// The last value comes from the bucket's newest timestamp, looked up by index afterwards:
	// TimescaleDB last(value, recorded_at) gave the same result about 3.5x slower over 30 days
	// of a machine (3 million rows)
	query := `WITH buckets AS (
			SELECT variable_id,
				time_bucket(make_interval(secs => $5), recorded_at, $3::timestamptz) AS bucket_at,
				MIN(value) AS min_value,
				MAX(value) AS max_value,
				COUNT(*) AS row_count,
				MAX(recorded_at) AS last_at
			FROM machine_telementry
			WHERE machine_id = $1 AND variable_id = ANY($2)
				AND recorded_at >= $3 AND recorded_at < $4
				AND value IS NOT NULL
			GROUP BY variable_id, bucket_at
		)
		SELECT b.variable_id, b.bucket_at, last_row.value AS last_value, b.min_value, b.max_value, b.row_count
		FROM buckets b
		CROSS JOIN LATERAL (
			SELECT t.value
			FROM machine_telementry t
			WHERE t.machine_id = $1 AND t.variable_id = b.variable_id
				AND t.recorded_at = b.last_at AND t.value IS NOT NULL
			LIMIT 1
		) AS last_row
		ORDER BY b.variable_id, b.bucket_at`

	rows, err := db.Execute(ctx, healthRepository.database, query, machineID, variableIDs, from, to, bucketSeconds)
	if err != nil {
		return nil, fmt.Errorf("find axis buckets: %w", err)
	}
	buckets, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AxisBucket])
	if err != nil {
		return nil, fmt.Errorf("find axis buckets: %w", err)
	}
	return buckets, nil
}

// FindAxisStartValues returns, per variable, the last value stored before from (variables with none
// are left out): one index lookup each.
func (healthRepository *MachineHealthRepository) FindAxisStartValues(ctx context.Context, machineID int64, variableIDs []int64, from time.Time) ([]model.AxisStartValue, error) {
	query := `SELECT ids.variable_id, last_row.value, last_row.recorded_at
		FROM unnest($2::bigint[]) AS ids(variable_id)
		CROSS JOIN LATERAL (
			SELECT t.value, t.recorded_at
			FROM machine_telementry t
			WHERE t.machine_id = $1 AND t.variable_id = ids.variable_id
				AND t.recorded_at < $3 AND t.value IS NOT NULL
			ORDER BY t.recorded_at DESC
			LIMIT 1
		) AS last_row`

	rows, err := db.Execute(ctx, healthRepository.database, query, machineID, variableIDs, from)
	if err != nil {
		return nil, fmt.Errorf("find axis start values: %w", err)
	}
	values, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AxisStartValue])
	if err != nil {
		return nil, fmt.Errorf("find axis start values: %w", err)
	}
	return values, nil
}
