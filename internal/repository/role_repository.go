package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// RoleRepository runs the SQL queries for the roles table. Roles are seed data
// and are never deleted, so the table has no is_deleted.
type RoleRepository struct {
	database db.Querier
}

// NewRoleRepository creates a RoleRepository that runs its queries on database
// (a connection pool or a transaction).
func NewRoleRepository(database db.Querier) *RoleRepository {
	return &RoleRepository{database: database}
}

// FindAll returns every role, by id.
func (roleRepository *RoleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	query := `SELECT id, name, COALESCE(description, '') AS description
		FROM roles
		ORDER BY id`

	rows, err := db.Execute(ctx, roleRepository.database, query)
	if err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
	}
	roles, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
	}
	return roles, nil
}

// FindByID returns the role with id, or ErrNotFound.
func (roleRepository *RoleRepository) FindByID(ctx context.Context, id int64) (model.Role, error) {
	query := `SELECT id, name, COALESCE(description, '') AS description
		FROM roles
		WHERE id = $1`

	rows, err := db.Execute(ctx, roleRepository.database, query, id)
	if err != nil {
		return model.Role{}, fmt.Errorf("find role: %w", err)
	}
	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Role])
	return role, mapRoleError("find role", err)
}

// mapRoleError turns driver errors into repository errors. A nil err stays nil.
func mapRoleError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
