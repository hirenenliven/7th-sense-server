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

// UserLineRepository runs the SQL queries for the user_line table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserLineRepository struct {
	database db.Querier
}

// NewUserLineRepository creates a UserLineRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserLineRepository(database db.Querier) *UserLineRepository {
	return &UserLineRepository{database: database}
}

// FindAll returns the user's lines ordered by name. A non-nil userShopID keeps
// only the lines of that shop.
func (userLineRepository *UserLineRepository) FindAll(ctx context.Context, userID int64, userShopID *int64) ([]model.UserLine, error) {
	query := `SELECT u.id, u.user_id, u.user_shop_id, p.name AS shop_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_line u
		JOIN user_shop p ON p.id = u.user_shop_id
		WHERE u.user_id = $1 AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_shop_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userLineRepository.database, query, userID, userShopID)
	if err != nil {
		return nil, fmt.Errorf("find all lines: %w", err)
	}
	lines, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserLine])
	if err != nil {
		return nil, fmt.Errorf("find all lines: %w", err)
	}
	return lines, nil
}

func (userLineRepository *UserLineRepository) FindByID(ctx context.Context, userID, id int64) (model.UserLine, error) {
	query := `SELECT u.id, u.user_id, u.user_shop_id, p.name AS shop_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_line u
		JOIN user_shop p ON p.id = u.user_shop_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userLineRepository.database, query, id, userID)
	if err != nil {
		return model.UserLine{}, fmt.Errorf("find line: %w", err)
	}
	line, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLine])
	return line, mapUserLineError("find line", err)
}

func (userLineRepository *UserLineRepository) Create(ctx context.Context, line model.UserLine) (model.UserLine, error) {
	query := `WITH saved AS (
			INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_shop_id, p.name AS shop_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_shop p ON p.id = u.user_shop_id`

	rows, err := db.Execute(ctx, userLineRepository.database, query, line.UserID, line.UserShopID, line.Name)
	if err != nil {
		return model.UserLine{}, fmt.Errorf("create line: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLine])
	return created, mapUserLineError("create line", err)
}

// Update replaces the name and shop of line.ID owned by line.UserID.
func (userLineRepository *UserLineRepository) Update(ctx context.Context, line model.UserLine) (model.UserLine, error) {
	query := `WITH saved AS (
			UPDATE user_line
			SET user_shop_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_shop_id, p.name AS shop_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_shop p ON p.id = u.user_shop_id`

	rows, err := db.Execute(ctx, userLineRepository.database, query, line.ID, line.UserID, line.UserShopID, line.Name)
	if err != nil {
		return model.UserLine{}, fmt.Errorf("update line: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserLine])
	return updated, mapUserLineError("update line", err)
}

// Delete soft-deletes the line: the row is kept with is_deleted = true.
func (userLineRepository *UserLineRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_line
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userLineRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete line: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete line: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasCells reports whether the line id owned by userID has active cells.
func (userLineRepository *UserLineRepository) HasCells(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_cell WHERE user_line_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userLineRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check cells: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check cells: %w", err)
	}
	return found, nil
}

// mapUserLineError turns driver errors into repository errors. A nil err stays nil.
func mapUserLineError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_line_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
