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

// operatorTaskSelect reads model.OperatorTask: the operator_tasks row (ot) plus the display names
// of its customer, machine model (with controller and brand) and task.
const operatorTaskSelect = `SELECT ot.id, ot.user_id, ot.model_id, ot.maintenance_task_id,
		ot.is_deleted, ot.deleted_at, ot.created_at, ot.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name,
		COALESCE(mm.name, '') AS model_name,
		COALESCE(c.name, '') AS controller_name,
		COALESCE(b.name, '') AS brand_name,
		COALESCE(mt.description, '') AS task
	FROM operator_tasks ot
	LEFT JOIN users u ON u.id = ot.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = ot.user_id AND cp.is_deleted = false
	LEFT JOIN models mm ON mm.id = ot.model_id
	LEFT JOIN controllers c ON c.id = mm.controller_id
	LEFT JOIN brands b ON b.id = c.brand_id
	LEFT JOIN maintenance_tasks mt ON mt.id = ot.maintenance_task_id`

// OperatorTaskRepository runs the SQL queries for the operator_tasks table. Soft-deleted rows
// are invisible to every method. A nil ownerID does not narrow a query; otherwise only the
// tasks of that customer (operator_tasks.user_id) are found.
type OperatorTaskRepository struct {
	database db.Querier
}

// NewOperatorTaskRepository creates an OperatorTaskRepository that runs its queries on database
// (a connection pool or a transaction).
func NewOperatorTaskRepository(database db.Querier) *OperatorTaskRepository {
	return &OperatorTaskRepository{database: database}
}

// FindAll returns the operator tasks of ownerID (every customer's when nil), oldest first.
func (operatorTaskRepository *OperatorTaskRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.OperatorTask, error) {
	query := operatorTaskSelect + `
		WHERE ot.is_deleted = false AND ($1::integer IS NULL OR ot.user_id = $1)
		ORDER BY ot.id`

	rows, err := db.Execute(ctx, operatorTaskRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all operator tasks: %w", err)
	}
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.OperatorTask])
	if err != nil {
		return nil, fmt.Errorf("find all operator tasks: %w", err)
	}
	return tasks, nil
}

// FindByID returns the operator task, or ErrNotFound when it is not one of ownerID's.
func (operatorTaskRepository *OperatorTaskRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.OperatorTask, error) {
	query := operatorTaskSelect + `
		WHERE ot.id = $1 AND ot.is_deleted = false AND ($2::integer IS NULL OR ot.user_id = $2)`

	rows, err := db.Execute(ctx, operatorTaskRepository.database, query, id, ownerID)
	if err != nil {
		return model.OperatorTask{}, fmt.Errorf("find operator task: %w", err)
	}
	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.OperatorTask])
	return task, mapOperatorTaskError("find operator task", err)
}

// Create inserts the operator task (UserID, ModelID, MaintenanceTaskID) and returns its id.
func (operatorTaskRepository *OperatorTaskRepository) Create(ctx context.Context, task model.OperatorTask) (int64, error) {
	query := `INSERT INTO operator_tasks (user_id, model_id, maintenance_task_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	rows, err := db.Execute(ctx, operatorTaskRepository.database, query, task.UserID, task.ModelID, task.MaintenanceTaskID)
	if err != nil {
		return 0, mapOperatorTaskError("create operator task", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapOperatorTaskError("create operator task", err)
}

// Update replaces the customer, model and task of task.ID.
func (operatorTaskRepository *OperatorTaskRepository) Update(ctx context.Context, task model.OperatorTask) error {
	query := `UPDATE operator_tasks
		SET user_id = $2, model_id = $3, maintenance_task_id = $4, updated_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, operatorTaskRepository.database, query, task.ID, task.UserID, task.ModelID, task.MaintenanceTaskID)
	if err != nil {
		return mapOperatorTaskError("update operator task", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the operator task.
func (operatorTaskRepository *OperatorTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE operator_tasks
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, operatorTaskRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete operator task: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapOperatorTaskError turns driver errors into repository errors. A nil err stays nil.
func mapOperatorTaskError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "operator_tasks_customer_model_task_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "operator_tasks_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		case pgError.Code == "23503" && pgError.ConstraintName == "operator_tasks_model_id_fkey":
			return ErrInvalidModel
		case pgError.Code == "23503" && pgError.ConstraintName == "operator_tasks_maintenance_task_id_fkey":
			return ErrInvalidParent
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
