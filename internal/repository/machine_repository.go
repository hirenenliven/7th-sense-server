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

// machineSelect reads model.Machine: the machines row (m) plus the display names of
// its customer, model, brand, controller and line. DATE and NUMERIC columns are cast
// to text because the model holds them as strings. Soft-deleted joined rows still
// give their name, so a machine keeps showing what it was saved with.
const machineSelect = `SELECT m.id, m.user_id, m.machine_type_id, m.model_id, m.name, m.serial_number,
		m.live_ip, m.port, m.username, m.password, m.password <> '' AS has_password, m.local_ip,
		m.brand_id, m.controller_id, m.controller_version, m.accuracy_id, m.line_id,
		m.installation_date::text AS installation_date, m.bill_date::text AS bill_date,
		m.warranty_date::text AS warranty_date,
		COALESCE(m.amc, '') AS amc, m.is_insured, m.online_status, m.maintenance_required,
		m.hourly_cost::text AS hourly_cost, m.status, m.is_monitor,
		m.online_date::text AS online_date, m.offline_date, m.mode_id, m.controller_version_type,
		COALESCE(m.health_ids, '') AS health_ids,
		COALESCE(m.longitude, '') AS longitude,
		COALESCE(m.latitude, '') AS latitude,
		m.cycle_analysis, m.machine_system_type, m.security_mode,
		m.is_deleted, m.deleted_at, m.created_at, m.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name,
		COALESCE(mo.name, '') AS model_name,
		COALESCE(b.name, '') AS brand_name,
		COALESCE(c.name, '') AS controller_name,
		COALESCE(l.name, '') AS line_name
	FROM machines m
	LEFT JOIN users u ON u.id = m.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = m.user_id AND cp.is_deleted = false
	LEFT JOIN models mo ON mo.id = m.model_id
	LEFT JOIN brands b ON b.id = m.brand_id
	LEFT JOIN controllers c ON c.id = m.controller_id
	LEFT JOIN user_line l ON l.id = m.line_id`

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

// FindAll returns the machines in scope, newest first.
func (machineRepository *MachineRepository) FindAll(ctx context.Context, scope model.MachineScope) ([]model.Machine, error) {
	query := machineSelect + `
		WHERE m.is_deleted = false
			AND ($1::integer IS NULL OR m.user_id = $1)
			AND ($2::integer IS NULL OR EXISTS (
				SELECT 1 FROM user_machine_access a WHERE a.machine_id = m.id AND a.user_id = $2))
		ORDER BY m.id DESC`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, scope.OwnerID, scope.AssignedUserID)
	if err != nil {
		return nil, fmt.Errorf("find all machines: %w", err)
	}
	machines, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Machine])
	if err != nil {
		return nil, fmt.Errorf("find all machines: %w", err)
	}
	return machines, nil
}

// FindByID returns the machine, or ErrNotFound when it is not in scope.
func (machineRepository *MachineRepository) FindByID(ctx context.Context, scope model.MachineScope, id int64) (model.Machine, error) {
	query := machineSelect + `
		WHERE m.id = $1 AND m.is_deleted = false
			AND ($2::integer IS NULL OR m.user_id = $2)
			AND ($3::integer IS NULL OR EXISTS (
				SELECT 1 FROM user_machine_access a WHERE a.machine_id = m.id AND a.user_id = $3))`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, id, scope.OwnerID, scope.AssignedUserID)
	if err != nil {
		return model.Machine{}, fmt.Errorf("find machine: %w", err)
	}
	machine, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Machine])
	return machine, mapMachineError("find machine", err)
}

