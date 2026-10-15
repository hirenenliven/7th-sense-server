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

// partSelect reads model.Part: the parts row (p) plus the display names of its customer
// and clientele and the number of its processes.
const partSelect = `SELECT p.id, p.user_id, p.clientele_id, p.part_no, p.part_name, p.part_count,
		p.production_time, p.status, p.is_deleted, p.deleted_at, p.created_at, p.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name,
		(SELECT COUNT(*) FROM part_processes pp WHERE pp.part_id = p.id AND pp.is_deleted = false)::integer AS process_count,
		COALESCE(cl.name, '') AS clientele_name
	FROM parts p
	LEFT JOIN users u ON u.id = p.user_id
	LEFT JOIN clienteles cl ON cl.id = p.clientele_id
	LEFT JOIN company_profiles cp ON cp.user_id = p.user_id AND cp.is_deleted = false`

// PartRepository runs the SQL queries for the parts, part_processes,
// part_process_machines and part_sub_programs tables. Soft-deleted rows
// (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the parts of that customer
// (parts.user_id) are found.
type PartRepository struct {
	database db.Querier
}

// NewPartRepository creates a PartRepository that runs its queries on database
// (a connection pool or a transaction).
func NewPartRepository(database db.Querier) *PartRepository {
	return &PartRepository{database: database}
}

// FindAll returns the parts of ownerID (every part when nil), newest first.
func (partRepository *PartRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.Part, error) {
	query := partSelect + `
		WHERE p.is_deleted = false AND ($1::integer IS NULL OR p.user_id = $1)
		ORDER BY p.id DESC`

	rows, err := db.Execute(ctx, partRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all parts: %w", err)
	}
	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Part])
	if err != nil {
		return nil, fmt.Errorf("find all parts: %w", err)
	}
	return parts, nil
}

// FindByID returns the part, or ErrNotFound when it is not a part of ownerID.
func (partRepository *PartRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.Part, error) {
	query := partSelect + `
		WHERE p.id = $1 AND p.is_deleted = false AND ($2::integer IS NULL OR p.user_id = $2)`

	rows, err := db.Execute(ctx, partRepository.database, query, id, ownerID)
	if err != nil {
		return model.Part{}, fmt.Errorf("find part: %w", err)
	}
	part, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Part])
	return part, mapPartError("find part", err)
}

// FindProcesses returns the processes of partID in form order, each with the ids of
// its applicable machines (soft-deleted machines left out).
func (partRepository *PartRepository) FindProcesses(ctx context.Context, partID int64) ([]model.PartProcess, error) {
	query := `SELECT pp.id, pp.part_id, pp.sort_order, pp.process, pp.program_name, pp.drawing_no,
			pp.drawing_revision_no, pp.drawing_file, pp.drawing_name, pp.accuracy_id, pp.setup_time,
			pp.cycle_time, pp.loading_time, pp.production_time, pp.has_sub_program,
			ARRAY(SELECT ppm.machine_id::bigint FROM part_process_machines ppm
				JOIN machines m ON m.id = ppm.machine_id AND m.is_deleted = false
				WHERE ppm.process_id = pp.id ORDER BY ppm.machine_id) AS machine_ids,
			pp.is_deleted, pp.deleted_at, pp.created_at, pp.updated_at
		FROM part_processes pp
		WHERE pp.part_id = $1 AND pp.is_deleted = false
		ORDER BY pp.sort_order, pp.id`

	rows, err := db.Execute(ctx, partRepository.database, query, partID)
	if err != nil {
		return nil, fmt.Errorf("find part processes: %w", err)
	}
	processes, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.PartProcess])
	if err != nil {
		return nil, fmt.Errorf("find part processes: %w", err)
	}
	return processes, nil
}

