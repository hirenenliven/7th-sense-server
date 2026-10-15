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

// UserCellRepository runs the SQL queries for the user_cell table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserCellRepository struct {
	database db.Querier
}

// NewUserCellRepository creates a UserCellRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserCellRepository(database db.Querier) *UserCellRepository {
	return &UserCellRepository{database: database}
}

// FindAll returns the user's cells ordered by name. A non-nil userLineID keeps
// only the cells of that line.
func (userCellRepository *UserCellRepository) FindAll(ctx context.Context, userID *int64, userLineID *int64) ([]model.UserCell, error) {
	query := `SELECT u.id, u.user_id, u.user_line_id, p.name AS line_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_cell u
		JOIN user_line p ON p.id = u.user_line_id
		WHERE ($1::bigint IS NULL OR u.user_id = $1) AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_line_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userCellRepository.database, query, userID, userLineID)
	if err != nil {
		return nil, fmt.Errorf("find all cells: %w", err)
	}
	cells, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserCell])
	if err != nil {
		return nil, fmt.Errorf("find all cells: %w", err)
	}
	return cells, nil
}

func (userCellRepository *UserCellRepository) FindByID(ctx context.Context, userID, id int64) (model.UserCell, error) {
	query := `SELECT u.id, u.user_id, u.user_line_id, p.name AS line_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_cell u
		JOIN user_line p ON p.id = u.user_line_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userCellRepository.database, query, id, userID)
	if err != nil {
		return model.UserCell{}, fmt.Errorf("find cell: %w", err)
	}
	cell, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCell])
	return cell, mapUserCellError("find cell", err)
}

func (userCellRepository *UserCellRepository) Create(ctx context.Context, cell model.UserCell) (model.UserCell, error) {
	query := `WITH saved AS (
			INSERT INTO user_cell (user_id, user_line_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_line_id, p.name AS line_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_line p ON p.id = u.user_line_id`

	rows, err := db.Execute(ctx, userCellRepository.database, query, cell.UserID, cell.UserLineID, cell.Name)
	if err != nil {
		return model.UserCell{}, fmt.Errorf("create cell: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCell])
	return created, mapUserCellError("create cell", err)
}

// Update replaces the name and line of cell.ID owned by cell.UserID.
func (userCellRepository *UserCellRepository) Update(ctx context.Context, cell model.UserCell) (model.UserCell, error) {
	query := `WITH saved AS (
			UPDATE user_cell
			SET user_line_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_line_id, p.name AS line_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_line p ON p.id = u.user_line_id`

	rows, err := db.Execute(ctx, userCellRepository.database, query, cell.ID, cell.UserID, cell.UserLineID, cell.Name)
	if err != nil {
		return model.UserCell{}, fmt.Errorf("update cell: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserCell])
	return updated, mapUserCellError("update cell", err)
}

// Delete soft-deletes the cell: the row is kept with is_deleted = true.
func (userCellRepository *UserCellRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_cell
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userCellRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete cell: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete cell: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasSupervisors reports whether the cell id owned by userID has active supervisors.
func (userCellRepository *UserCellRepository) HasSupervisors(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_supervisor WHERE user_cell_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userCellRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check supervisors: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check supervisors: %w", err)
	}
	return found, nil
}

// mapUserCellError turns driver errors into repository errors. A nil err stays nil.
func mapUserCellError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_cell_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
