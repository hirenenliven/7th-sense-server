package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// MachineDiagnosticsRepository reads what the CNC collector recorded about one machine, for the
// superadmin's diagnostics page. It writes nothing.
type MachineDiagnosticsRepository struct {
	database db.Querier
}

// NewMachineDiagnosticsRepository creates a MachineDiagnosticsRepository that runs its queries on
// database.
func NewMachineDiagnosticsRepository(database db.Querier) *MachineDiagnosticsRepository {
	return &MachineDiagnosticsRepository{database: database}
}

// FindEngine returns the cnc_engine of the machine's controller, nil when it has none.
func (diagnosticsRepository *MachineDiagnosticsRepository) FindEngine(ctx context.Context, machineID int64) (*int, error) {
	query := `SELECT c.cnc_engine::int
		FROM machines m
		LEFT JOIN controllers c ON c.id = m.controller_id AND c.is_deleted = false
		WHERE m.id = $1`

	rows, err := db.Execute(ctx, diagnosticsRepository.database, query, machineID)
	if err != nil {
		return nil, fmt.Errorf("find machine engine: %w", err)
	}
	engine, err := pgx.CollectOneRow(rows, pgx.RowTo[*int])
	if err != nil {
		return nil, fmt.Errorf("find machine engine: %w", err)
	}
	return engine, nil
}

// FindLiveStatus returns the machine's machine_live_status row, nil before its first report.
func (diagnosticsRepository *MachineDiagnosticsRepository) FindLiveStatus(ctx context.Context, machineID int64) (*model.DiagnosticLiveStatus, error) {
	query := `SELECT is_connected, connected_since, last_seen_at, last_change_at,
			callback_count, stored_count, dropped_count,
			error_count, COALESCE(last_error, '') AS last_error, updated_at
		FROM machine_live_status
		WHERE machine_id = $1`

	rows, err := db.Execute(ctx, diagnosticsRepository.database, query, machineID)
	if err != nil {
		return nil, fmt.Errorf("find machine live status: %w", err)
	}
	status, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.DiagnosticLiveStatus])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find machine live status: %w", err)
	}
	return &status, nil
}

// FindCurrentMode returns the machine's open machine_mode_log period, nil when it has none.
func (diagnosticsRepository *MachineDiagnosticsRepository) FindCurrentMode(ctx context.Context, machineID int64) (*model.DiagnosticMode, error) {
	query := `SELECT l.mode_id, COALESCE(mm.name, '') AS mode_name,
			COALESCE(l.program_name, '') AS program_name, l.started_at
		FROM machine_mode_log l
		LEFT JOIN machine_modes mm ON mm.id = l.mode_id
		WHERE l.machine_id = $1 AND l.ended_at IS NULL`

	rows, err := db.Execute(ctx, diagnosticsRepository.database, query, machineID)
	if err != nil {
		return nil, fmt.Errorf("find machine current mode: %w", err)
	}
	mode, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.DiagnosticMode])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find machine current mode: %w", err)
	}
	return &mode, nil
}

// FindVariables returns every active variable as the machine would read it. It picks the address
// the collector picks (MachineTelemetryRepository.FindVariableConfig), but keeps variables turned
// off for the machine and adds the last stored number, so the page can show why each is or is not
// read. sources are what the collector can read; engine may be nil (no address then).
func (diagnosticsRepository *MachineDiagnosticsRepository) FindVariables(ctx context.Context, machineID int64, engine *int, versionType int, sources []string) ([]model.DiagnosticVariable, error) {
	query := `SELECT v.id, v.code, v.name, COALESCE(v.unit, '') AS unit, v.category, v.data_type,
			COALESCE(mv.is_enabled, true) AS enabled,
			COALESCE(mv.source, a.source, '') AS source,
			COALESCE(mv.address, a.address, '') AS address,
			COALESCE((SELECT string_agg(DISTINCT o.source::text, ', ')
				FROM cnc_variable_addresses o
				WHERE o.variable_id = v.id AND o.cnc_engine = $2 AND o.is_deleted = false), '') AS known_sources,
			val.value AS stored_value, val.updated_at AS stored_at
		FROM cnc_variables v
		LEFT JOIN machine_variables mv ON mv.id = v.id AND mv.machine_id = $1
		LEFT JOIN machine_variable_values val ON val.variable_id = v.id AND val.machine_id = $1
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
		ORDER BY v.id`

	rows, err := db.Execute(ctx, diagnosticsRepository.database, query, machineID, engine, versionType, sources)
	if err != nil {
		return nil, fmt.Errorf("find machine variables: %w", err)
	}
	variables, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.DiagnosticVariable])
	if err != nil {
		return nil, fmt.Errorf("find machine variables: %w", err)
	}
	return variables, nil
}

// FindRecentEvents returns the machine's newest machine_events rows, newest first.
func (diagnosticsRepository *MachineDiagnosticsRepository) FindRecentEvents(ctx context.Context, machineID int64, limit int) ([]model.DiagnosticEvent, error) {
	query := `SELECT id, event_type, code, COALESCE(message, '') AS message, created_at
		FROM machine_events
		WHERE machine_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2`

	rows, err := db.Execute(ctx, diagnosticsRepository.database, query, machineID, limit)
	if err != nil {
		return nil, fmt.Errorf("find machine events: %w", err)
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.DiagnosticEvent])
	if err != nil {
		return nil, fmt.Errorf("find machine events: %w", err)
	}
	return events, nil
}
