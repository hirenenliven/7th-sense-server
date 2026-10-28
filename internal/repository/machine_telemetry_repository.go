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

// MachineTelemetryRepository runs the SQL queries of the CNC collectors: the machines they
// connect to (active, not deleted), the variables they read, and the telemetry tables they
// write (migration 028). The data tables have no soft delete: they are append-only or one
// row per key.
//
// Selected column names must match the db tags on model.MonitoredMachine and
// model.VariableConfig, so nullable columns are read with COALESCE and named with AS.
type MachineTelemetryRepository struct {
	database db.Querier
}

// NewMachineTelemetryRepository creates a MachineTelemetryRepository that runs its queries on
// database (a connection pool or a transaction).
func NewMachineTelemetryRepository(database db.Querier) *MachineTelemetryRepository {
	return &MachineTelemetryRepository{database: database}
}

// activeMachineSelect reads the machines a collector connects to: active and not deleted.
// The engine comes from the machine's controller (NULL = not set).
const activeMachineSelect = `SELECT m.id, m.name, m.live_ip, m.port,
		COALESCE(m.username, '') AS username,
		COALESCE(m.password, '') AS password,
		COALESCE(m.controller_version_type, 0) AS controller_version_type,
		c.cnc_engine
	FROM machines m
	LEFT JOIN controllers c ON c.id = m.controller_id AND c.is_deleted = false
	WHERE m.is_deleted = false AND m.status = 'active'`

// FindActiveMachines returns the active machines of shard (id % shards = shard), by id.
func (telemetryRepository *MachineTelemetryRepository) FindActiveMachines(ctx context.Context, shard, shards int) ([]model.MonitoredMachine, error) {
	query := activeMachineSelect + ` AND m.id % $2::int = $1::int ORDER BY m.id`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, shard, shards)
	if err != nil {
		return nil, fmt.Errorf("find active machines: %w", err)
	}
	machines, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MonitoredMachine])
	if err != nil {
		return nil, fmt.Errorf("find active machines: %w", err)
	}
	return machines, nil
}

// FindActiveMachineByID returns the machine when it is active, else ErrNotFound.
func (telemetryRepository *MachineTelemetryRepository) FindActiveMachineByID(ctx context.Context, id int64) (model.MonitoredMachine, error) {
	query := activeMachineSelect + ` AND m.id = $1`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, id)
	if err != nil {
		return model.MonitoredMachine{}, fmt.Errorf("find active machine: %w", err)
	}
	machine, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MonitoredMachine])
	return machine, mapMachineTelemetryError("find active machine", err)
}

// mapMachineTelemetryError turns driver errors into repository errors. A nil err stays nil.
func mapMachineTelemetryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation: the machine is gone
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// FindVariableConfig returns the variables machine reads, with its overrides applied. Per
// variable the address is the one for the machine's controller version before the one for all
// versions, then the cheapest source; only sources in sources (what the collector can read)
// are used. A variable without a usable address comes back with an empty source and address,
// so the collector can report it instead of dropping it silently.
func (telemetryRepository *MachineTelemetryRepository) FindVariableConfig(ctx context.Context, machine model.MonitoredMachine, sources []string) ([]model.VariableConfig, error) {
	query := `SELECT v.id, v.code, v.data_type, v.store_mode,
			COALESCE(mv.deadband, v.deadband) AS deadband,
			COALESCE(mv.min_interval_ms, v.min_interval_ms) AS min_interval_ms,
			COALESCE(mv.source, a.source, '') AS source,
			COALESCE(mv.address, a.address, '') AS address
		FROM cnc_variables v
		LEFT JOIN machine_variables mv ON mv.id = v.id AND mv.machine_id = $1
		LEFT JOIN LATERAL (
			SELECT a.source, a.address
			FROM cnc_variable_addresses a
			WHERE a.variable_id = v.id AND a.cnc_engine = $2 AND a.is_deleted = false
			  AND (a.controller_version_type = $3 OR a.controller_version_type IS NULL)
			  AND a.source = ANY($4::text[])
			ORDER BY a.controller_version_type IS NULL,
			         array_position(ARRAY['snapshot', 'nc', 'plc'], a.source::text)
			LIMIT 1
		) a ON true
		WHERE v.is_deleted = false AND v.status = 'active'
		  AND COALESCE(mv.is_enabled, true)
		ORDER BY v.id`

	rows, err := db.Execute(ctx, telemetryRepository.database, query,
		machine.ID, machine.CNCEngine, machine.ControllerVersionType, sources)
	if err != nil {
		return nil, fmt.Errorf("find variable config: %w", err)
	}
	variables, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.VariableConfig])
	if err != nil {
		return nil, fmt.Errorf("find variable config: %w", err)
	}
	return variables, nil
}

