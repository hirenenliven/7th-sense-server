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

func (controllerRepository *ControllerRepository) FindAll(ctx context.Context) ([]model.Controller, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(variant, '') AS variant,
			brand_id,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM controllers
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query)
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
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(variant, '') AS variant,
			brand_id,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM controllers
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query, id)
	if err != nil {
		return model.Controller{}, fmt.Errorf("find controller: %w", err)
	}
	controller, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Controller])
	return controller, mapControllerError("find controller", err)
}

func (controllerRepository *ControllerRepository) Create(ctx context.Context, controller model.Controller) (model.Controller, error) {
	query := `INSERT INTO controllers (name, code, variant, brand_id, description)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), $4, NULLIF($5, ''))
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(variant, '') AS variant,
			brand_id,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query,
		controller.Name, controller.Code, controller.Variant, controller.BrandID, controller.Description,
	)
	if err != nil {
		return model.Controller{}, fmt.Errorf("create controller: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Controller])
	return created, mapControllerError("create controller", err)
}

// Update replaces the fields of controller.ID.
func (controllerRepository *ControllerRepository) Update(ctx context.Context, controller model.Controller) (model.Controller, error) {
	query := `UPDATE controllers
		SET name = $2, code = NULLIF($3, ''), variant = NULLIF($4, ''),
			brand_id = $5, description = NULLIF($6, '')
		WHERE id = $1 AND is_deleted = false
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(variant, '') AS variant,
			brand_id,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, controllerRepository.dbPool, query,
		controller.ID, controller.Name, controller.Code, controller.Variant, controller.BrandID, controller.Description,
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
		case pgError.Code == "23505" && pgError.ConstraintName == "controllers_code_key":
			return ErrDuplicateCode
		case pgError.Code == "23503" && pgError.ConstraintName == "controllers_brand_id_fkey": // foreign_key_violation
			return ErrInvalidBrand
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
