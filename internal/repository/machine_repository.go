package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// machineColumns is the column list read into model.Machine. DATE and NUMERIC
// columns are cast to text because the model holds them as strings.
const machineColumns = `id, user_id, machine_type_id, product_id, name, serial_number, live_ip, port, username, password, local_ip,
	brand_id, controller_id, controller_version, accuracy_id, line_id,
	installation_date::text AS installation_date, bill_date::text AS bill_date, warranty_date::text AS warranty_date,
	COALESCE(amc, '') AS amc, is_insured, online_status, maintenance_required,
	hourly_cost::text AS hourly_cost, status, is_monitor,
	online_date::text AS online_date, offline_date, mode_id, controller_version_type,
	COALESCE(health_ids, '') AS health_ids,
	COALESCE(longitude, '') AS longitude,
	COALESCE(latitude, '') AS latitude,
	cycle_analysis, machine_system_type, security_mode,
	is_deleted, deleted_at, created_at, updated_at`

// MachineRepository runs the SQL queries for the machines table. Soft-deleted
// machines (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Machine, so nullable
// columns are read with COALESCE and given their column name with AS.
type MachineRepository struct {
	dbPool *pgxpool.Pool
}

// NewMachineRepository creates a MachineRepository that uses dbPool.
func NewMachineRepository(dbPool *pgxpool.Pool) *MachineRepository {
	return &MachineRepository{dbPool: dbPool}
}

func (machineRepository *MachineRepository) FindAll(ctx context.Context) ([]model.Machine, error) {
	query := `SELECT ` + machineColumns + `
		FROM machines
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query)
	if err != nil {
		return nil, fmt.Errorf("find all machines: %w", err)
	}
	machines, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Machine])
	if err != nil {
		return nil, fmt.Errorf("find all machines: %w", err)
	}
	return machines, nil
}

func (machineRepository *MachineRepository) FindByID(ctx context.Context, id int64) (model.Machine, error) {
	query := `SELECT ` + machineColumns + `
		FROM machines
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, id)
	if err != nil {
		return model.Machine{}, fmt.Errorf("find machine: %w", err)
	}
	machine, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Machine])
	return machine, mapMachineError("find machine", err)
}

func (machineRepository *MachineRepository) Create(ctx context.Context, machine model.Machine) (model.Machine, error) {
	query := `INSERT INTO machines (user_id, machine_type_id, product_id, name, serial_number, live_ip, port, username, password, local_ip,
			brand_id, controller_id, controller_version, accuracy_id, line_id, installation_date, bill_date, warranty_date,
			amc, is_insured, online_status, maintenance_required, hourly_cost, status, is_monitor, online_date, mode_id,
			controller_version_type, health_ids, longitude, latitude, cycle_analysis, machine_system_type, security_mode)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34)
		RETURNING ` + machineColumns

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.UserID, machine.MachineTypeID, machine.ProductID, machine.Name, machine.SerialNumber, machine.LiveIP,
		machine.Port, machine.Username, machine.Password, machine.LocalIP, machine.BrandID, machine.ControllerID,
		machine.ControllerVersion, machine.AccuracyID, machine.LineID, machine.InstallationDate, machine.BillDate,
		machine.WarrantyDate, machine.AMC, machine.IsInsured, machine.OnlineStatus, machine.MaintenanceRequired,
		machine.HourlyCost, machine.Status, machine.IsMonitor, machine.OnlineDate, machine.ModeID,
		machine.ControllerVersionType, machine.HealthIDs, machine.Longitude, machine.Latitude, machine.CycleAnalysis,
		machine.MachineSystemType, machine.SecurityMode,
	)
	if err != nil {
		return model.Machine{}, fmt.Errorf("create machine: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Machine])
	return created, mapMachineError("create machine", err)
}

// Update replaces the fields of machine.ID.
func (machineRepository *MachineRepository) Update(ctx context.Context, machine model.Machine) (model.Machine, error) {
	query := `UPDATE machines
		SET user_id = $2, machine_type_id = $3, product_id = $4, name = $5, serial_number = $6, live_ip = $7, port = $8,
			username = $9, password = $10, local_ip = $11, brand_id = $12, controller_id = $13, controller_version = $14,
			accuracy_id = $15, line_id = $16, installation_date = $17, bill_date = $18, warranty_date = $19,
			amc = NULLIF($20, ''), is_insured = $21, online_status = $22, maintenance_required = $23,
			hourly_cost = $24, status = $25, is_monitor = $26, online_date = $27, mode_id = $28,
			controller_version_type = $29, health_ids = NULLIF($30, ''), longitude = NULLIF($31, ''),
			latitude = NULLIF($32, ''), cycle_analysis = $33, machine_system_type = $34, security_mode = $35
		WHERE id = $1 AND is_deleted = false
		RETURNING ` + machineColumns

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.ID, machine.UserID, machine.MachineTypeID, machine.ProductID, machine.Name, machine.SerialNumber,
		machine.LiveIP, machine.Port, machine.Username, machine.Password, machine.LocalIP, machine.BrandID,
		machine.ControllerID, machine.ControllerVersion, machine.AccuracyID, machine.LineID, machine.InstallationDate,
		machine.BillDate, machine.WarrantyDate, machine.AMC, machine.IsInsured, machine.OnlineStatus,
		machine.MaintenanceRequired, machine.HourlyCost, machine.Status, machine.IsMonitor, machine.OnlineDate,
		machine.ModeID, machine.ControllerVersionType, machine.HealthIDs, machine.Longitude, machine.Latitude,
		machine.CycleAnalysis, machine.MachineSystemType, machine.SecurityMode,
	)
	if err != nil {
		return model.Machine{}, fmt.Errorf("update machine: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Machine])
	return updated, mapMachineError("update machine", err)
}

// Delete soft-deletes the machine: the row is kept with is_deleted = true.
func (machineRepository *MachineRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE machines
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete machine: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete machine: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapMachineError turns driver errors into repository errors. A nil err stays nil.
func mapMachineError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "machines_name_key":
			return ErrDuplicateName
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