// IsHypertable reports whether machine_telementry is a TimescaleDB hypertable.
func (telemetryRepository *MachineTelemetryRepository) IsHypertable(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb')`
	rows, err := db.Execute(ctx, telemetryRepository.database, query)
	if err != nil {
		return false, fmt.Errorf("check timescaledb: %w", err)
	}
	installed, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil || !installed {
		return false, err
	}

	query = `SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables
		WHERE hypertable_name = 'machine_telementry')`
	rows, err = db.Execute(ctx, telemetryRepository.database, query)
	if err != nil {
		return false, fmt.Errorf("check hypertable: %w", err)
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[bool])
}

// UpsertLiveValues upserts machine_variable_values rows: one per machine-variable,
// with the latest numeric value. Text variables are not cached here.
func (telemetryRepository *MachineTelemetryRepository) UpsertLiveValues(ctx context.Context, values []model.LiveValue) error {
	if len(values) == 0 {
		return nil
	}
	machineIDs := make([]int64, len(values))
	variableIDs := make([]int32, len(values))
	numbers := make([]*float64, len(values))
	for i, value := range values {
		machineIDs[i], variableIDs[i] = value.MachineID, int32(value.VariableID)
		numbers[i] = value.Value
	}

	query := `INSERT INTO machine_variable_values (machine_id, variable_id, value, updated_at)
		SELECT machine_id, variable_id, value, CURRENT_TIMESTAMP
		FROM unnest($1::bigint[], $2::int[], $3::float8[])
		AS t(machine_id, variable_id, value)
		ON CONFLICT (machine_id, variable_id) DO UPDATE
		SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP`
	if err := execute(ctx, telemetryRepository.database, query, machineIDs, variableIDs, numbers); err != nil {
		return fmt.Errorf("upsert live values: %w", err)
	}
	return nil
}

// GetLiveValues returns the current live values for a machine's variables.
func (telemetryRepository *MachineTelemetryRepository) GetLiveValues(ctx context.Context, machineID int64, variableIDs []int) (map[int]model.LiveValue, error) {
	if len(variableIDs) == 0 {
		return nil, nil
	}
	query := `SELECT machine_id, variable_id, value, updated_at
		FROM machine_variable_values
		WHERE machine_id = $1 AND variable_id = ANY($2::int[])`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, machineID, variableIDs)
	if err != nil {
		return nil, fmt.Errorf("get live values: %w", err)
	}
	liveValues, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.LiveValue])
	if err != nil {
		return nil, fmt.Errorf("get live values: %w", err)
	}

	result := make(map[int]model.LiveValue, len(liveValues))
	for _, lv := range liveValues {
		result[lv.VariableID] = lv
	}
	return result, nil
}

// InsertTelemetry adds rows to machine_telementry in one statement.
func (telemetryRepository *MachineTelemetryRepository) InsertTelemetry(ctx context.Context, values []model.TelemetryValue) error {
	if len(values) == 0 {
		return nil
	}
	machineIDs := make([]int64, len(values))
	variableIDs := make([]int32, len(values))
	numbers := make([]*float64, len(values))
	recordedAts := make([]time.Time, len(values))
	for i, value := range values {
		machineIDs[i], variableIDs[i] = value.MachineID, int32(value.VariableID)
		numbers[i], recordedAts[i] = value.Value, value.RecordedAt
	}

	query := `INSERT INTO machine_telementry (machine_id, variable_id, value, recorded_at)
		SELECT machine_id, variable_id, value, recorded_at
		FROM unnest($1::bigint[], $2::int[], $3::float8[], $4::timestamptz[])
		AS t(machine_id, variable_id, value, recorded_at)`
	if err := execute(ctx, telemetryRepository.database, query, machineIDs, variableIDs, numbers, recordedAts); err != nil {
		return fmt.Errorf("insert telemetry: %w", err)
	}
	return nil
}

