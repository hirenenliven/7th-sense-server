package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// MachineCycleRepository runs the SQL queries of machine_cycles (migration 030) and the job a
// cycle is matched to. Cycles are collector data: no soft delete, so no is_deleted filter on
// machine_cycles itself.
type MachineCycleRepository struct {
	database db.Querier
}

// NewMachineCycleRepository creates a MachineCycleRepository that runs its queries on database
// (a connection pool or a transaction).
func NewMachineCycleRepository(database db.Querier) *MachineCycleRepository {
	return &MachineCycleRepository{database: database}
}

// machineCycleSelect reads cycles with the names of their machine, job, part, process and shift.
// duration_sec runs up to now while the cycle is open.
const machineCycleSelect = `SELECT c.id, c.machine_id, COALESCE(m.name, '') AS machine_name,
		COALESCE(c.program_name, '') AS program_name,
		c.schedule_item_id, COALESCE(i.job_id, '') AS job_id,
		c.part_id, COALESCE(p.part_name, '') AS part_name,
		c.process_id, COALESCE(pp.process, '') AS process_name,
		c.shift_id, COALESCE(s.name, '') AS shift_name,
		COALESCE(to_char(c.shift_date, 'YYYY-MM-DD'), '') AS shift_date,
		c.standard_cycle_sec, c.status, COALESCE(c.end_reason, '') AS end_reason,
		c.started_at, c.ended_at,
		EXTRACT(EPOCH FROM COALESCE(c.ended_at, now()) - c.started_at)::int AS duration_sec,
		c.production_sec, c.hold_sec, c.downtime_sec, c.offline_sec, c.loading_sec,
		c.part_count_start, c.part_count_end, c.parts_made, c.created_at, c.updated_at
	FROM machine_cycles c
	JOIN machines m ON m.id = c.machine_id
	LEFT JOIN production_schedule_items i ON i.id = c.schedule_item_id
	LEFT JOIN parts p ON p.id = c.part_id
	LEFT JOIN part_processes pp ON pp.id = c.process_id
	LEFT JOIN shifts s ON s.id = c.shift_id`

// processStandardCycleSec is a part process's planned cycle time: its cycle time plus the cycle
// times of its sub programs (seconds). %s is the part_processes alias.
const processStandardCycleSec = `(%[1]s.cycle_time + COALESCE((SELECT SUM(sp.cycle_time) FROM part_sub_programs sp
		WHERE sp.process_id = %[1]s.id AND sp.is_deleted = false), 0))::int`

// FindAll returns machine machineID's cycles that started in from … to (to excluded), oldest first.
func (cycleRepository *MachineCycleRepository) FindAll(ctx context.Context, machineID int64, from, to time.Time) ([]model.MachineCycle, error) {
	query := machineCycleSelect + ` WHERE c.machine_id = $1 AND c.started_at >= $2 AND c.started_at < $3
		ORDER BY c.started_at, c.id`
	return cycleRepository.findCycles(ctx, "find machine cycles", query, machineID, from, to)
}

// FindByID returns cycle id, or ErrNotFound.
func (cycleRepository *MachineCycleRepository) FindByID(ctx context.Context, id int64) (model.MachineCycle, error) {
	return cycleRepository.findCycle(ctx, "find machine cycle", machineCycleSelect+` WHERE c.id = $1`, id)
}

// FindOpenByMachineID returns machine machineID's running cycle, or ErrNotFound.
func (cycleRepository *MachineCycleRepository) FindOpenByMachineID(ctx context.Context, machineID int64) (model.MachineCycle, error) {
	query := machineCycleSelect + ` WHERE c.machine_id = $1 AND c.status = $2`
	return cycleRepository.findCycle(ctx, "find open machine cycle", query, machineID, model.CycleRunning)
}

// FindLatestClosed returns machine machineID's cycle that ended last at or before before, or
// ErrNotFound.
func (cycleRepository *MachineCycleRepository) FindLatestClosed(ctx context.Context, machineID int64, before time.Time) (model.MachineCycle, error) {
	query := machineCycleSelect + ` WHERE c.machine_id = $1 AND c.status <> $2 AND c.ended_at <= $3
		ORDER BY c.ended_at DESC, c.id DESC LIMIT 1`
	return cycleRepository.findCycle(ctx, "find latest closed machine cycle", query, machineID, model.CycleRunning, before)
}

