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

// UserBranchRepository runs the SQL queries for the user_branch table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserBranchRepository struct {
	database db.Querier
}

// NewUserBranchRepository creates a UserBranchRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserBranchRepository(database db.Querier) *UserBranchRepository {
	return &UserBranchRepository{database: database}
}

// FindAll returns the user's branches ordered by name. A non-nil userCityID keeps
// only the branches of that city.
func (userBranchRepository *UserBranchRepository) FindAll(ctx context.Context, userID int64, userCityID *int64) ([]model.UserBranch, error) {
	query := `SELECT u.id, u.user_id, u.user_city_id, p.name AS city_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_branch u
		JOIN user_city p ON p.id = u.user_city_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_city_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, userID, userCityID)
	if err != nil {
		return nil, fmt.Errorf("find all branches: %w", err)
	}
	branches, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserBranch])
	if err != nil {
		return nil, fmt.Errorf("find all branches: %w", err)
	}
	return branches, nil
}

func (userBranchRepository *UserBranchRepository) FindByID(ctx context.Context, userID, id int64) (model.UserBranch, error) {
	query := `SELECT u.id, u.user_id, u.user_city_id, p.name AS city_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_branch u
		JOIN user_city p ON p.id = u.user_city_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, id, userID)
	if err != nil {
		return model.UserBranch{}, fmt.Errorf("find branch: %w", err)
	}
	branch, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserBranch])
	return branch, mapUserBranchError("find branch", err)
}

func (userBranchRepository *UserBranchRepository) Create(ctx context.Context, branch model.UserBranch) (model.UserBranch, error) {
	query := `WITH saved AS (
			INSERT INTO user_branch (user_id, user_city_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_city_id, p.name AS city_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_city p ON p.id = u.user_city_id`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, branch.UserID, branch.UserCityID, branch.Name)
	if err != nil {
		return model.UserBranch{}, fmt.Errorf("create branch: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserBranch])
	return created, mapUserBranchError("create branch", err)
}

// Update replaces the name and city of branch.ID owned by branch.UserID.
func (userBranchRepository *UserBranchRepository) Update(ctx context.Context, branch model.UserBranch) (model.UserBranch, error) {
	query := `WITH saved AS (
			UPDATE user_branch
			SET user_city_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_city_id, p.name AS city_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_city p ON p.id = u.user_city_id`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, branch.ID, branch.UserID, branch.UserCityID, branch.Name)
	if err != nil {
		return model.UserBranch{}, fmt.Errorf("update branch: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserBranch])
	return updated, mapUserBranchError("update branch", err)
}

// Delete soft-deletes the branch: the row is kept with is_deleted = true.
func (userBranchRepository *UserBranchRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_branch
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, id, userID)
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

// HasLocations reports whether the branch id owned by userID has active locations.
func (userBranchRepository *UserBranchRepository) HasLocations(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_location WHERE user_branch_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userBranchRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check locations: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check locations: %w", err)
	}
	return found, nil
}

// mapUserBranchError turns driver errors into repository errors. A nil err stays nil.
func mapUserBranchError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_branch_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
