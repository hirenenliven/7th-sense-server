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

// BranchRepository runs the SQL queries for the branches table. Soft-deleted
// branches (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Branch, so nullable
// columns are read with COALESCE and given their column name with AS.
type BranchRepository struct {
	dbPool *pgxpool.Pool
}

// NewBranchRepository creates a BranchRepository that uses dbPool.
func NewBranchRepository(dbPool *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{dbPool: dbPool}
}

// branchColumns is the column list read into model.Branch from branches b,
// with the city name read from user_city c (LEFT JOIN: user_city_id is optional).
const branchColumns = `b.id, b.code, b.name,
			COALESCE(b.email, '') AS email,
			COALESCE(b.address, '') AS address,
			COALESCE(b.city, '') AS city,
			b.state_id, b.user_city_id,
			COALESCE(c.name, '') AS city_name,
			b.customer_id,
			COALESCE(b.gst_number, '') AS gst_number,
			COALESCE(b.pan_number, '') AS pan_number,
			COALESCE(b.color_code, '') AS color_code,
			COALESCE(b.status, '') AS status,
			b.is_deleted, b.deleted_at, b.created_at, b.updated_at`

// FindAll returns the branches ordered by id. A non-nil userCityID keeps only the
// branches of that city.
func (branchRepository *BranchRepository) FindAll(ctx context.Context, userCityID *int64) ([]model.Branch, error) {
	query := `SELECT ` + branchColumns + `
		FROM branches b
		LEFT JOIN user_city c ON c.id = b.user_city_id
		WHERE b.is_deleted = false
			AND ($1::bigint IS NULL OR b.user_city_id = $1)
		ORDER BY b.id`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query, userCityID)
	if err != nil {
		return nil, fmt.Errorf("find all branches: %w", err)
	}
	branches, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Branch])
	if err != nil {
		return nil, fmt.Errorf("find all branches: %w", err)
	}
	return branches, nil
}

func (branchRepository *BranchRepository) FindByID(ctx context.Context, id int64) (model.Branch, error) {
	query := `SELECT ` + branchColumns + `
		FROM branches b
		LEFT JOIN user_city c ON c.id = b.user_city_id
		WHERE b.id = $1 AND b.is_deleted = false`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query, id)
	if err != nil {
		return model.Branch{}, fmt.Errorf("find branch: %w", err)
	}
	branch, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Branch])
	return branch, mapBranchError("find branch", err)
}

func (branchRepository *BranchRepository) Create(ctx context.Context, branch model.Branch) (model.Branch, error) {
	query := `WITH saved AS (
			INSERT INTO branches (code, name, email, address, city, state_id, user_city_id, customer_id, gst_number, pan_number, color_code, status)
			VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6, $7, $8, NULLIF($9, ''), NULLIF($10, ''), NULLIF($11, ''), NULLIF($12, ''))
			RETURNING *
		)
		SELECT ` + branchColumns + `
		FROM saved b
		LEFT JOIN user_city c ON c.id = b.user_city_id`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query,
		branch.Code, branch.Name, branch.Email, branch.Address, branch.City, branch.StateID, branch.UserCityID, branch.CustomerID, branch.GSTNumber, branch.PANNumber, branch.ColorCode, branch.Status,
	)
	if err != nil {
		return model.Branch{}, fmt.Errorf("create branch: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Branch])
	return created, mapBranchError("create branch", err)
}

// Update replaces the profile fields of branch.ID. An empty Status keeps the
// current status.
func (branchRepository *BranchRepository) Update(ctx context.Context, branch model.Branch) (model.Branch, error) {
	query := `WITH saved AS (
			UPDATE branches
			SET code = $2, name = $3, email = NULLIF($4, ''), address = NULLIF($5, ''), city = NULLIF($6, ''),
				state_id = $7, user_city_id = $8, customer_id = $9, gst_number = NULLIF($10, ''), pan_number = NULLIF($11, ''),
				color_code = NULLIF($12, ''), status = COALESCE(NULLIF($13, ''), status)
			WHERE id = $1 AND is_deleted = false
			RETURNING *
		)
		SELECT ` + branchColumns + `
		FROM saved b
		LEFT JOIN user_city c ON c.id = b.user_city_id`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query,
		branch.ID, branch.Code, branch.Name, branch.Email, branch.Address, branch.City, branch.StateID, branch.UserCityID, branch.CustomerID, branch.GSTNumber, branch.PANNumber, branch.ColorCode, branch.Status,
	)
	if err != nil {
		return model.Branch{}, fmt.Errorf("update branch: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Branch])
	return updated, mapBranchError("update branch", err)
}

// Delete soft-deletes the branch: the row is kept with is_deleted = true.
func (branchRepository *BranchRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE branches
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete branch: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete branch: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasLocations reports whether the branch id has active locations (user_location.branch_id).
func (branchRepository *BranchRepository) HasLocations(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_location WHERE branch_id = $1 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query, id)
	if err != nil {
		return false, fmt.Errorf("check locations: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check locations: %w", err)
	}
	return found, nil
}

// mapBranchError turns driver errors into repository errors. A nil err stays nil.
func mapBranchError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "branches_code_key":
			return ErrDuplicateCode
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" && pgError.ConstraintName == "branches_user_city_id_fkey" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
