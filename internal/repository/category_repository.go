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

// CategoryRepository runs the SQL queries for the categories table. Soft-deleted
// categories (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Category, so nullable
// columns are read with COALESCE and given their column name with AS.
type CategoryRepository struct {
	dbPool *pgxpool.Pool
}

// NewCategoryRepository creates a CategoryRepository that uses dbPool.
func NewCategoryRepository(dbPool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{dbPool: dbPool}
}

func (categoryRepository *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM categories
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, categoryRepository.dbPool, query)
	if err != nil {
		return nil, fmt.Errorf("find all categories: %w", err)
	}
	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
	if err != nil {
		return nil, fmt.Errorf("find all categories: %w", err)
	}
	return categories, nil
}

func (categoryRepository *CategoryRepository) FindByID(ctx context.Context, id int64) (model.Category, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM categories
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, categoryRepository.dbPool, query, id)
	if err != nil {
		return model.Category{}, fmt.Errorf("find category: %w", err)
	}
	category, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Category])
	return category, mapCategoryError("find category", err)
}

func (categoryRepository *CategoryRepository) Create(ctx context.Context, category model.Category) (model.Category, error) {
	query := `INSERT INTO categories (name, code, description)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''))
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, categoryRepository.dbPool, query,
		category.Name, category.Code, category.Description,
	)
	if err != nil {
		return model.Category{}, fmt.Errorf("create category: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Category])
	return created, mapCategoryError("create category", err)
}

// Update replaces the fields of category.ID.
func (categoryRepository *CategoryRepository) Update(ctx context.Context, category model.Category) (model.Category, error) {
	query := `UPDATE categories
		SET name = $2, code = NULLIF($3, ''), description = NULLIF($4, '')
		WHERE id = $1 AND is_deleted = false
		RETURNING id, name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, categoryRepository.dbPool, query,
		category.ID, category.Name, category.Code, category.Description,
	)
	if err != nil {
		return model.Category{}, fmt.Errorf("update category: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Category])
	return updated, mapCategoryError("update category", err)
}

// Delete soft-deletes the category: the row is kept with is_deleted = true.
func (categoryRepository *CategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE categories
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, categoryRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapCategoryError turns driver errors into repository errors. A nil err stays nil.
func mapCategoryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "categories_name_key":
			return ErrDuplicateName
		case "categories_code_key":
			return ErrDuplicateCode
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
