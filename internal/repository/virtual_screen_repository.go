package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// The queries of the Virtual Screen (GET /productivity/virtual-screen): a machine's live values,
// its last alarm and the progress of its job. Read only.

// FindLiveVariables returns machine machineID's latest values of the variables codes (the ones
// it has a value for).
func (telemetryRepository *MachineTelemetryRepository) FindLiveVariables(ctx context.Context, machineID int64, codes []string) ([]model.LiveVariable, error) {
	query := `SELECT cv.code, v.value, v.updated_at
		FROM machine_variable_values v
		JOIN cnc_variables cv ON cv.id = v.variable_id AND cv.is_deleted = false
		WHERE v.machine_id = $1 AND cv.code = ANY($2::text[])`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, machineID, codes)
	if err != nil {
		return nil, fmt.Errorf("find live variables: %w", err)
	}
	variables, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.LiveVariable])
	if err != nil {
		return nil, fmt.Errorf("find live variables: %w", err)
	}
	return variables, nil
}

// FindLatestAlarm returns the alarm machine machineID raised last, or nil when it never raised one.
func (telemetryRepository *MachineTelemetryRepository) FindLatestAlarm(ctx context.Context, machineID int64) (*model.MachineAlarm, error) {
	query := `SELECT e.code, COALESCE(e.message, '') AS message, e.created_at,
			NOT EXISTS (SELECT 1 FROM machine_events o
				WHERE o.machine_id = e.machine_id AND o.event_type = 'alarm_off' AND o.created_at >= e.created_at) AS active
		FROM machine_events e
		WHERE e.machine_id = $1 AND e.event_type = 'alarm_on'
		ORDER BY e.created_at DESC, e.id DESC
		LIMIT 1`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, machineID)
	if err != nil {
		return nil, fmt.Errorf("find latest alarm: %w", err)
	}
	alarms, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineAlarm])
	if err != nil {
		return nil, fmt.Errorf("find latest alarm: %w", err)
	}
	if len(alarms) == 0 {
		return nil, nil
	}
	return &alarms[0], nil
}

// FindJobProgress returns production schedule item scheduleItemID with its names, planned cycle
// time and completed cycles, or ErrNotFound.
func (cycleRepository *MachineCycleRepository) FindJobProgress(ctx context.Context, scheduleItemID int64) (model.JobProgress, error) {
	query := `SELECT i.id AS schedule_item_id, COALESCE(i.job_id, '') AS job_id,
			to_char(i.schedule_date, 'YYYY-MM-DD') AS schedule_date,
			to_char(i.start_time, 'HH24:MI') AS start_time, to_char(i.end_time, 'HH24:MI') AS end_time,
			i.quantity, COALESCE(p.part_name, '') AS part_name, COALESCE(pp.process, '') AS process_name,
			COALESCE(cl.name, '') AS clientele_name,
			` + fmt.Sprintf(processStandardCycleSec, "pp") + ` AS standard_cycle_sec,
			(SELECT COUNT(*) FROM machine_cycles c
				WHERE c.schedule_item_id = i.id AND c.status = '` + model.CycleCompleted + `')::int AS completed_cycles
		FROM production_schedule_items i
		JOIN part_processes pp ON pp.id = i.process_id
		LEFT JOIN parts p ON p.id = i.part_id
		LEFT JOIN clienteles cl ON cl.id = p.clientele_id
		WHERE i.id = $1 AND i.is_deleted = false`

	rows, err := db.Execute(ctx, cycleRepository.database, query, scheduleItemID)
	if err != nil {
		return model.JobProgress{}, fmt.Errorf("find job progress: %w", err)
	}
	job, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.JobProgress])
	return job, mapMachineCycleError("find job progress", err)
}
