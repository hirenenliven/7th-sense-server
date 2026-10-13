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

// BrandRepository runs the SQL queries for the brands table. Soft-deleted
// brands (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Brand, so nullable
// columns are read with COALESCE and given their column name with AS.
type BrandRepository struct {
	dbPool *pgxpool.Pool
}

// NewBrandRepository creates a BrandRepository that uses dbPool.
func NewBrandRepository(dbPool *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{dbPool: dbPool}
}

func (brandRepository *BrandRepository) FindAll(ctx context.Context) ([]model.Brand, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM brands
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, brandRepository.dbPool, query)
	if err != nil {
		return nil, fmt.Errorf("find all brands: %w", err)
	}
	brands, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Brand])
	if err != nil {
		return nil, fmt.Errorf("find all brands: %w", err)
	}
	return brands, nil
}

func (brandRepository *BrandRepository) FindByID(ctx context.Context, id int64) (model.Brand, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM brands
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, brandRepository.dbPool, query, id)
	if err != nil {
		return model.Brand{}, fmt.Errorf("find brand: %w", err)
	}
	brand, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Brand])
	return brand, mapBrandError("find brand", err)
}

func (brandRepository *BrandRepository) Create(ctx context.Context, brand model.Brand) (model.Brand, error) {
	query := `INSERT INTO brands (name, code, description)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''))
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, brandRepository.dbPool, query,
		brand.Name, brand.Code, brand.Description,
	)
	if err != nil {
		return model.Brand{}, fmt.Errorf("create brand: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Brand])
	return created, mapBrandError("create brand", err)
}

// Update replaces the fields of brand.ID.
func (brandRepository *BrandRepository) Update(ctx context.Context, brand model.Brand) (model.Brand, error) {
	query := `UPDATE brands
		SET name = $2, code = NULLIF($3, ''), description = NULLIF($4, '')
		WHERE id = $1 AND is_deleted = false
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, brandRepository.dbPool, query,
		brand.ID, brand.Name, brand.Code, brand.Description,
	)
	if err != nil {
		return model.Brand{}, fmt.Errorf("update brand: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Brand])
	return updated, mapBrandError("update brand", err)
}

// Delete soft-deletes the brand: the row is kept with is_deleted = true.
func (brandRepository *BrandRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE brands
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, brandRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete brand: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete brand: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapBrandError turns driver errors into repository errors. A nil err stays nil.
func mapBrandError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "brands_name_key":
			return ErrDuplicateName
		case "brands_code_key":
			return ErrDuplicateCode
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
