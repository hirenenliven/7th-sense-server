package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// MachinePlcVariableRepository runs the SQL queries of a machine's PLC node ids: the defaults in
// cnc_variable_addresses (source plc) and the machine's own in machine_variables (source plc).
// machine_variables has no soft delete; its trigger NOTIFYs cnc_config_changed, so the machine's
// collector worker reloads its variables after every write.
type MachinePlcVariableRepository struct {
	database db.Querier
}

// NewMachinePlcVariableRepository creates a MachinePlcVariableRepository that runs its queries on
// database (a connection pool or a transaction).
func NewMachinePlcVariableRepository(database db.Querier) *MachinePlcVariableRepository {
	return &MachinePlcVariableRepository{database: database}
}

// plcDefaultOrder ranks a variable's plc addresses (alias a) as the machine's default node id, for
// a machine of engine engineSQL and controller version versionSQL: its own engine before any other
// (PLC defaults exist for Siemens only, and every machine reads a default until it has its own),
// then its version before all versions before another version, so every PLC variable has one. MachineTelemetryRepository.FindVariableConfig ranks plc addresses the
// same way, so the page shows the node id the collector uses.
func plcDefaultOrder(engineSQL, versionSQL string) string {
	return `a.cnc_engine IS DISTINCT FROM ` + engineSQL + `,
		a.controller_version_type IS DISTINCT FROM ` + versionSQL + `,
		a.controller_version_type IS NULL DESC, a.id`
}

// FindAll returns every PLC variable for the machine, by variable id: each variable with a plc
// address in cnc_variable_addresses for any engine. default_node_id is the best ranked plc address
// (plcDefaultOrder), custom_node_id the machine's own (machine_variables, source plc).
func (plcVariableRepository *MachinePlcVariableRepository) FindAll(ctx context.Context, machineID int64) ([]model.MachinePlcVariable, error) {
	query := `SELECT v.id, v.code, v.name, COALESCE(v.unit, '') AS unit, v.category, v.data_type, v.status,
			d.address AS default_node_id,
			COALESCE(CASE WHEN mv.source = 'plc' THEN mv.address END, '') AS custom_node_id
		FROM machines m
		LEFT JOIN controllers c ON c.id = m.controller_id AND c.is_deleted = false
		JOIN cnc_variables v ON v.is_deleted = false
		JOIN LATERAL (
			SELECT a.address
			FROM cnc_variable_addresses a
			WHERE a.variable_id = v.id AND a.source = 'plc' AND a.is_deleted = false
			ORDER BY ` + plcDefaultOrder("c.cnc_engine", "COALESCE(m.controller_version_type, 0)") + `
			LIMIT 1
		) d ON true
		LEFT JOIN machine_variables mv ON mv.machine_id = m.id AND mv.id = v.id
		WHERE m.id = $1 AND m.is_deleted = false
		ORDER BY v.id`

	rows, err := db.Execute(ctx, plcVariableRepository.database, query, machineID)
	if err != nil {
		return nil, fmt.Errorf("find machine plc variables: %w", err)
	}
	variables, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachinePlcVariable])
	if err != nil {
		return nil, fmt.Errorf("find machine plc variables: %w", err)
	}
	return variables, nil
}

// SaveNodeID stores nodeID as the machine's PLC node id of the variable, keeping the row's other
// overrides (is_enabled, deadband, min_interval_ms).
func (plcVariableRepository *MachinePlcVariableRepository) SaveNodeID(ctx context.Context, machineID int64, variableID int, nodeID string) error {
	query := `INSERT INTO machine_variables (machine_id, id, source, address)
		VALUES ($1, $2, 'plc', $3)
		ON CONFLICT (machine_id, id) DO UPDATE
		SET source = 'plc', address = EXCLUDED.address, updated_at = CURRENT_TIMESTAMP`

	if _, err := executeAffected(ctx, plcVariableRepository.database, query, machineID, variableID, nodeID); err != nil {
		return fmt.Errorf("save machine plc node id: %w", err)
	}
	return nil
}

// DeleteNodeIDs drops the machine's PLC node ids of the variables, so they read their default
// again. A row that only held the node id is deleted; one with other overrides keeps them. Run it
// in a transaction: it is two statements.
func (plcVariableRepository *MachinePlcVariableRepository) DeleteNodeIDs(ctx context.Context, machineID int64, variableIDs []int) error {
	deleteQuery := `DELETE FROM machine_variables
		WHERE machine_id = $1 AND id = ANY($2::int[]) AND source = 'plc'
		  AND is_enabled AND deadband IS NULL AND min_interval_ms IS NULL`
	if _, err := executeAffected(ctx, plcVariableRepository.database, deleteQuery, machineID, variableIDs); err != nil {
		return fmt.Errorf("delete machine plc node ids: %w", err)
	}

	clearQuery := `UPDATE machine_variables
		SET source = NULL, address = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE machine_id = $1 AND id = ANY($2::int[]) AND source = 'plc'`
	if _, err := executeAffected(ctx, plcVariableRepository.database, clearQuery, machineID, variableIDs); err != nil {
		return fmt.Errorf("clear machine plc node ids: %w", err)
	}
	return nil
}
