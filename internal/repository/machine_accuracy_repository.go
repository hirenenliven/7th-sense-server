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

// machineAccuracyColumns are the columns of a model.MachineAccuracy read from machine_accuracies ma,
// with the number of machines and part processes (not deleted) that use it.
const machineAccuracyColumns = `ma.id, ma.name, ma.is_deleted, ma.deleted_at, ma.created_at, ma.updated_at,
	((SELECT COUNT(*) FROM machines m WHERE m.accuracy_id = ma.id AND m.is_deleted = false)
	 + (SELECT COUNT(*) FROM part_processes pp WHERE pp.accuracy_id = ma.id AND pp.is_deleted = false))::integer AS usage_count`

// MachineAccuracyRepository runs the SQL queries for the machine_accuracies table. Soft-deleted
// rows are invisible to every method.
type MachineAccuracyRepository struct {
	database db.Querier
}

// NewMachineAccuracyRepository creates a MachineAccuracyRepository that runs its queries on
// database (a connection pool or a transaction).
func NewMachineAccuracyRepository(database db.Querier) *MachineAccuracyRepository {
	return &MachineAccuracyRepository{database: database}
}

// FindAll returns the accuracies ordered by id (the order they were added).
func (accuracyRepository *MachineAccuracyRepository) FindAll(ctx context.Context) ([]model.MachineAccuracy, error) {
	query := `SELECT ` + machineAccuracyColumns + `
		FROM machine_accuracies ma
		WHERE ma.is_deleted = false
		ORDER BY ma.id`

	rows, err := db.Execute(ctx, accuracyRepository.database, query)
	if err != nil {
		return nil, fmt.Errorf("find all machine accuracies: %w", err)
	}
	accuracies, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineAccuracy])
	if err != nil {
		return nil, fmt.Errorf("find all machine accuracies: %w", err)
	}
	return accuracies, nil
}

// FindByID returns the accuracy, or ErrNotFound.
func (accuracyRepository *MachineAccuracyRepository) FindByID(ctx context.Context, id int64) (model.MachineAccuracy, error) {
	query := `SELECT ` + machineAccuracyColumns + `
		FROM machine_accuracies ma
		WHERE ma.id = $1 AND ma.is_deleted = false`

	rows, err := db.Execute(ctx, accuracyRepository.database, query, id)
	if err != nil {
		return model.MachineAccuracy{}, fmt.Errorf("find machine accuracy: %w", err)
	}
	accuracy, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineAccuracy])
	return accuracy, mapMachineAccuracyError("find machine accuracy", err)
}

// Create adds an accuracy and returns it; ErrDuplicateName when the name is taken
// (case-insensitive).
func (accuracyRepository *MachineAccuracyRepository) Create(ctx context.Context, name string) (model.MachineAccuracy, error) {
	query := `WITH saved AS (
			INSERT INTO machine_accuracies (name)
			VALUES ($1)
			RETURNING *
		)
		SELECT ` + machineAccuracyColumns + ` FROM saved ma`

	rows, err := db.Execute(ctx, accuracyRepository.database, query, name)
	if err != nil {
		return model.MachineAccuracy{}, mapMachineAccuracyError("create machine accuracy", err)
	}
	accuracy, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineAccuracy])
	return accuracy, mapMachineAccuracyError("create machine accuracy", err)
}

// Update renames accuracy id; ErrNotFound when there is no such accuracy, ErrDuplicateName when
// another accuracy has the name.
func (accuracyRepository *MachineAccuracyRepository) Update(ctx context.Context, id int64, name string) error {
	query := `UPDATE machine_accuracies
		SET name = $2, updated_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, accuracyRepository.database, query, id, name)
	if err != nil {
		return mapMachineAccuracyError("update machine accuracy", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes accuracy id; ErrNotFound when there is no such accuracy.
func (accuracyRepository *MachineAccuracyRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE machine_accuracies
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, accuracyRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete machine accuracy: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapMachineAccuracyError turns driver errors into repository errors. A nil err stays nil.
func mapMachineAccuracyError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" && pgError.ConstraintName == "machine_accuracies_name_key" { // unique_violation
		return ErrDuplicateName
	}
	return fmt.Errorf("%s: %w", operation, err)
}
