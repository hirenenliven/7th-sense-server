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

// ControllerRepository runs the SQL queries for the controllers table. Soft-deleted
// controllers (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Controller, so nullable
// columns are read with COALESCE and given their column name with AS.
type ControllerRepository struct {
	dbPool *pgxpool.Pool
}

// NewControllerRepository creates a ControllerRepository that uses dbPool.
func NewControllerRepository(dbPool *pgxpool.Pool) *ControllerRepository {
	return &ControllerRepository{dbPool: dbPool}
}

// controllerColumns reads a controllers row aliased c plus its brand name (brands b).
const controllerColumns = `c.id, c.name,
		c.brand_id,
		COALESCE(c.description, '') AS description,
		c.is_deleted, c.deleted_at, c.created_at, c.updated_at,
		COALESCE(b.name, '') AS brand_name`

// FindAll returns the controllers by brand and name. A non-nil brandID keeps only
// that brand's controllers.
func (controllerRepository *ControllerRepository) FindAll(ctx context.Context, brandID *int64) ([]model.Controller, error) {
	query := `SELECT ` + controllerColumns + `
		FROM controllers c
		LEFT JOIN brands b ON b.id = c.brand_id
		WHERE c.is_deleted = false AND ($1::int IS NULL OR c.brand_id = $1)
		ORDER BY b.name, c.name`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("find all controllers: %w", err)
	}
	controllers, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Controller])
	if err != nil {
		return nil, fmt.Errorf("find all controllers: %w", err)
	}
	return controllers, nil
}

func (controllerRepository *ControllerRepository) FindByID(ctx context.Context, id int64) (model.Controller, error) {
	query := `SELECT ` + controllerColumns + `
		FROM controllers c
		LEFT JOIN brands b ON b.id = c.brand_id
		WHERE c.id = $1 AND c.is_deleted = false`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query, id)
	if err != nil {
		return model.Controller{}, fmt.Errorf("find controller: %w", err)
	}
	controller, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Controller])
	return controller, mapControllerError("find controller", err)
}

func (controllerRepository *ControllerRepository) Create(ctx context.Context, controller model.Controller) (model.Controller, error) {
	query := `WITH c AS (
			INSERT INTO controllers (name, brand_id, description)
			VALUES ($1, $2, NULLIF($3, ''))
			RETURNING *
		)
		SELECT ` + controllerColumns + `
		FROM c
		LEFT JOIN brands b ON b.id = c.brand_id`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query,
		controller.Name, controller.BrandID, controller.Description,
	)
	if err != nil {
		return model.Controller{}, fmt.Errorf("create controller: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Controller])
	return created, mapControllerError("create controller", err)
}

// Update replaces the fields of controller.ID.
func (controllerRepository *ControllerRepository) Update(ctx context.Context, controller model.Controller) (model.Controller, error) {
	query := `WITH c AS (
			UPDATE controllers
			SET name = $2, brand_id = $3, description = NULLIF($4, '')
			WHERE id = $1 AND is_deleted = false
			RETURNING *
		)
		SELECT ` + controllerColumns + `
		FROM c
		LEFT JOIN brands b ON b.id = c.brand_id`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query,
		controller.ID, controller.Name, controller.BrandID, controller.Description,
	)
	if err != nil {
		return model.Controller{}, fmt.Errorf("update controller: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Controller])
	return updated, mapControllerError("update controller", err)
}

// Delete soft-deletes the controller: the row is kept with is_deleted = true.
func (controllerRepository *ControllerRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE controllers
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete controller: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete controller: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasModels reports whether the controller id has active models.
func (controllerRepository *ControllerRepository) HasModels(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM models WHERE controller_id = $1 AND is_deleted = false
		)`
	return controllerRepository.exists(ctx, "check controller models", query, id)
}

// HasMachines reports whether active machines use the controller id.
func (controllerRepository *ControllerRepository) HasMachines(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM machines WHERE controller_id = $1 AND is_deleted = false
		)`
	return controllerRepository.exists(ctx, "check controller machines", query, id)
}

// exists runs a SELECT EXISTS query with id.
func (controllerRepository *ControllerRepository) exists(ctx context.Context, operation, query string, id int64) (bool, error) {
	rows, err := db.Execute(ctx, controllerRepository.dbPool, query, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", operation, err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("%s: %w", operation, err)
	}
	return found, nil
}

// mapControllerError turns driver errors into repository errors. A nil err stays nil.
func mapControllerError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "controllers_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "controllers_brand_id_fkey": // foreign_key_violation
			return ErrInvalidBrand
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