func (cycleRepository *MachineCycleRepository) findCycles(ctx context.Context, operation, query string, args ...any) ([]model.MachineCycle, error) {
	rows, err := db.Execute(ctx, cycleRepository.database, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	cycles, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineCycle])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return cycles, nil
}

func (cycleRepository *MachineCycleRepository) findCycle(ctx context.Context, operation, query string, args ...any) (model.MachineCycle, error) {
	rows, err := db.Execute(ctx, cycleRepository.database, query, args...)
	if err != nil {
		return model.MachineCycle{}, fmt.Errorf("%s: %w", operation, err)
	}
	cycle, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineCycle])
	return cycle, mapMachineCycleError(operation, err)
}

// Create adds cycle as a running cycle and returns its id.
func (cycleRepository *MachineCycleRepository) Create(ctx context.Context, cycle model.MachineCycle) (int64, error) {
	query := `INSERT INTO machine_cycles (machine_id, program_name, schedule_item_id, part_id, process_id,
			shift_id, shift_date, standard_cycle_sec, status, started_at, loading_sec, part_count_start,
			part_count_end, parts_made)
		VALUES ($1, NULLIF(LEFT($2, 64), ''), $3, $4, $5, $6, NULLIF($7, '')::date, $8, $9, $10, $11, $12, $12, 0)
		RETURNING id`
	rows, err := db.Execute(ctx, cycleRepository.database, query, cycle.MachineID, cycle.ProgramName,
		cycle.ScheduleItemID, cycle.PartID, cycle.ProcessID, cycle.ShiftID, cycle.ShiftDate,
		cycle.StandardCycleSec, model.CycleRunning, cycle.StartedAt, cycle.LoadingSec, cycle.PartCountStart)
	if err != nil {
		return 0, fmt.Errorf("create machine cycle: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapMachineCycleError("create machine cycle", err)
}

// Update writes what the collector counts while a cycle runs and how it ended: the seconds per
// mode, the part counter, and status, end reason and end time. Returns ErrNotFound when the
// cycle is gone.
func (cycleRepository *MachineCycleRepository) Update(ctx context.Context, cycle model.MachineCycle) error {
	query := `UPDATE machine_cycles
		SET status = $2, end_reason = NULLIF($3, ''), ended_at = $4,
			production_sec = $5, hold_sec = $6, downtime_sec = $7, offline_sec = $8,
			part_count_end = $9, parts_made = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`
	affected, err := executeAffected(ctx, cycleRepository.database, query, cycle.ID, cycle.Status, cycle.EndReason,
		cycle.EndedAt, cycle.ProductionSec, cycle.HoldSec, cycle.DowntimeSec, cycle.OfflineSec,
		cycle.PartCountEnd, cycle.PartsMade)
	if err != nil {
		return mapMachineCycleError("update machine cycle", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// FindScheduledJob returns the job planned on machine machineID at at: an active production
// schedule item whose run (schedule date = the day its shift starts, so a 02:00 job of a
// 22:00 – 06:00 shift runs the next morning) holds at, the one whose process has programName
// first. ErrNotFound when none. Times are compared in the database session's time zone.
func (cycleRepository *MachineCycleRepository) FindScheduledJob(ctx context.Context, machineID int64, at time.Time, programName string) (model.CycleJob, error) {
	query := `SELECT i.id AS schedule_item_id, i.part_id, i.process_id, i.shift_id,
			to_char(i.schedule_date, 'YYYY-MM-DD') AS shift_date,
			` + fmt.Sprintf(processStandardCycleSec, "pp") + ` AS standard_cycle_sec
		FROM production_schedule_items i
		JOIN production_schedules ps ON ps.id = i.schedule_id AND ps.is_deleted = false AND ps.status = 'active'
		JOIN shifts s ON s.id = i.shift_id
		JOIN part_processes pp ON pp.id = i.process_id
		CROSS JOIN LATERAL (SELECT
				i.schedule_date + s.start_time + make_interval(mins =>
					mod((EXTRACT(EPOCH FROM i.start_time - s.start_time) / 60)::int + 1440, 1440)) AS job_start,
				mod((EXTRACT(EPOCH FROM i.end_time - i.start_time) / 60)::int + 1440, 1440) AS job_minutes) span
		WHERE ps.machine_id = $1 AND i.is_deleted = false
		  AND $2::timestamptz::timestamp >= span.job_start
		  AND $2::timestamptz::timestamp < span.job_start + make_interval(mins => span.job_minutes)
		ORDER BY lower(pp.program_name) = lower($3) DESC, span.job_start DESC
		LIMIT 1`
	rows, err := db.Execute(ctx, cycleRepository.database, query, machineID, at, programName)
	if err != nil {
		return model.CycleJob{}, fmt.Errorf("find scheduled job: %w", err)
	}
	job, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CycleJob])
	return job, mapMachineCycleError("find scheduled job", err)
}

// FindProcessesByProgram returns up to two part processes of active parts that list machine
// machineID among their applicable machines and whose program is programName
// (case-insensitive), without shift.
func (cycleRepository *MachineCycleRepository) FindProcessesByProgram(ctx context.Context, machineID int64, programName string) ([]model.CycleJob, error) {
	query := `SELECT NULL::bigint AS schedule_item_id, p.id::bigint AS part_id, pp.id::bigint AS process_id,
			NULL::bigint AS shift_id, '' AS shift_date,
			` + fmt.Sprintf(processStandardCycleSec, "pp") + ` AS standard_cycle_sec
		FROM part_processes pp
		JOIN part_process_machines ppm ON ppm.process_id = pp.id AND ppm.machine_id = $1
		JOIN parts p ON p.id = pp.part_id AND p.is_deleted = false AND p.status = 'active'
		WHERE pp.is_deleted = false AND lower(pp.program_name) = lower($2)
		ORDER BY pp.id
		LIMIT 2`
	rows, err := db.Execute(ctx, cycleRepository.database, query, machineID, programName)
	if err != nil {
		return nil, fmt.Errorf("find processes by program: %w", err)
	}
	jobs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CycleJob])
	if err != nil {
		return nil, fmt.Errorf("find processes by program: %w", err)
	}
	return jobs, nil
}

// FindShift returns the active shift of machine machineID's customer that runs at at, with the
// day it started (a shift runs on its working days; one past midnight started the day
// before). ErrNotFound when none.
func (cycleRepository *MachineCycleRepository) FindShift(ctx context.Context, machineID int64, at time.Time) (model.CycleJob, error) {
	query := `SELECT NULL::bigint AS schedule_item_id, NULL::bigint AS part_id, NULL::bigint AS process_id,
			s.id::bigint AS shift_id, to_char(d.day, 'YYYY-MM-DD') AS shift_date, NULL::int AS standard_cycle_sec
		FROM machines m
		JOIN shifts s ON s.user_id = m.user_id AND s.is_deleted = false AND s.status = 'active'
		CROSS JOIN LATERAL (VALUES ($2::timestamptz::timestamp::date), ($2::timestamptz::timestamp::date - 1)) d(day)
		WHERE m.id = $1
		  AND EXTRACT(ISODOW FROM d.day)::int = ANY(s.working_days)
		  AND $2::timestamptz::timestamp >= d.day + s.start_time
		  AND $2::timestamptz::timestamp < d.day + s.start_time + make_interval(mins =>
				mod((EXTRACT(EPOCH FROM s.end_time - s.start_time) / 60)::int + 1440, 1440))
		ORDER BY d.day DESC, s.start_time
		LIMIT 1`
	rows, err := db.Execute(ctx, cycleRepository.database, query, machineID, at)
	if err != nil {
		return model.CycleJob{}, fmt.Errorf("find shift: %w", err)
	}
	shift, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CycleJob])
	return shift, mapMachineCycleError("find shift", err)
}

// HasModeBetween reports whether machine machineID was in mode modeID at any time between from
// and to (machine_mode_log).
func (cycleRepository *MachineCycleRepository) HasModeBetween(ctx context.Context, machineID int64, modeID int, from, to time.Time) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM machine_mode_log
		WHERE machine_id = $1 AND mode_id = $2 AND started_at < $4 AND (ended_at IS NULL OR ended_at > $3))`
	rows, err := db.Execute(ctx, cycleRepository.database, query, machineID, modeID, from, to)
	if err != nil {
		return false, fmt.Errorf("check mode between: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check mode between: %w", err)
	}
	return found, nil
}

// mapMachineCycleError turns driver errors into repository errors. A nil err stays nil.
func mapMachineCycleError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23503": // foreign_key_violation: the machine, job, part or shift is gone
			return ErrNotFound
		case "23505": // unique_violation: machine_cycles_open_key, the machine already has a running cycle
			return ErrDuplicateCycle
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
