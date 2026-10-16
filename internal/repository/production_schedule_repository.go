package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// productionScheduleSelect reads model.ProductionSchedule: the production_schedules row
// (ps) plus the display names of its customer and machine and its item totals.
const productionScheduleSelect = `SELECT ps.id, ps.user_id, ps.machine_id, ps.schedule_type,
		ps.start_date::text AS start_date, ps.end_date::text AS end_date,
		ps.status, ps.is_deleted, ps.deleted_at, ps.created_at, ps.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name,
		COALESCE(m.name, '') AS machine_name,
		COALESCE(items.item_count, 0)::integer AS item_count,
		COALESCE(items.total_quantity, 0)::integer AS total_quantity
	FROM production_schedules ps
	LEFT JOIN users u ON u.id = ps.user_id
	LEFT JOIN machines m ON m.id = ps.machine_id
	LEFT JOIN company_profiles cp ON cp.user_id = ps.user_id AND cp.is_deleted = false
	LEFT JOIN LATERAL (
		SELECT COUNT(*) AS item_count, SUM(psi.quantity) AS total_quantity
		FROM production_schedule_items psi
		WHERE psi.schedule_id = ps.id AND psi.is_deleted = false
	) items ON true`

// userDisplayName is the display name of users row alias, as in the customer names.
const userDisplayName = `COALESCE(NULLIF(TRIM(CONCAT_WS(' ', %[1]s.first_name, %[1]s.last_name)), ''), %[1]s.username, '')`

// ProductionScheduleRepository runs the SQL queries for the production_schedules and
// production_schedule_items tables. Soft-deleted rows (is_deleted = true) are invisible
// to every method.
//
// A nil ownerID does not narrow a query; otherwise only the schedules of that customer
// (production_schedules.user_id) are found.
type ProductionScheduleRepository struct {
	database db.Querier
}

// NewProductionScheduleRepository creates a ProductionScheduleRepository that runs its
// queries on database (a connection pool or a transaction).
func NewProductionScheduleRepository(database db.Querier) *ProductionScheduleRepository {
	return &ProductionScheduleRepository{database: database}
}

// FindAll returns the schedules of ownerID (every schedule when nil), latest date first.
func (scheduleRepository *ProductionScheduleRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.ProductionSchedule, error) {
	query := productionScheduleSelect + `
		WHERE ps.is_deleted = false AND ($1::integer IS NULL OR ps.user_id = $1)
		ORDER BY ps.start_date DESC, ps.id DESC`

	rows, err := db.Execute(ctx, scheduleRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all production schedules: %w", err)
	}
	schedules, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ProductionSchedule])
	if err != nil {
		return nil, fmt.Errorf("find all production schedules: %w", err)
	}
	return schedules, nil
}

// FindByID returns the schedule, or ErrNotFound when it is not a schedule of ownerID.
func (scheduleRepository *ProductionScheduleRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.ProductionSchedule, error) {
	query := productionScheduleSelect + `
		WHERE ps.id = $1 AND ps.is_deleted = false AND ($2::integer IS NULL OR ps.user_id = $2)`

	rows, err := db.Execute(ctx, scheduleRepository.database, query, id, ownerID)
	if err != nil {
		return model.ProductionSchedule{}, fmt.Errorf("find production schedule: %w", err)
	}
	schedule, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ProductionSchedule])
	return schedule, mapProductionScheduleError("find production schedule", err)
}

