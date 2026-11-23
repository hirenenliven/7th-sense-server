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
// its customer, model, brand, controller and organization path (op: the level the machine
// is saved at and every level above it, migration 040; line_id is the path's line) and
// its customer's filters. DATE and NUMERIC columns are cast to text because the model
// holds them as strings. Soft-deleted joined rows still give their name, so a machine
// keeps showing what it was saved with.
const machineSelect = `SELECT m.id, m.user_id, m.machine_type_id, m.model_id, m.name, m.serial_number,
		m.live_ip, m.port, m.username, m.password, m.password <> '' AS has_password, m.local_ip,
		m.brand_id, m.controller_id, m.controller_version, m.accuracy_id, op.user_line_id AS line_id,
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
		COALESCE(op.user_line_name, '') AS line_name,
		op.user_country_id, COALESCE(op.user_country_name, '') AS user_country_name,
		op.user_state_id, COALESCE(op.user_state_name, '') AS user_state_name,
		op.user_city_id, COALESCE(op.user_city_name, '') AS user_city_name,
		op.branch_id, COALESCE(op.branch_name, '') AS branch_name,
		op.user_location_id, COALESCE(op.user_location_name, '') AS user_location_name,
		op.user_shop_id, COALESCE(op.user_shop_name, '') AS user_shop_name,
		op.user_cell_id, COALESCE(op.user_cell_name, '') AS user_cell_name,
		op.user_supervisor_id, COALESCE(op.user_supervisor_name, '') AS user_supervisor_name,
		COALESCE(up.filters, '{}') AS customer_filters
	FROM machines m
	LEFT JOIN users u ON u.id = m.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = m.user_id AND cp.is_deleted = false
	LEFT JOIN user_profiles up ON up.user_id = m.user_id AND up.is_deleted = false
	LEFT JOIN models mo ON mo.id = m.model_id
	LEFT JOIN brands b ON b.id = m.brand_id
	LEFT JOIN controllers c ON c.id = m.controller_id
	LEFT JOIN machine_organization_paths op ON op.machine_id = m.id`

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

// accessibleMachine is the condition "user <parameter> may access machine m" (NULL: every
// machine): m is assigned to the user (user_machine_access), or the user's machine filter has
// all machines on and m's path lies under every level the filter names (machine_organization_paths,
// so a machine saved at a cell or above its line counts too; a filter without levels takes every
// machine), so machines added later are included.
func accessibleMachine(parameter string) string {
	return `(` + parameter + `::integer IS NULL
		OR EXISTS (SELECT 1 FROM user_machine_access a WHERE a.machine_id = m.id AND a.user_id = ` + parameter + `)
		OR EXISTS (
			SELECT 1 FROM user_machine_filters f
			WHERE f.user_id = ` + parameter + ` AND f.all_machines
				AND (num_nonnulls(f.user_country_id, f.user_state_id, f.user_city_id, f.branch_id,
						f.user_location_id, f.user_shop_id, f.user_line_id) = 0
					OR EXISTS (
						SELECT 1 FROM machine_organization_paths p
						WHERE p.machine_id = m.id
							AND (f.user_country_id IS NULL OR p.user_country_id = f.user_country_id)
							AND (f.user_state_id IS NULL OR p.user_state_id = f.user_state_id)
							AND (f.user_city_id IS NULL OR p.user_city_id = f.user_city_id)
							AND (f.branch_id IS NULL OR p.branch_id = f.branch_id)
							AND (f.user_location_id IS NULL OR p.user_location_id = f.user_location_id)
							AND (f.user_shop_id IS NULL OR p.user_shop_id = f.user_shop_id)
							AND (f.user_line_id IS NULL OR p.user_line_id = f.user_line_id)))))`
}

