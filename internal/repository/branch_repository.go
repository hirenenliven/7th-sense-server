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

func (branchRepository *BranchRepository) FindAll(ctx context.Context) ([]model.Branch, error) {
	query := `SELECT id, code, name, email, address, city, state_id, user_city_id, customer_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(color_code, '') AS color_code,
			COALESCE(status, '') AS status,
			is_deleted, deleted_at, created_at, updated_at
		FROM branches
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query)
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
	query := `SELECT id, code, name, email, address, city, state_id, user_city_id, customer_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(color_code, '') AS color_code,
			COALESCE(status, '') AS status,
			is_deleted, deleted_at, created_at, updated_at
		FROM branches
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, branchRepository.dbPool, query, id)
	if err != nil {
		return model.Branch{}, fmt.Errorf("find branch: %w", err)
	}
	branch, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Branch])
	return branch, mapBranchError("find branch", err)
}

func (branchRepository *BranchRepository) Create(ctx context.Context, branch model.Branch) (model.Branch, error) {
	query := `INSERT INTO branches (code, name, email, address, city, state_id, user_city_id, customer_id, gst_number, pan_number, color_code, status)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6, $7, $8, NULLIF($9, ''), NULLIF($10, ''), NULLIF($11, ''), NULLIF($12, ''))
		RETURNING id, code, name, email, address, city, state_id, user_city_id, customer_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(color_code, '') AS color_code,
			COALESCE(status, '') AS status,
			is_deleted, deleted_at, created_at, updated_at`

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
	query := `UPDATE branches
		SET code = $2, name = $3, email = NULLIF($4, ''), address = NULLIF($5, ''), city = NULLIF($6, ''),
			state_id = $7, user_city_id = $8, customer_id = $9, gst_number = NULLIF($10, ''), pan_number = NULLIF($11, ''),
			color_code = NULLIF($12, ''), status = COALESCE(NULLIF($13, ''), status)
		WHERE id = $1 AND is_deleted = false
		RETURNING id, code, name, email, address, city, state_id, user_city_id, customer_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(color_code, '') AS color_code,
			COALESCE(status, '') AS status,
			is_deleted, deleted_at, created_at, updated_at`

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
	return fmt.Errorf("%s: %w", operation, err)
}