// FindItems returns the items of scheduleID in form order, with display names.
func (scheduleRepository *ProductionScheduleRepository) FindItems(ctx context.Context, scheduleID int64) ([]model.ProductionScheduleItem, error) {
	query := `SELECT psi.id, psi.schedule_id, psi.sort_order, psi.schedule_date::text AS schedule_date,
			psi.job_id, psi.supervisor_id, psi.operator_id, psi.shift_id, psi.part_id, psi.process_id,
			psi.quantity, to_char(psi.start_time, 'HH24:MI') AS start_time,
			to_char(psi.end_time, 'HH24:MI') AS end_time, psi.total_minutes,
			psi.is_deleted, psi.deleted_at, psi.created_at, psi.updated_at,
			` + fmt.Sprintf(userDisplayName, "su") + ` AS supervisor_name,
			` + fmt.Sprintf(userDisplayName, "op") + ` AS operator_name,
			COALESCE(s.name, '') AS shift_name,
			COALESCE(CONCAT_WS(' - ', p.part_no, p.part_name), '') AS part_name,
			COALESCE(pp.process, '') AS process_name
		FROM production_schedule_items psi
		LEFT JOIN users su ON su.id = psi.supervisor_id
		LEFT JOIN users op ON op.id = psi.operator_id
		LEFT JOIN shifts s ON s.id = psi.shift_id
		LEFT JOIN parts p ON p.id = psi.part_id
		LEFT JOIN part_processes pp ON pp.id = psi.process_id
		WHERE psi.schedule_id = $1 AND psi.is_deleted = false
		ORDER BY psi.sort_order, psi.id`

	rows, err := db.Execute(ctx, scheduleRepository.database, query, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("find production schedule items: %w", err)
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ProductionScheduleItem])
	if err != nil {
		return nil, fmt.Errorf("find production schedule items: %w", err)
	}
	return items, nil
}