// UpsertLiveStatuses saves one machine_live_status row per machine.
func (telemetryRepository *MachineTelemetryRepository) UpsertLiveStatuses(ctx context.Context, statuses []model.MachineLiveStatus) error {
	if len(statuses) == 0 {
		return nil
	}
	count := len(statuses)
	machineIDs, shards, modeIDs := make([]int64, count), make([]int32, count), make([]int32, count)
	connected := make([]bool, count)
	connectedSince, lastSeen, lastChange, sampleAt := make([]*time.Time, count), make([]*time.Time, count), make([]*time.Time, count), make([]*time.Time, count)
	missingNodes, lastErrors := make([]string, count), make([]string, count)
	callbacks, stored, dropped, errorCounts := make([]int64, count), make([]int64, count), make([]int64, count), make([]int64, count)
	for i, status := range statuses {
		machineIDs[i], shards[i], modeIDs[i], connected[i] = status.MachineID, int32(status.Shard), int32(status.ModeID), status.IsConnected
		connectedSince[i], lastSeen[i], lastChange[i], sampleAt[i] = status.ConnectedSince, status.LastSeenAt, status.LastChangeAt, status.SampleAt
		missingNodes[i], lastErrors[i] = status.MissingNodes, status.LastError
		callbacks[i], stored[i], dropped[i], errorCounts[i] = status.CallbackCount, status.StoredCount, status.DroppedCount, status.ErrorCount
	}

	query := `INSERT INTO machine_live_status (machine_id, collector_shard, is_connected, mode_id,
			connected_since, last_seen_at, last_change_at, sample_ts, missing_nodes,
			callback_count, stored_count, dropped_count, error_count, last_error, updated_at)
		SELECT s.machine_id, s.shard, s.connected, s.mode_id, s.connected_since, s.last_seen,
			s.last_change, s.sample_at, NULLIF(s.missing_nodes, ''), s.callbacks, s.stored,
			s.dropped, s.error_count, NULLIF(s.last_error, ''), CURRENT_TIMESTAMP
		FROM unnest($1::bigint[], $2::smallint[], $3::bool[], $4::int[], $5::timestamptz[],
			$6::timestamptz[], $7::timestamptz[], $8::timestamptz[], $9::text[], $10::bigint[],
			$11::bigint[], $12::bigint[], $13::bigint[], $14::text[])
			AS s(machine_id, shard, connected, mode_id, connected_since, last_seen, last_change,
				sample_at, missing_nodes, callbacks, stored, dropped, error_count, last_error)
		ON CONFLICT (machine_id) DO UPDATE
		SET collector_shard = EXCLUDED.collector_shard, is_connected = EXCLUDED.is_connected,
			mode_id = EXCLUDED.mode_id, connected_since = EXCLUDED.connected_since,
			last_seen_at = EXCLUDED.last_seen_at, last_change_at = EXCLUDED.last_change_at,
			sample_ts = EXCLUDED.sample_ts, missing_nodes = EXCLUDED.missing_nodes,
			callback_count = EXCLUDED.callback_count, stored_count = EXCLUDED.stored_count,
			dropped_count = EXCLUDED.dropped_count, error_count = EXCLUDED.error_count,
			last_error = EXCLUDED.last_error, updated_at = EXCLUDED.updated_at`
	err := execute(ctx, telemetryRepository.database, query, machineIDs, shards, connected, modeIDs,
		connectedSince, lastSeen, lastChange, sampleAt, missingNodes, callbacks, stored, dropped,
		errorCounts, lastErrors)
	if err != nil {
		return fmt.Errorf("upsert live statuses: %w", err)
	}
	return nil
}

// ChangeMode closes the machine's open mode period at change.At, opens one in change.ModeID
// and sets machines.mode_id / online_status (offline_date when it goes offline).
func (telemetryRepository *MachineTelemetryRepository) ChangeMode(ctx context.Context, change model.ModeChange) error {
	closeQuery := `UPDATE machine_mode_log
		SET ended_at = GREATEST(started_at, $2::timestamptz),
			duration_sec = EXTRACT(EPOCH FROM GREATEST(started_at, $2::timestamptz) - started_at)::int
		WHERE machine_id = $1 AND ended_at IS NULL`
	if err := execute(ctx, telemetryRepository.database, closeQuery, change.MachineID, change.At); err != nil {
		return fmt.Errorf("close mode period: %w", err)
	}

	openQuery := `INSERT INTO machine_mode_log (machine_id, mode_id, program_name, started_at)
		VALUES ($1, $2, NULLIF(LEFT($3, 64), ''), $4)`
	if err := execute(ctx, telemetryRepository.database, openQuery, change.MachineID, change.ModeID, change.ProgramName, change.At); err != nil {
		return fmt.Errorf("open mode period: %w", err)
	}

	machineQuery := `UPDATE machines
		SET mode_id = $2::int, online_status = ($2::int <> $3::int),
			offline_date = CASE WHEN $2::int = $3::int THEN $4::timestamptz ELSE offline_date END
		WHERE id = $1`
	if err := execute(ctx, telemetryRepository.database, machineQuery, change.MachineID, change.ModeID, model.ModeOffline, change.At); err != nil {
		return fmt.Errorf("set machine mode: %w", err)
	}
	return nil
}