// FindSubPrograms returns the sub programs of every process of partID, in form order.
func (partRepository *PartRepository) FindSubPrograms(ctx context.Context, partID int64) ([]model.PartSubProgram, error) {
	query := `SELECT sp.id, sp.process_id, sp.sort_order, sp.name, sp.cycle_time,
			sp.is_deleted, sp.deleted_at, sp.created_at, sp.updated_at
		FROM part_sub_programs sp
		JOIN part_processes pp ON pp.id = sp.process_id AND pp.is_deleted = false
		WHERE pp.part_id = $1 AND sp.is_deleted = false
		ORDER BY sp.process_id, sp.sort_order, sp.id`

	rows, err := db.Execute(ctx, partRepository.database, query, partID)
	if err != nil {
		return nil, fmt.Errorf("find part sub programs: %w", err)
	}
	subPrograms, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.PartSubProgram])
	if err != nil {
		return nil, fmt.Errorf("find part sub programs: %w", err)
	}
	return subPrograms, nil
}

// Create inserts part and returns the new id.
func (partRepository *PartRepository) Create(ctx context.Context, part model.Part) (int64, error) {
	query := `INSERT INTO parts (user_id, clientele_id, part_no, part_name, part_count, production_time, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	rows, err := db.Execute(ctx, partRepository.database, query,
		part.UserID, part.ClienteleID, part.PartNo, part.PartName, part.PartCount, part.ProductionTime, part.Status)
	if err != nil {
		return 0, fmt.Errorf("create part: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapPartError("create part", err)
}

// Update replaces the fields of part.ID. A blank Status keeps the saved status.
func (partRepository *PartRepository) Update(ctx context.Context, part model.Part) error {
	query := `UPDATE parts
		SET user_id = $2, clientele_id = $3, part_no = $4, part_name = $5, part_count = $6,
			production_time = $7, status = COALESCE(NULLIF($8, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, partRepository.database, query,
		part.ID, part.UserID, part.ClienteleID, part.PartNo, part.PartName, part.PartCount, part.ProductionTime, part.Status)
	if err != nil {
		return mapPartError("update part", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the part of ownerID (any part when nil). Its processes stay as
// they are: they are only ever read through the part.
func (partRepository *PartRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE parts
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, partRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete part: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateProcess inserts process and returns the new id.
func (partRepository *PartRepository) CreateProcess(ctx context.Context, process model.PartProcess) (int64, error) {
	query := `INSERT INTO part_processes (part_id, sort_order, process, program_name, drawing_no,
			drawing_revision_no, drawing_file, drawing_name, accuracy_id, setup_time, cycle_time,
			loading_time, production_time, has_sub_program)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	rows, err := db.Execute(ctx, partRepository.database, query,
		process.PartID, process.SortOrder, process.Process, process.ProgramName, process.DrawingNo,
		process.DrawingRevisionNo, process.DrawingFile, process.DrawingName, process.AccuracyID, process.SetupTime,
		process.CycleTime, process.LoadingTime, process.ProductionTime, process.HasSubProgram)
	if err != nil {
		return 0, fmt.Errorf("create part process: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return 0, fmt.Errorf("create part process: %w", err)
	}
	return id, nil
}

// UpdateProcess replaces the fields of process.ID, or returns ErrNotFound when it is
// not a process of process.PartID.
func (partRepository *PartRepository) UpdateProcess(ctx context.Context, process model.PartProcess) error {
	query := `UPDATE part_processes
		SET sort_order = $3, process = $4, program_name = $5, drawing_no = $6, drawing_revision_no = $7,
			drawing_file = $8, drawing_name = $9, accuracy_id = $10, setup_time = $11, cycle_time = $12,
			loading_time = $13, production_time = $14, has_sub_program = $15
		WHERE id = $1 AND part_id = $2 AND is_deleted = false`

	affected, err := executeAffected(ctx, partRepository.database, query,
		process.ID, process.PartID, process.SortOrder, process.Process, process.ProgramName, process.DrawingNo,
		process.DrawingRevisionNo, process.DrawingFile, process.DrawingName, process.AccuracyID, process.SetupTime,
		process.CycleTime, process.LoadingTime, process.ProductionTime, process.HasSubProgram)
	if err != nil {
		return fmt.Errorf("update part process: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteProcessesExcept soft-deletes the processes of partID whose id is not in keepIDs.
func (partRepository *PartRepository) DeleteProcessesExcept(ctx context.Context, partID int64, keepIDs []int64) error {
	if keepIDs == nil {
		keepIDs = []int64{}
	}
	query := `UPDATE part_processes
		SET is_deleted = true, deleted_at = NOW()
		WHERE part_id = $1 AND is_deleted = false AND id <> ALL($2::bigint[])`

	if _, err := executeAffected(ctx, partRepository.database, query, partID, keepIDs); err != nil {
		return fmt.Errorf("delete part processes: %w", err)
	}
	return nil
}

// ReplaceProcessMachines makes machineIDs the applicable machines of processID:
// machines not in the list are removed and missing ones are added.
func (partRepository *PartRepository) ReplaceProcessMachines(ctx context.Context, processID int64, machineIDs []int64) error {
	if machineIDs == nil {
		machineIDs = []int64{}
	}
	deleteQuery := `DELETE FROM part_process_machines
		WHERE process_id = $1 AND machine_id <> ALL($2::bigint[])`
	if _, err := executeAffected(ctx, partRepository.database, deleteQuery, processID, machineIDs); err != nil {
		return fmt.Errorf("remove part process machines: %w", err)
	}

	insertQuery := `INSERT INTO part_process_machines (process_id, machine_id)
		SELECT $1, machine_id FROM unnest($2::bigint[]) AS machine_id
		ON CONFLICT (process_id, machine_id) DO NOTHING`
	if _, err := executeAffected(ctx, partRepository.database, insertQuery, processID, machineIDs); err != nil {
		return fmt.Errorf("add part process machines: %w", err)
	}
	return nil
}

// ReplaceSubPrograms makes subPrograms the sub programs of processID: the saved ones
// are soft-deleted and subPrograms inserted in their order.
func (partRepository *PartRepository) ReplaceSubPrograms(ctx context.Context, processID int64, subPrograms []model.PartSubProgram) error {
	deleteQuery := `UPDATE part_sub_programs
		SET is_deleted = true, deleted_at = NOW()
		WHERE process_id = $1 AND is_deleted = false`
	if _, err := executeAffected(ctx, partRepository.database, deleteQuery, processID); err != nil {
		return fmt.Errorf("delete part sub programs: %w", err)
	}
	if len(subPrograms) == 0 {
		return nil
	}

	names := make([]string, len(subPrograms))
	cycleTimes := make([]int32, len(subPrograms))
	for index, subProgram := range subPrograms {
		names[index] = subProgram.Name
		cycleTimes[index] = int32(subProgram.CycleTime)
	}
	insertQuery := `INSERT INTO part_sub_programs (process_id, sort_order, name, cycle_time)
		SELECT $1, sub.ordinality, sub.name, sub.cycle_time
		FROM unnest($2::text[], $3::integer[]) WITH ORDINALITY AS sub(name, cycle_time, ordinality)`
	if _, err := executeAffected(ctx, partRepository.database, insertQuery, processID, names, cycleTimes); err != nil {
		return fmt.Errorf("add part sub programs: %w", err)
	}
	return nil
}

// executeAffected runs a query whose rows are not needed and returns how many rows it changed.
func executeAffected(ctx context.Context, database db.Querier, query string, args ...any) (int64, error) {
	rows, err := db.Execute(ctx, database, query, args...)
	if err != nil {
		return 0, err
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return rows.CommandTag().RowsAffected(), nil
}

// mapPartError turns driver errors into repository errors. A nil err stays nil.
func mapPartError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "parts_user_part_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "parts_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		case pgError.Code == "23503" && pgError.ConstraintName == "parts_clientele_id_fkey":
			return ErrInvalidClientele
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
