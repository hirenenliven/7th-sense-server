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

// UserSupervisorRepository runs the SQL queries for the user_supervisor table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserSupervisorRepository struct {
	database db.Querier
}

// NewUserSupervisorRepository creates a UserSupervisorRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserSupervisorRepository(database db.Querier) *UserSupervisorRepository {
	return &UserSupervisorRepository{database: database}
}

// FindAll returns the user's supervisors ordered by name. A non-nil userCellID keeps
// only the supervisors of that cell.
func (userSupervisorRepository *UserSupervisorRepository) FindAll(ctx context.Context, userID int64, userCellID *int64) ([]model.UserSupervisor, error) {
	query := `SELECT u.id, u.user_id, u.user_cell_id, p.name AS cell_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_supervisor u
		JOIN user_cell p ON p.id = u.user_cell_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_cell_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userSupervisorRepository.database, query, userID, userCellID)
	if err != nil {
		return nil, fmt.Errorf("find all supervisors: %w", err)
	}
	supervisors, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserSupervisor])
	if err != nil {
		return nil, fmt.Errorf("find all supervisors: %w", err)
	}
	return supervisors, nil
}

func (userSupervisorRepository *UserSupervisorRepository) FindByID(ctx context.Context, userID, id int64) (model.UserSupervisor, error) {
	query := `SELECT u.id, u.user_id, u.user_cell_id, p.name AS cell_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_supervisor u
		JOIN user_cell p ON p.id = u.user_cell_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userSupervisorRepository.database, query, id, userID)
	if err != nil {
		return model.UserSupervisor{}, fmt.Errorf("find supervisor: %w", err)
	}
	supervisor, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserSupervisor])
	return supervisor, mapUserSupervisorError("find supervisor", err)
}

func (userSupervisorRepository *UserSupervisorRepository) Create(ctx context.Context, supervisor model.UserSupervisor) (model.UserSupervisor, error) {
	query := `WITH saved AS (
			INSERT INTO user_supervisor (user_id, user_cell_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_cell_id, p.name AS cell_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_cell p ON p.id = u.user_cell_id`

	rows, err := db.Execute(ctx, userSupervisorRepository.database, query, supervisor.UserID, supervisor.UserCellID, supervisor.Name)
	if err != nil {
		return model.UserSupervisor{}, fmt.Errorf("create supervisor: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserSupervisor])
	return created, mapUserSupervisorError("create supervisor", err)
}

// Update replaces the name and cell of supervisor.ID owned by supervisor.UserID.
func (userSupervisorRepository *UserSupervisorRepository) Update(ctx context.Context, supervisor model.UserSupervisor) (model.UserSupervisor, error) {
	query := `WITH saved AS (
			UPDATE user_supervisor
			SET user_cell_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_cell_id, p.name AS cell_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_cell p ON p.id = u.user_cell_id`

	rows, err := db.Execute(ctx, userSupervisorRepository.database, query, supervisor.ID, supervisor.UserID, supervisor.UserCellID, supervisor.Name)
	if err != nil {
		return model.UserSupervisor{}, fmt.Errorf("update supervisor: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserSupervisor])
	return updated, mapUserSupervisorError("update supervisor", err)
}

// Delete soft-deletes the supervisor: the row is kept with is_deleted = true.
func (userSupervisorRepository *UserSupervisorRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_supervisor
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userSupervisorRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete supervisor: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete supervisor: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapUserSupervisorError turns driver errors into repository errors. A nil err stays nil.
func mapUserSupervisorError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_supervisor_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
