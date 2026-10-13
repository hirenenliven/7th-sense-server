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

// UserCityRepository runs the SQL queries for the user_city table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserCityRepository struct {
	database db.Querier
}

// NewUserCityRepository creates a UserCityRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserCityRepository(database db.Querier) *UserCityRepository {
	return &UserCityRepository{database: database}
}

// FindAll returns the user's cities ordered by name. A non-nil userStateID keeps
// only the cities of that state.
func (userCityRepository *UserCityRepository) FindAll(ctx context.Context, userID int64, userStateID *int64) ([]model.UserCity, error) {
	query := `SELECT u.id, u.user_id, u.user_state_id, p.name AS state_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_city u
		JOIN user_state p ON p.id = u.user_state_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_state_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userCityRepository.database, query, userID, userStateID)
	if err != nil {
		return nil, fmt.Errorf("find all cities: %w", err)
	}
	cities, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserCity])
	if err != nil {
		return nil, fmt.Errorf("find all cities: %w", err)
	}
	return cities, nil
}

func (userCityRepository *UserCityRepository) FindByID(ctx context.Context, userID, id int64) (model.UserCity, error) {
	query := `SELECT u.id, u.user_id, u.user_state_id, p.name AS state_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_city u
		JOIN user_state p ON p.id = u.user_state_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userCityRepository.database, query, id, userID)
	if err != nil {
		return model.UserCity{}, fmt.Errorf("find city: %w", err)
	}
	city, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCity])
	return city, mapUserCityError("find city", err)
}

func (userCityRepository *UserCityRepository) Create(ctx context.Context, city model.UserCity) (model.UserCity, error) {
	query := `WITH saved AS (
			INSERT INTO user_city (user_id, user_state_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_state_id, p.name AS state_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_state p ON p.id = u.user_state_id`

	rows, err := db.Execute(ctx, userCityRepository.database, query, city.UserID, city.UserStateID, city.Name)
	if err != nil {
		return model.UserCity{}, fmt.Errorf("create city: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCity])
	return created, mapUserCityError("create city", err)
}

// Update replaces the name and state of city.ID owned by city.UserID.
func (userCityRepository *UserCityRepository) Update(ctx context.Context, city model.UserCity) (model.UserCity, error) {
	query := `WITH saved AS (
			UPDATE user_city
			SET user_state_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_state_id, p.name AS state_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_state p ON p.id = u.user_state_id`

	rows, err := db.Execute(ctx, userCityRepository.database, query, city.ID, city.UserID, city.UserStateID, city.Name)
	if err != nil {
		return model.UserCity{}, fmt.Errorf("update city: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCity])
	return updated, mapUserCityError("update city", err)
}

// Delete soft-deletes the city: the row is kept with is_deleted = true.
func (userCityRepository *UserCityRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_city
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userCityRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete city: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete city: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasBranches reports whether the city id owned by userID has active branches.
func (userCityRepository *UserCityRepository) HasBranches(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_branch WHERE user_city_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userCityRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check branches: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check branches: %w", err)
	}
	return found, nil
}

// mapUserCityError turns driver errors into repository errors. A nil err stays nil.
func mapUserCityError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_city_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
