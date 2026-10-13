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

// UserCountryRepository runs the SQL queries for the user_country table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserCountryRepository struct {
	database db.Querier
}

// NewUserCountryRepository creates a UserCountryRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserCountryRepository(database db.Querier) *UserCountryRepository {
	return &UserCountryRepository{database: database}
}

// FindAll returns the user's countries ordered by name.
func (userCountryRepository *UserCountryRepository) FindAll(ctx context.Context, userID int64) ([]model.UserCountry, error) {
	query := `SELECT u.id, u.user_id, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_country u
		WHERE u.user_id = $1 AND u.is_deleted = false
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userCountryRepository.database, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find all countries: %w", err)
	}
	countries, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserCountry])
	if err != nil {
		return nil, fmt.Errorf("find all countries: %w", err)
	}
	return countries, nil
}

func (userCountryRepository *UserCountryRepository) FindByID(ctx context.Context, userID, id int64) (model.UserCountry, error) {
	query := `SELECT u.id, u.user_id, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_country u
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userCountryRepository.database, query, id, userID)
	if err != nil {
		return model.UserCountry{}, fmt.Errorf("find country: %w", err)
	}
	country, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCountry])
	return country, mapUserCountryError("find country", err)
}

func (userCountryRepository *UserCountryRepository) Create(ctx context.Context, country model.UserCountry) (model.UserCountry, error) {
	query := `WITH saved AS (
			INSERT INTO user_country (user_id, name) VALUES ($1, $2)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u`

	rows, err := db.Execute(ctx, userCountryRepository.database, query, country.UserID, country.Name)
	if err != nil {
		return model.UserCountry{}, fmt.Errorf("create country: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCountry])
	return created, mapUserCountryError("create country", err)
}

// Delete soft-deletes the country: the row is kept with is_deleted = true.
func (userCountryRepository *UserCountryRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_country
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userCountryRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete country: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete country: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasStates reports whether the country id owned by userID has active states.
func (userCountryRepository *UserCountryRepository) HasStates(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_state WHERE user_country_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userCountryRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check states: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check states: %w", err)
	}
	return found, nil
}

// mapUserCountryError turns driver errors into repository errors. A nil err stays nil.
func mapUserCountryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_country_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation: user_id
		return ErrInvalidUser
	}
	return fmt.Errorf("%s: %w", operation, err)
}