// MachineBooked reports whether another schedule (not excludeID) of machineID shares a
// day with startDate … endDate.
func (scheduleRepository *ProductionScheduleRepository) MachineBooked(ctx context.Context, machineID int64, startDate, endDate string, excludeID int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM production_schedules
			WHERE machine_id = $1 AND is_deleted = false AND id <> $4
				AND start_date <= $3::date AND end_date >= $2::date
		)`

	rows, err := db.Execute(ctx, scheduleRepository.database, query, machineID, startDate, endDate, excludeID)
	if err != nil {
		return false, fmt.Errorf("check machine schedule: %w", err)
	}
	booked, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check machine schedule: %w", err)
	}
	return booked, nil
}

// CountCompanyMembers returns how many of ids are members of customer companyID with role.
func (scheduleRepository *ProductionScheduleRepository) CountCompanyMembers(ctx context.Context, companyID int64, role string, ids []int64) (int, error) {
	query := `SELECT COUNT(*) FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.company_id = $1 AND r.name = $2 AND u.id = ANY($3::bigint[]) AND u.is_deleted = false`
	return scheduleRepository.count(ctx, "count company members", query, companyID, role, ids)
}

// CountCustomerParts returns how many of ids are parts of customer ownerID.
func (scheduleRepository *ProductionScheduleRepository) CountCustomerParts(ctx context.Context, ownerID int64, ids []int64) (int, error) {
	query := `SELECT COUNT(*) FROM parts
		WHERE user_id = $1 AND id = ANY($2::bigint[]) AND is_deleted = false`
	return scheduleRepository.count(ctx, "count customer parts", query, ownerID, ids)
}

// CountMachineProcesses returns how many of the (partIDs[i], processIDs[i]) pairs are a
// process of that part with machineID among its applicable machines.
func (scheduleRepository *ProductionScheduleRepository) CountMachineProcesses(ctx context.Context, machineID int64, partIDs, processIDs []int64) (int, error) {
	query := `SELECT COUNT(DISTINCT pair.process_id) FROM unnest($2::bigint[], $3::bigint[]) AS pair(part_id, process_id)
		JOIN part_processes pp ON pp.id = pair.process_id AND pp.part_id = pair.part_id AND pp.is_deleted = false
		JOIN part_process_machines ppm ON ppm.process_id = pp.id AND ppm.machine_id = $1`
	return scheduleRepository.count(ctx, "count machine processes", query, machineID, partIDs, processIDs)
}

func (scheduleRepository *ProductionScheduleRepository) count(ctx context.Context, operation, query string, args ...any) (int, error) {
	rows, err := db.Execute(ctx, scheduleRepository.database, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	return count, nil
}

// Create inserts schedule and returns the new id.
func (scheduleRepository *ProductionScheduleRepository) Create(ctx context.Context, schedule model.ProductionSchedule) (int64, error) {
	query := `INSERT INTO production_schedules (user_id, machine_id, schedule_type, start_date, end_date, status)
		VALUES ($1, $2, $3, $4::date, $5::date, $6)
		RETURNING id`

	rows, err := db.Execute(ctx, scheduleRepository.database, query,
		schedule.UserID, schedule.MachineID, schedule.ScheduleType, schedule.StartDate, schedule.EndDate, schedule.Status)
	if err != nil {
		return 0, fmt.Errorf("create production schedule: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapProductionScheduleError("create production schedule", err)
}

// Update replaces the fields of schedule.ID. A blank Status keeps the saved status.
func (scheduleRepository *ProductionScheduleRepository) Update(ctx context.Context, schedule model.ProductionSchedule) error {
	query := `UPDATE production_schedules
		SET user_id = $2, machine_id = $3, schedule_type = $4, start_date = $5::date, end_date = $6::date,
			status = COALESCE(NULLIF($7, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, scheduleRepository.database, query,
		schedule.ID, schedule.UserID, schedule.MachineID, schedule.ScheduleType, schedule.StartDate, schedule.EndDate, schedule.Status)
	if err != nil {
		return mapProductionScheduleError("update production schedule", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the schedule of ownerID (any schedule when nil). Its items stay as
// they are: they are only ever read through the schedule.
func (scheduleRepository *ProductionScheduleRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE production_schedules
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, scheduleRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete production schedule: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateItem inserts item.
func (scheduleRepository *ProductionScheduleRepository) CreateItem(ctx context.Context, item model.ProductionScheduleItem) error {
	query := `INSERT INTO production_schedule_items (schedule_id, sort_order, schedule_date, job_id,
			supervisor_id, operator_id, shift_id, part_id, process_id, quantity, start_time, end_time, total_minutes)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7, $8, $9, $10, $11::time, $12::time, $13)`

	if _, err := executeAffected(ctx, scheduleRepository.database, query,
		item.ScheduleID, item.SortOrder, item.ScheduleDate, item.JobID, item.SupervisorID, item.OperatorID,
		item.ShiftID, item.PartID, item.ProcessID, item.Quantity, item.StartTime, item.EndTime, item.TotalMinutes); err != nil {
		return fmt.Errorf("create production schedule item: %w", err)
	}
	return nil
}

// UpdateItem replaces the fields of item.ID, or returns ErrNotFound when it is not an
// item of item.ScheduleID.
func (scheduleRepository *ProductionScheduleRepository) UpdateItem(ctx context.Context, item model.ProductionScheduleItem) error {
	query := `UPDATE production_schedule_items
		SET sort_order = $3, schedule_date = $4::date, job_id = $5, supervisor_id = $6, operator_id = $7,
			shift_id = $8, part_id = $9, process_id = $10, quantity = $11, start_time = $12::time,
			end_time = $13::time, total_minutes = $14
		WHERE id = $1 AND schedule_id = $2 AND is_deleted = false`

	affected, err := executeAffected(ctx, scheduleRepository.database, query,
		item.ID, item.ScheduleID, item.SortOrder, item.ScheduleDate, item.JobID, item.SupervisorID, item.OperatorID,
		item.ShiftID, item.PartID, item.ProcessID, item.Quantity, item.StartTime, item.EndTime, item.TotalMinutes)
	if err != nil {
		return fmt.Errorf("update production schedule item: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteItemsExcept soft-deletes the items of scheduleID whose id is not in keepIDs.
func (scheduleRepository *ProductionScheduleRepository) DeleteItemsExcept(ctx context.Context, scheduleID int64, keepIDs []int64) error {
	if keepIDs == nil {
		keepIDs = []int64{}
	}
	query := `UPDATE production_schedule_items
		SET is_deleted = true, deleted_at = NOW()
		WHERE schedule_id = $1 AND is_deleted = false AND id <> ALL($2::bigint[])`

	if _, err := executeAffected(ctx, scheduleRepository.database, query, scheduleID, keepIDs); err != nil {
		return fmt.Errorf("delete production schedule items: %w", err)
	}
	return nil
}

// mapProductionScheduleError turns driver errors into repository errors. A nil err stays nil.
func mapProductionScheduleError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23503" && pgError.ConstraintName == "production_schedules_user_id_fkey" {
		return ErrInvalidUser
	}
	return fmt.Errorf("%s: %w", operation, err)
}
