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

// UserShopRepository runs the SQL queries for the user_shop table. Every method is
// scoped to one user, and soft-deleted rows (is_deleted = true) are invisible
// to every method.
type UserShopRepository struct {
	database db.Querier
}

// NewUserShopRepository creates a UserShopRepository that runs its queries on database
// (a connection pool or a transaction).
func NewUserShopRepository(database db.Querier) *UserShopRepository {
	return &UserShopRepository{database: database}
}

// FindAll returns the user's shops ordered by name. A non-nil userLocationID keeps
// only the shops of that location.
func (userShopRepository *UserShopRepository) FindAll(ctx context.Context, userID *int64, userLocationID *int64) ([]model.UserShop, error) {
	query := `SELECT u.id, u.user_id, u.user_location_id, p.name AS location_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_shop u
		JOIN user_location p ON p.id = u.user_location_id
		WHERE ($1::bigint IS NULL OR u.user_id = $1) AND u.is_deleted = false
			AND ($2::bigint IS NULL OR u.user_location_id = $2)
		ORDER BY lower(u.name), u.id`

	rows, err := db.Execute(ctx, userShopRepository.database, query, userID, userLocationID)
	if err != nil {
		return nil, fmt.Errorf("find all shops: %w", err)
	}
	shops, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserShop])
	if err != nil {
		return nil, fmt.Errorf("find all shops: %w", err)
	}
	return shops, nil
}

func (userShopRepository *UserShopRepository) FindByID(ctx context.Context, userID, id int64) (model.UserShop, error) {
	query := `SELECT u.id, u.user_id, u.user_location_id, p.name AS location_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM user_shop u
		JOIN user_location p ON p.id = u.user_location_id
		WHERE u.id = $1 AND u.user_id = $2 AND u.is_deleted = false`

	rows, err := db.Execute(ctx, userShopRepository.database, query, id, userID)
	if err != nil {
		return model.UserShop{}, fmt.Errorf("find shop: %w", err)
	}
	shop, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserShop])
	return shop, mapUserShopError("find shop", err)
}

func (userShopRepository *UserShopRepository) Create(ctx context.Context, shop model.UserShop) (model.UserShop, error) {
	query := `WITH saved AS (
			INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_location_id, p.name AS location_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_location p ON p.id = u.user_location_id`

	rows, err := db.Execute(ctx, userShopRepository.database, query, shop.UserID, shop.UserLocationID, shop.Name)
	if err != nil {
		return model.UserShop{}, fmt.Errorf("create shop: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserShop])
	return created, mapUserShopError("create shop", err)
}

// Update replaces the name and location of shop.ID owned by shop.UserID.
func (userShopRepository *UserShopRepository) Update(ctx context.Context, shop model.UserShop) (model.UserShop, error) {
	query := `WITH saved AS (
			UPDATE user_shop
			SET user_location_id = $3, name = $4
			WHERE id = $1 AND user_id = $2 AND is_deleted = false
			RETURNING *
		)
		SELECT u.id, u.user_id, u.user_location_id, p.name AS location_name, u.name,
			u.is_deleted, u.deleted_at, u.created_at, u.updated_at
		FROM saved u
		JOIN user_location p ON p.id = u.user_location_id`

	rows, err := db.Execute(ctx, userShopRepository.database, query, shop.ID, shop.UserID, shop.UserLocationID, shop.Name)
	if err != nil {
		return model.UserShop{}, fmt.Errorf("update shop: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserShop])
	return updated, mapUserShopError("update shop", err)
}

// Delete soft-deletes the shop: the row is kept with is_deleted = true.
func (userShopRepository *UserShopRepository) Delete(ctx context.Context, userID, id int64) error {
	query := `UPDATE user_shop
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false`

	rows, err := db.Execute(ctx, userShopRepository.database, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete shop: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete shop: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// HasLines reports whether the shop id owned by userID has active lines.
func (userShopRepository *UserShopRepository) HasLines(ctx context.Context, userID, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_line WHERE user_shop_id = $1 AND user_id = $2 AND is_deleted = false
		)`

	rows, err := db.Execute(ctx, userShopRepository.database, query, id, userID)
	if err != nil {
		return false, fmt.Errorf("check lines: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check lines: %w", err)
	}
	return found, nil
}

// mapUserShopError turns driver errors into repository errors. A nil err stays nil.
func mapUserShopError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "user_shop_name_key":
			return ErrDuplicateName
		}
	}
	if errors.As(err, &pgError) && pgError.Code == "23503" { // foreign_key_violation
		return ErrInvalidParent
	}
	return fmt.Errorf("%s: %w", operation, err)
}