// InsertEvents adds rows to machine_events in one statement.
func (telemetryRepository *MachineTelemetryRepository) InsertEvents(ctx context.Context, events []model.MachineEvent) error {
	if len(events) == 0 {
		return nil
	}
	count := len(events)
	machineIDs, shards, codes := make([]*int64, count), make([]int32, count), make([]*int32, count)
	types, messages, createdAts := make([]string, count), make([]string, count), make([]time.Time, count)
	for i, event := range events {
		machineIDs[i], shards[i], types[i] = event.MachineID, int32(event.Shard), event.Type
		messages[i], createdAts[i] = event.Message, event.At
		if event.Code != nil {
			code := int32(*event.Code)
			codes[i] = &code
		}
	}

	query := `INSERT INTO machine_events (machine_id, collector_shard, event_type, code, message, created_at)
		SELECT e.machine_id, e.shard, e.event_type, e.code, NULLIF(e.message, ''), e.created_at
		FROM unnest($1::bigint[], $2::smallint[], $3::text[], $4::int[], $5::text[], $6::timestamptz[])
			AS e(machine_id, shard, event_type, code, message, created_at)`
	if err := execute(ctx, telemetryRepository.database, query, machineIDs, shards, types, codes, messages, createdAts); err != nil {
		return fmt.Errorf("insert machine events: %w", err)
	}
	return nil
}

// MarkShardOffline makes every machine of shard (id % shards = shard) offline at at: an open
// period in another mode is closed and an Offline period opened; active machines without an
// open period get one; machine_live_status and machines.online_status say offline.
//
// When starting (a collector run begins), a period left open by a previous run closes at the
// machine's last_seen_at, so the time the collector was down counts as offline, and the live
// status counters restart (they count since the collector started). When stopping, the
// counters are kept for diagnosis.
func (telemetryRepository *MachineTelemetryRepository) MarkShardOffline(ctx context.Context, shard, shards int, at time.Time, starting bool) error {
	closeQuery := `WITH ended AS (
			SELECT l.id, l.machine_id,
				GREATEST(l.started_at, LEAST($3::timestamptz,
					CASE WHEN $4 THEN COALESCE(s.last_seen_at, $3::timestamptz) ELSE $3::timestamptz END)) AS ended_at
			FROM machine_mode_log l
			LEFT JOIN machine_live_status s ON s.machine_id = l.machine_id
			WHERE l.ended_at IS NULL AND l.mode_id <> $5::int AND l.machine_id % $2::int = $1::int
		), closed AS (
			UPDATE machine_mode_log l
			SET ended_at = e.ended_at,
				duration_sec = EXTRACT(EPOCH FROM e.ended_at - l.started_at)::int
			FROM ended e
			WHERE l.id = e.id
			RETURNING l.machine_id, e.ended_at
		)
		INSERT INTO machine_mode_log (machine_id, mode_id, started_at)
		SELECT machine_id, $5::int, ended_at FROM closed`
	if err := execute(ctx, telemetryRepository.database, closeQuery, shard, shards, at, starting, model.ModeOffline); err != nil {
		return fmt.Errorf("close mode periods: %w", err)
	}

	openQuery := `INSERT INTO machine_mode_log (machine_id, mode_id, started_at)
		SELECT m.id, $4::int, $3::timestamptz
		FROM machines m
		WHERE m.is_deleted = false AND m.status = 'active' AND m.id % $2::int = $1::int
		  AND NOT EXISTS (SELECT 1 FROM machine_mode_log l WHERE l.machine_id = m.id AND l.ended_at IS NULL)`
	if err := execute(ctx, telemetryRepository.database, openQuery, shard, shards, at, model.ModeOffline); err != nil {
		return fmt.Errorf("open offline periods: %w", err)
	}

	statusQuery := `INSERT INTO machine_live_status (machine_id, collector_shard, is_connected, mode_id, updated_at)
		SELECT m.id, $1::int, false, $3, CURRENT_TIMESTAMP
		FROM machines m
		WHERE m.is_deleted = false AND m.id % $2::int = $1::int
		  AND (m.status = 'active' OR EXISTS (SELECT 1 FROM machine_live_status s WHERE s.machine_id = m.id))
		ON CONFLICT (machine_id) DO UPDATE
		SET collector_shard = EXCLUDED.collector_shard, is_connected = false, mode_id = EXCLUDED.mode_id,
			connected_since = NULL,
			callback_count = CASE WHEN $4 THEN 0 ELSE machine_live_status.callback_count END,
			stored_count = CASE WHEN $4 THEN 0 ELSE machine_live_status.stored_count END,
			dropped_count = CASE WHEN $4 THEN 0 ELSE machine_live_status.dropped_count END,
			error_count = CASE WHEN $4 THEN 0 ELSE machine_live_status.error_count END,
			updated_at = EXCLUDED.updated_at`
	if err := execute(ctx, telemetryRepository.database, statusQuery, shard, shards, model.ModeOffline, starting); err != nil {
		return fmt.Errorf("reset live status: %w", err)
	}

	machineQuery := `UPDATE machines
		SET online_status = false, mode_id = $3::int,
			offline_date = CASE WHEN online_status THEN $4::timestamptz ELSE offline_date END
		WHERE is_deleted = false AND id % $2::int = $1::int
		  AND (status = 'active' OR online_status = true)`
	if err := execute(ctx, telemetryRepository.database, machineQuery, shard, shards, model.ModeOffline, at); err != nil {
		return fmt.Errorf("set machines offline: %w", err)
	}
	return nil
}

