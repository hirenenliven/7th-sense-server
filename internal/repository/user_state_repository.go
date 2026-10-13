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

// UserStateRepository runs the SQL queries for the user_state table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserStateRepository struct {
	database db.Querier
}

// NewUserStateRepository creates a UserStateRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserStateRepository(database db.Querier) *UserStateRepository {
	return &UserStateRepository{database: database}
}

// FindAll returns the user's states ordered by name. A non-nil userCountryID keeps
// only the states of that country.
func (userStateRepository *UserStateRepository) FindAll(ctx context.Context, userID int64, userCountryID *int64) ([]model.UserState, error) {
	query := `SELECT u.id, u.user_id, u.user_country_id, p.name AS country_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_state u
		JOIN user_country p ON p.id = u.user_country_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_country_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userStateRepository.database, query, userID, userCountryID)
	if err != nil {
		return nil, fmt.Errorf("find all states: %w", err)
	}
	states, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserState])
	if err != nil {
		return nil, fmt.Errorf("find all states: %w", err)
	}
	return states, nil
}

func (userStateRepository *UserStateRepository) FindByID(ctx context.Context, userID, id int64) (model.UserState, error) {
	query := `SELECT u.id, u.user_id, u.user_country_id, p.name AS country_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_state u
		JOIN user_country p ON p.id = u.user_country_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userStateRepository.database, query, id, userID)
	if err != nil {
		return model.UserState{}, fmt.Errorf("find state: %w", err)
	}
	state, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserState])
	return state, mapUserStateError("find state", err)
}

func (userStateRepository *UserStateRepository) Create(ctx context.Context, state model.UserState) (model.UserState, error) {
	query := `WITH saved AS (
			INSERT INTO user_state (user_id, user_country_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_country_id, p.name AS country_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_country p ON p.id = u.user_country_id`

	rows, err := db.Execute(ctx, userStateRepository.database, query, state.UserID, state.UserCountryID, state.Name)
	if err != nil {
		return model.UserState{}, fmt.Errorf("create state: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserState])
	return created, mapUserStateError("create state", err)
}

// Delete soft-deletes the state: the row is kept with is_deleted = true.
func (userStateRepository *UserStateRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_state
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userStateRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasCities reports whether the state id owned by userID has active cities.
func (userStateRepository *UserStateRepository) HasCities(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_city WHERE user_state_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userStateRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check cities: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check cities: %w", err)
	}
	return found, nil
}

// mapUserStateError turns driver errors into repository errors. A nil err stays nil.
func mapUserStateError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_state_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