// FindAll returns the machines in scope, newest first.
func (machineRepository *MachineRepository) FindAll(ctx context.Context, scope model.MachineScope) ([]model.Machine, error) {
	query := machineSelect + `
		WHERE m.is_deleted = false
			AND ($1::integer IS NULL OR m.user_id = $1)
			AND ` + accessibleMachine("$2") + `
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
			AND ` + accessibleMachine("$3")

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

// Create inserts the form fields of machine and returns the new id. place is the one level
// the machine is saved at (model.OrganizationPath.Deepest; empty = none). Runtime columns
// (online status, mode, …) keep their defaults.
func (machineRepository *MachineRepository) Create(ctx context.Context, machine model.Machine, place model.OrganizationPath) (int64, error) {
	query := `INSERT INTO machines (user_id, model_id, name, serial_number, brand_id, controller_id,
			controller_version, accuracy_id, line_id, is_insured, maintenance_required,
			installation_date, bill_date, warranty_date, hourly_cost, amc, machine_system_type, security_mode,
			local_ip, live_ip, port, username, password, controller_version_type,
			longitude, latitude, health_ids, status,
			user_country_id, user_state_id, user_city_id, branch_id, user_location_id, user_shop_id,
			user_cell_id, user_supervisor_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12::date, $13::date, $14::date, $15::numeric,
			NULLIF($16, ''), $17, $18, $19, $20, $21, $22, $23, $24,
			NULLIF($25, ''), NULLIF($26, ''), NULLIF($27, ''), $28,
			$29, $30, $31, $32, $33, $34, $35, $36)
		RETURNING id`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.UserID, machine.ModelID, machine.Name, machine.SerialNumber, machine.BrandID, machine.ControllerID,
		machine.ControllerVersion, machine.AccuracyID, place.UserLineID, machine.IsInsured, machine.MaintenanceRequired,
		machine.InstallationDate, machine.BillDate, machine.WarrantyDate, machine.HourlyCost, machine.AMC,
		machine.MachineSystemType, machine.SecurityMode,
		machine.LocalIP, machine.LiveIP, machine.Port, machine.Username, machine.Password, machine.ControllerVersionType,
		machine.Longitude, machine.Latitude, machine.HealthIDs, machine.Status,
		place.UserCountryID, place.UserStateID, place.UserCityID, place.BranchID, place.UserLocationID,
		place.UserShopID, place.UserCellID, place.UserSupervisorID,
	)
	if err != nil {
		return 0, fmt.Errorf("create machine: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapMachineError("create machine", err)
}

// Update replaces the form fields of machine.ID and saves it at place (one level, as for
// Create). A blank Password or Status keeps the saved value; runtime columns are left as
// they are.
func (machineRepository *MachineRepository) Update(ctx context.Context, machine model.Machine, place model.OrganizationPath) error {
	query := `UPDATE machines
		SET user_id = $2, model_id = $3, name = $4, serial_number = $5, brand_id = $6, controller_id = $7,
			controller_version = $8, accuracy_id = $9, line_id = $10, is_insured = $11, maintenance_required = $12,
			installation_date = $13::date, bill_date = $14::date, warranty_date = $15::date,
			hourly_cost = $16::numeric, amc = NULLIF($17, ''), machine_system_type = $18, security_mode = $19,
			local_ip = $20, live_ip = $21, port = $22, username = $23,
			password = COALESCE(NULLIF($24, ''), password), controller_version_type = $25,
			longitude = NULLIF($26, ''), latitude = NULLIF($27, ''), health_ids = NULLIF($28, ''),
			status = COALESCE(NULLIF($29, ''), status),
			user_country_id = $30, user_state_id = $31, user_city_id = $32, branch_id = $33,
			user_location_id = $34, user_shop_id = $35, user_cell_id = $36, user_supervisor_id = $37
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query,
		machine.ID, machine.UserID, machine.ModelID, machine.Name, machine.SerialNumber, machine.BrandID,
		machine.ControllerID, machine.ControllerVersion, machine.AccuracyID, place.UserLineID, machine.IsInsured,
		machine.MaintenanceRequired, machine.InstallationDate, machine.BillDate, machine.WarrantyDate,
		machine.HourlyCost, machine.AMC, machine.MachineSystemType, machine.SecurityMode,
		machine.LocalIP, machine.LiveIP, machine.Port, machine.Username,
		machine.Password, machine.ControllerVersionType,
		machine.Longitude, machine.Latitude, machine.HealthIDs, machine.Status,
		place.UserCountryID, place.UserStateID, place.UserCityID, place.BranchID,
		place.UserLocationID, place.UserShopID, place.UserCellID, place.UserSupervisorID,
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

// FindCustomerFilters returns the organization levels customerID uses (user_profiles.filters,
// as saved; empty when none or no profile).
func (machineRepository *MachineRepository) FindCustomerFilters(ctx context.Context, customerID int64) ([]string, error) {
	query := `SELECT COALESCE((SELECT filters FROM user_profiles WHERE user_id = $1 AND is_deleted = false), '{}')`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("find customer filters: %w", err)
	}
	filters, err := pgx.CollectOneRow(rows, pgx.RowTo[[]string])
	if err != nil {
		return nil, fmt.Errorf("find customer filters: %w", err)
	}
	return filters, nil
}

// IsPathInOrganization reports whether path is one chain of customerID's organization: walked up
// from its lowest level, every level it names is the row found there (the customer's, not
// deleted) and the walk reaches one of the customer's countries. Branches are shared rows: a
// branch counts when its city is the customer's.
func (machineRepository *MachineRepository) IsPathInOrganization(ctx context.Context, customerID int64, path model.OrganizationPath) (bool, error) {
	query := `WITH picked AS (
			SELECT $2::integer AS country, $3::integer AS state, $4::integer AS city, $5::integer AS branch,
				$6::integer AS location, $7::integer AS shop, $8::integer AS line, $9::integer AS cell,
				$10::integer AS supervisor
		)
		SELECT COALESCE(co.id IS NOT NULL
			AND (p.country IS NULL OR p.country = co.id)
			AND (p.state IS NULL OR p.state = st.id)
			AND (p.city IS NULL OR p.city = ci.id)
			AND (p.branch IS NULL OR p.branch = br.id)
			AND (p.location IS NULL OR p.location = lo.id)
			AND (p.shop IS NULL OR p.shop = sh.id)
			AND (p.line IS NULL OR p.line = li.id)
			AND (p.cell IS NULL OR p.cell = ce.id)
			AND (p.supervisor IS NULL OR p.supervisor = su.id), false)
		FROM picked p
		LEFT JOIN user_supervisor su ON su.id = p.supervisor AND su.user_id = $1 AND su.is_deleted = false
		LEFT JOIN user_cell ce ON ce.id = COALESCE(su.user_cell_id, p.cell) AND ce.user_id = $1 AND ce.is_deleted = false
		LEFT JOIN user_line li ON li.id = COALESCE(ce.user_line_id, p.line) AND li.user_id = $1 AND li.is_deleted = false
		LEFT JOIN user_shop sh ON sh.id = COALESCE(li.user_shop_id, p.shop) AND sh.user_id = $1 AND sh.is_deleted = false
		LEFT JOIN user_location lo ON lo.id = COALESCE(sh.user_location_id, p.location) AND lo.user_id = $1 AND lo.is_deleted = false
		LEFT JOIN branches br ON br.id = COALESCE(lo.branch_id, p.branch) AND br.is_deleted = false
		LEFT JOIN user_city ci ON ci.id = COALESCE(br.user_city_id, p.city) AND ci.user_id = $1 AND ci.is_deleted = false
		LEFT JOIN user_state st ON st.id = COALESCE(ci.user_state_id, p.state) AND st.user_id = $1 AND st.is_deleted = false
		LEFT JOIN user_country co ON co.id = COALESCE(st.user_country_id, p.country) AND co.user_id = $1 AND co.is_deleted = false`

	rows, err := db.Execute(ctx, machineRepository.dbPool, query, customerID,
		path.UserCountryID, path.UserStateID, path.UserCityID, path.BranchID, path.UserLocationID,
		path.UserShopID, path.UserLineID, path.UserCellID, path.UserSupervisorID)
	if err != nil {
		return false, fmt.Errorf("check machine organization: %w", err)
	}
	valid, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check machine organization: %w", err)
	}
	return valid, nil
}

// levelHasMachines reports whether machines that are not deleted are saved at the organization
// level id in column (line_id or one of the level columns of migration 040).
func levelHasMachines(ctx context.Context, database db.Querier, column string, id int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM machines WHERE ` + column + ` = $1 AND is_deleted = false)`

	rows, err := db.Execute(ctx, database, query, id)
	if err != nil {
		return false, fmt.Errorf("check level machines: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check level machines: %w", err)
	}
	return found, nil
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
