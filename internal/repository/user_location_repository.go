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

// UserLocationRepository runs the SQL queries for the user_location table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserLocationRepository struct {
	database db.Querier
}

// NewUserLocationRepository creates a UserLocationRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserLocationRepository(database db.Querier) *UserLocationRepository {
	return &UserLocationRepository{database: database}
}

// FindAll returns the user's locations ordered by name. A non-nil userBranchID keeps
// only the locations of that branch.
func (userLocationRepository *UserLocationRepository) FindAll(ctx context.Context, userID int64, userBranchID *int64) ([]model.UserLocation, error) {
	query := `SELECT u.id, u.user_id, u.user_branch_id, p.name AS branch_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_location u
		JOIN user_branch p ON p.id = u.user_branch_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_branch_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, userID, userBranchID)
	if err != nil {
		return nil, fmt.Errorf("find all locations: %w", err)
	}
	locations, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserLocation])
	if err != nil {
		return nil, fmt.Errorf("find all locations: %w", err)
	}
	return locations, nil
}

func (userLocationRepository *UserLocationRepository) FindByID(ctx context.Context, userID, id int64) (model.UserLocation, error) {
	query := `SELECT u.id, u.user_id, u.user_branch_id, p.name AS branch_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_location u
		JOIN user_branch p ON p.id = u.user_branch_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, id, userID)
	if err != nil {
		return model.UserLocation{}, fmt.Errorf("find location: %w", err)
	}
	location, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLocation])
	return location, mapUserLocationError("find location", err)
}

func (userLocationRepository *UserLocationRepository) Create(ctx context.Context, location model.UserLocation) (model.UserLocation, error) {
	query := `WITH saved AS (
			INSERT INTO user_location (user_id, user_branch_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_branch_id, p.name AS branch_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_branch p ON p.id = u.user_branch_id`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, location.UserID, location.UserBranchID, location.Name)
	if err != nil {
		return model.UserLocation{}, fmt.Errorf("create location: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLocation])
	return created, mapUserLocationError("create location", err)
}

// Update replaces the name and branch of location.ID owned by location.UserID.
func (userLocationRepository *UserLocationRepository) Update(ctx context.Context, location model.UserLocation) (model.UserLocation, error) {
	query := `WITH saved AS (
			UPDATE user_location
			SET user_branch_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_branch_id, p.name AS branch_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_branch p ON p.id = u.user_branch_id`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, location.ID, location.UserID, location.UserBranchID, location.Name)
	if err != nil {
		return model.UserLocation{}, fmt.Errorf("update location: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLocation])
	return updated, mapUserLocationError("update location", err)
}

// Delete soft-deletes the location: the row is kept with is_deleted = true.
func (userLocationRepository *UserLocationRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_location
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete location: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete location: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasShops reports whether the location id owned by userID has active shops.
func (userLocationRepository *UserLocationRepository) HasShops(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_shop WHERE user_location_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userLocationRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check shops: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check shops: %w", err)
	}
	return found, nil
}

// mapUserLocationError turns driver errors into repository errors. A nil err stays nil.
func mapUserLocationError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_location_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