// CountOwnedBy returns how many of ids are machines of customer ownerID.
func (machineRepository *MachineRepository) CountOwnedBy(ctx context.Context, ownerID int64, ids []int64) (int, error) {
	query := `SELECT COUNT(*) FROM machines
		WHERE user_id = $1 AND id = ANY($2::bigint[]) AND is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, ownerID, ids)
	if err != nil {
		return 0, fmt.Errorf("count customer machines: %w", err)
	}
	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, fmt.Errorf("count customer machines: %w", err)
	}
	return count, nil
}

// Create inserts the form fields of machine and returns the new id. Runtime columns
// (online status, mode, …) keep their defaults.
func (machineRepository *MachineRepository) Create(ctx context.Context, machine model.Machine) (int64, error) {
	query := `INSERT INTO machines (user_id, model_id, name, serial_number, brand_id, controller_id,
			controller_version, accuracy_id, line_id, is_insured, maintenance_required,
			installation_date, bill_date, warranty_date, hourly_cost, amc, machine_system_type, security_mode,
			local_ip, live_ip, port, username, password, controller_version_type,
			longitude, latitude, health_ids, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12::date, $13::date, $14::date, $15::numeric,
			NULLIF($16, ''), $17, $18, $19, $20, $21, $22, $23, $24,
			NULLIF($25, ''), NULLIF($26, ''), NULLIF($27, ''), $28)
		RETURNING id`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.UserID, machine.ModelID, machine.Name, machine.SerialNumber, machine.BrandID, machine.ControllerID,
		machine.ControllerVersion, machine.AccuracyID, machine.LineID, machine.IsInsured, machine.MaintenanceRequired,
		machine.InstallationDate, machine.BillDate, machine.WarrantyDate, machine.HourlyCost, machine.AMC,
		machine.MachineSystemType, machine.SecurityMode,
		machine.LocalIP, machine.LiveIP, machine.Port, machine.Username, machine.Password, machine.ControllerVersionType,
		machine.Longitude, machine.Latitude, machine.HealthIDs, machine.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("create machine: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapMachineError("create machine", err)
}

// Update replaces the form fields of machine.ID. A blank Password or Status keeps the
// saved value; runtime columns are left as they are.
func (machineRepository *MachineRepository) Update(ctx context.Context, machine model.Machine) error {
	query := `UPDATE machines
		SET user_id = $2, model_id = $3, name = $4, serial_number = $5, brand_id = $6, controller_id = $7,
			controller_version = $8, accuracy_id = $9, line_id = $10, is_insured = $11, maintenance_required = $12,
			installation_date = $13::date, bill_date = $14::date, warranty_date = $15::date,
			hourly_cost = $16::numeric, amc = NULLIF($17, ''), machine_system_type = $18, security_mode = $19,
			local_ip = $20, live_ip = $21, port = $22, username = $23,
			password = COALESCE(NULLIF($24, ''), password), controller_version_type = $25,
			longitude = NULLIF($26, ''), latitude = NULLIF($27, ''), health_ids = NULLIF($28, ''),
			status = COALESCE(NULLIF($29, ''), status)
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.ID, machine.UserID, machine.ModelID, machine.Name, machine.SerialNumber, machine.BrandID,
		machine.ControllerID, machine.ControllerVersion, machine.AccuracyID, machine.LineID, machine.IsInsured,
		machine.MaintenanceRequired, machine.InstallationDate, machine.BillDate, machine.WarrantyDate,
		machine.HourlyCost, machine.AMC, machine.MachineSystemType, machine.SecurityMode,
		machine.LocalIP, machine.LiveIP, machine.Port, machine.Username,
		machine.Password, machine.ControllerVersionType,
		machine.Longitude, machine.Latitude, machine.HealthIDs, machine.Status,
	)
	if err != nil {
		return fmt.Errorf("update machine: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return mapMachineError("update machine", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "machines_user_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "machines_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		case pgError.Code == "23503" && pgError.ConstraintName == "machines_model_id_fkey":
			return ErrInvalidModel
		case pgError.Code == "23503" && pgError.ConstraintName == "machines_brand_id_fkey":
			return ErrInvalidBrand
		case pgError.Code == "23503" && pgError.ConstraintName == "machines_controller_id_fkey":
			return ErrInvalidController
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
