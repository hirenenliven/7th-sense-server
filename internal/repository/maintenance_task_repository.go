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

// maintenanceTaskColumns are the columns of a model.MaintenanceTask read from maintenance_tasks mt,
// with the number of operator tasks (not deleted) that use it.
const maintenanceTaskColumns = `mt.id, mt.description, mt.status, mt.is_deleted, mt.deleted_at, mt.created_at, mt.updated_at,
	(SELECT COUNT(*) FROM operator_tasks ot WHERE ot.maintenance_task_id = mt.id AND ot.is_deleted = false)::integer AS usage_count`

// MaintenanceTaskRepository runs the SQL queries for the maintenance_tasks table (the master list
// of operator tasks). Soft-deleted rows are invisible to every method.
type MaintenanceTaskRepository struct {
	database db.Querier
}

// NewMaintenanceTaskRepository creates a MaintenanceTaskRepository that runs its queries on
// database (a connection pool or a transaction).
func NewMaintenanceTaskRepository(database db.Querier) *MaintenanceTaskRepository {
	return &MaintenanceTaskRepository{database: database}
}

// FindAll returns the tasks ordered by description.
func (taskRepository *MaintenanceTaskRepository) FindAll(ctx context.Context) ([]model.MaintenanceTask, error) {
	query := `SELECT ` + maintenanceTaskColumns + `
		FROM maintenance_tasks mt
		WHERE mt.is_deleted = false
		ORDER BY LOWER(mt.description), mt.id`

	rows, err := db.Execute(ctx, taskRepository.database, query)
	if err != nil {
		return nil, fmt.Errorf("find all maintenance tasks: %w", err)
	}
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MaintenanceTask])
	if err != nil {
		return nil, fmt.Errorf("find all maintenance tasks: %w", err)
	}
	return tasks, nil
}

// FindByID returns the task, or ErrNotFound.
func (taskRepository *MaintenanceTaskRepository) FindByID(ctx context.Context, id int64) (model.MaintenanceTask, error) {
	query := `SELECT ` + maintenanceTaskColumns + `
		FROM maintenance_tasks mt
		WHERE mt.id = $1 AND mt.is_deleted = false`

	rows, err := db.Execute(ctx, taskRepository.database, query, id)
	if err != nil {
		return model.MaintenanceTask{}, fmt.Errorf("find maintenance task: %w", err)
	}
	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MaintenanceTask])
	return task, mapMaintenanceTaskError("find maintenance task", err)
}

// Create adds an active task and returns it; ErrDuplicateName when the description is taken
// (case-insensitive).
func (taskRepository *MaintenanceTaskRepository) Create(ctx context.Context, description, status string) (model.MaintenanceTask, error) {
	query := `WITH saved AS (
			INSERT INTO maintenance_tasks (description, status)
			VALUES ($1, $2)
			RETURNING *
		)
		SELECT ` + maintenanceTaskColumns + ` FROM saved mt`

	rows, err := db.Execute(ctx, taskRepository.database, query, description, status)
	if err != nil {
		return model.MaintenanceTask{}, mapMaintenanceTaskError("create maintenance task", err)
	}
	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MaintenanceTask])
	return task, mapMaintenanceTaskError("create maintenance task", err)
}

// Update changes the description and, unless status is blank, the status of task id; ErrNotFound
// when there is no such task, ErrDuplicateName when another task has the description.
func (taskRepository *MaintenanceTaskRepository) Update(ctx context.Context, id int64, description, status string) error {
	query := `UPDATE maintenance_tasks
		SET description = $2, status = COALESCE(NULLIF($3, ''), status), updated_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, taskRepository.database, query, id, description, status)
	if err != nil {
		return mapMaintenanceTaskError("update maintenance task", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes task id; ErrNotFound when there is no such task.
func (taskRepository *MaintenanceTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE maintenance_tasks
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, taskRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete maintenance task: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapMaintenanceTaskError turns driver errors into repository errors. A nil err stays nil.
func mapMaintenanceTaskError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" && pgError.ConstraintName == "maintenance_tasks_description_key" { // unique_violation
		return ErrDuplicateName
	}
	return fmt.Errorf("%s: %w", operation, err)
}
