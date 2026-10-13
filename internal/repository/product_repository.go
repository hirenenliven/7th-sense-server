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

// ProductRepository runs the SQL queries for the products table. Soft-deleted
// products (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Product, so nullable
// columns are read with COALESCE and given their column name with AS.
type ProductRepository struct {
	dbPool *pgxpool.Pool
}

// NewProductRepository creates a ProductRepository that uses dbPool.
func NewProductRepository(dbPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{dbPool: dbPool}
}

func (productRepository *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			category_id,
			COALESCE(image, '') AS image,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM products
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, productRepository.dbPool, query)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Product])
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	return products, nil
}

func (productRepository *ProductRepository) FindByID(ctx context.Context, id int64) (model.Product, error) {
	query := `SELECT id, name,
			COALESCE(code, '') AS code,
			category_id,
			COALESCE(image, '') AS image,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, productRepository.dbPool, query, id)
	if err != nil {
		return model.Product{}, fmt.Errorf("find product: %w", err)
	}
	product, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Product])
	return product, mapProductError("find product", err)
}

func (productRepository *ProductRepository) Create(ctx context.Context, product model.Product) (model.Product, error) {
	query := `INSERT INTO products (name, code, category_id, image, description)
		VALUES ($1, NULLIF($2, ''), $3, NULLIF($4, ''), NULLIF($5, ''))
		RETURNING id, name,
			COALESCE(code, '') AS code,
			category_id,
			COALESCE(image, '') AS image,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, productRepository.dbPool, query,
		product.Name, product.Code, product.CategoryID, product.Image, product.Description,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("create product: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Product])
	return created, mapProductError("create product", err)
}

// Update replaces the fields of product.ID.
func (productRepository *ProductRepository) Update(ctx context.Context, product model.Product) (model.Product, error) {
	query := `UPDATE products
		SET name = $2, code = NULLIF($3, ''), category_id = $4,
			image = NULLIF($5, ''), description = NULLIF($6, '')
		WHERE id = $1 AND is_deleted = false
		RETURNING id, name,
			COALESCE(code, '') AS code,
			category_id,
			COALESCE(image, '') AS image,
			COALESCE(description, '') AS description,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, productRepository.dbPool, query,
		product.ID, product.Name, product.Code, product.CategoryID, product.Image, product.Description,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("update product: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Product])
	return updated, mapProductError("update product", err)
}

// Delete soft-deletes the product: the row is kept with is_deleted = true.
func (productRepository *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE products
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, productRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapProductError turns driver errors into repository errors. A nil err stays nil.
func mapProductError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "products_name_key": // unique_violation
			return ErrDuplicateName
		case pgError.Code == "23505" && pgError.ConstraintName == "products_code_key":
			return ErrDuplicateCode
		case pgError.Code == "23503" && pgError.ConstraintName == "products_category_id_fkey": // foreign_key_violation
			return ErrInvalidCategory
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