// FindMachineModeStatuses returns every machine that is not deleted with the mode it is in now,
// ordered by name. The open machine_mode_log period (at most one per machine) gives the mode,
// its start and program; machines.mode_id is used when there is none.
func (telemetryRepository *MachineTelemetryRepository) FindMachineModeStatuses(ctx context.Context) ([]model.MachineModeStatus, error) {
	query := `SELECT m.id AS machine_id, COALESCE(m.name, '') AS machine_name,
			COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
			COALESCE(cp.company_name, '') AS company_name,
			COALESCE(l.mode_id, m.mode_id) AS mode_id, COALESCE(mm.name, '') AS mode_name,
			COALESCE(m.online_status, false) AS online_status,
			COALESCE(l.program_name, '') AS program_name, l.started_at AS mode_since,
			EXTRACT(EPOCH FROM now() - l.started_at)::int AS mode_duration_sec
		FROM machines m
		LEFT JOIN machine_mode_log l ON l.machine_id = m.id AND l.ended_at IS NULL
		LEFT JOIN machine_modes mm ON mm.id = COALESCE(l.mode_id, m.mode_id)
		LEFT JOIN users u ON u.id = m.user_id
		LEFT JOIN company_profiles cp ON cp.user_id = m.user_id AND cp.is_deleted = false
		WHERE m.is_deleted = false
		ORDER BY m.name, m.id`

	rows, err := db.Execute(ctx, telemetryRepository.database, query)
	if err != nil {
		return nil, fmt.Errorf("find machine mode statuses: %w", err)
	}
	statuses, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineModeStatus])
	if err != nil {
		return nil, fmt.Errorf("find machine mode statuses: %w", err)
	}
	return statuses, nil
}

// FindModePeriods returns machine machineID's mode periods that were running at or after since,
// oldest first. The open period (ended_at NULL) counts its duration up to now.
func (telemetryRepository *MachineTelemetryRepository) FindModePeriods(ctx context.Context, machineID int64, since time.Time) ([]model.MachineModePeriod, error) {
	query := `SELECT l.id, l.mode_id, COALESCE(mm.name, '') AS mode_name,
			COALESCE(l.program_name, '') AS program_name, l.started_at, l.ended_at,
			COALESCE(l.duration_sec, EXTRACT(EPOCH FROM now() - l.started_at)::int) AS duration_sec
		FROM machine_mode_log l
		LEFT JOIN machine_modes mm ON mm.id = l.mode_id
		WHERE l.machine_id = $1 AND (l.ended_at IS NULL OR l.ended_at >= $2)
		ORDER BY l.started_at`

	rows, err := db.Execute(ctx, telemetryRepository.database, query, machineID, since)
	if err != nil {
		return nil, fmt.Errorf("find mode periods: %w", err)
	}
	periods, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineModePeriod])
	if err != nil {
		return nil, fmt.Errorf("find mode periods: %w", err)
	}
	return periods, nil
}
