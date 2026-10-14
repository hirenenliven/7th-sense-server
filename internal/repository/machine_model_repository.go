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

// MachineModelRepository runs the SQL queries for the models table (machine models of
// a controller). Soft-deleted models (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.MachineModel.
type MachineModelRepository struct {
	dbPool *pgxpool.Pool
}

// NewMachineModelRepository creates a MachineModelRepository that uses dbPool.
func NewMachineModelRepository(dbPool *pgxpool.Pool) *MachineModelRepository {
	return &MachineModelRepository{dbPool: dbPool}
}

// machineModelColumns reads a models row aliased mo plus its controller (c) and brand (b).
const machineModelColumns = `mo.id, mo.controller_id, mo.name,
		mo.is_deleted, mo.deleted_at, mo.created_at, mo.updated_at,
		COALESCE(c.name, '') AS controller_name,
		c.brand_id,
		COALESCE(b.name, '') AS brand_name`

// machineModelJoins joins the controller and brand of the models row aliased mo.
const machineModelJoins = `
		LEFT JOIN controllers c ON c.id = mo.controller_id
		LEFT JOIN brands b ON b.id = c.brand_id`

// FindAll returns the models by brand, controller and name. A non-nil controllerID
// keeps only that controller's models.
func (machineModelRepository *MachineModelRepository) FindAll(ctx context.Context, controllerID *int64) ([]model.MachineModel, error) {
	query := `SELECT ` + machineModelColumns + `
		FROM models mo` + machineModelJoins + `
		WHERE mo.is_deleted = false AND ($1::int IS NULL OR mo.controller_id = $1)
		ORDER BY b.name, c.name, mo.name`

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query, controllerID)
	if err != nil {
		return nil, fmt.Errorf("find all models: %w", err)
	}
	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MachineModel])
	if err != nil {
		return nil, fmt.Errorf("find all models: %w", err)
	}
	return models, nil
}

func (machineModelRepository *MachineModelRepository) FindByID(ctx context.Context, id int64) (model.MachineModel, error) {
	query := `SELECT ` + machineModelColumns + `
		FROM models mo` + machineModelJoins + `
		WHERE mo.id = $1 AND mo.is_deleted = false`

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query, id)
	if err != nil {
		return model.MachineModel{}, fmt.Errorf("find model: %w", err)
	}
	machineModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineModel])
	return machineModel, mapMachineModelError("find model", err)
}

func (machineModelRepository *MachineModelRepository) Create(ctx context.Context, machineModel model.MachineModel) (model.MachineModel, error) {
	query := `WITH mo AS (
			INSERT INTO models (controller_id, name)
			VALUES ($1, $2)
			RETURNING *
		)
		SELECT ` + machineModelColumns + `
		FROM mo` + machineModelJoins

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query, machineModel.ControllerID, machineModel.Name)
	if err != nil {
		return model.MachineModel{}, fmt.Errorf("create model: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineModel])
	return created, mapMachineModelError("create model", err)
}

// Update replaces the controller and name of machineModel.ID.
func (machineModelRepository *MachineModelRepository) Update(ctx context.Context, machineModel model.MachineModel) (model.MachineModel, error) {
	query := `WITH mo AS (
			UPDATE models
			SET controller_id = $2, name = $3
			WHERE id = $1 AND is_deleted = false
			RETURNING *
		)
		SELECT ` + machineModelColumns + `
		FROM mo` + machineModelJoins

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query,
		machineModel.ID, machineModel.ControllerID, machineModel.Name,
	)
	if err != nil {
		return model.MachineModel{}, fmt.Errorf("update model: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MachineModel])
	return updated, mapMachineModelError("update model", err)
}

// Delete soft-deletes the model: the row is kept with is_deleted = true.
func (machineModelRepository *MachineModelRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE models
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete model: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete model: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasMachines reports whether active machines use the model id.
func (machineModelRepository *MachineModelRepository) HasMachines(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM machines WHERE model_id = $1 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, machineModelRepository.dbPool, query, id)
	if err != nil {
		return false, fmt.Errorf("check model machines: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check model machines: %w", err)
	}
	return found, nil
}

// mapMachineModelError turns driver errors into repository errors. A nil err stays nil.
func mapMachineModelError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "models_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "models_controller_id_fkey": // foreign_key_violation
			return ErrInvalidController
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
