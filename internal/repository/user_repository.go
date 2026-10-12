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

// UserRepository runs the SQL queries for the users table. Soft-deleted
// users (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.User, so nullable
// columns are read with COALESCE and given their column name with AS.
type UserRepository struct {
	dbPool *pgxpool.Pool
}

// NewUserRepository creates a UserRepository that uses dbPool.
func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{dbPool: dbPool}
}

func (userRepository *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at
		FROM users
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, userRepository.dbPool, query)
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, fmt.Errorf("find all users: %w", err)
	}
	return users, nil
}

func (userRepository *UserRepository) FindByID(ctx context.Context, id int64) (model.User, error) {
	query := `SELECT id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userRepository.dbPool, query, id)
	if err != nil {
		return model.User{}, fmt.Errorf("find user: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return user, mapUserError("find user", err)
}

func (userRepository *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users (username, email, password, first_name, last_name, mobile, role_id, branch_id, company_id)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), $7, $8, $9)
		RETURNING id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at`

	rows, err := db.Execute(ctx, userRepository.dbPool, query,
		user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Mobile, user.RoleID, user.BranchID, user.CompanyID,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return created, mapUserError("create user", err)
}

// Update replaces the profile fields of user.ID. An empty Status keeps the
// current status. The password is not changed.
func (userRepository *UserRepository) Update(ctx context.Context, user model.User) (model.User, error) {
	query := `UPDATE users
		SET username = $2, email = $3,
			first_name = NULLIF($4, ''), last_name = NULLIF($5, ''), mobile = NULLIF($6, ''),
			role_id = $7, branch_id = $8, company_id = $9,
			status = COALESCE(NULLIF($10, ''), status)
		WHERE id = $1 AND is_deleted = false
		RETURNING id, username, email, password,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			COALESCE(mobile, '') AS mobile,
			role_id, branch_id, company_id,
			COALESCE(status, '') AS status,
			created_at, updated_at`

	rows, err := db.Execute(ctx, userRepository.dbPool, query,
		user.ID, user.Username, user.Email, user.FirstName, user.LastName, user.Mobile, user.RoleID, user.BranchID, user.CompanyID, user.Status,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	return updated, mapUserError("update user", err)
}

// Delete soft-deletes the user: the row is kept with is_deleted = true.
func (userRepository *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapUserError turns driver errors into repository errors. A nil err stays nil.
func mapUserError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // unique_violation
		switch pgError.ConstraintName {
		case "users_email_key":
			return ErrDuplicateEmail
		case "users_username_key":
			return ErrDuplicateUsername
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
