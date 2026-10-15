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

// clienteleSelect reads model.Clientele: the clienteles row (c) plus its customer's display name.
const clienteleSelect = `SELECT c.id, c.user_id, c.client_code, c.name, COALESCE(c.email, '') AS email,
		c.contact_number, c.address, c.country_name, c.state_name, c.city, c.status,
		c.is_deleted, c.deleted_at, c.created_at, c.updated_at,
		COALESCE(NULLIF(cp.company_name, ''), NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name
	FROM clienteles c
	LEFT JOIN users u ON u.id = c.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = c.user_id AND cp.is_deleted = false`

// ClienteleRepository runs the SQL queries for the clienteles table. Soft-deleted
// clienteles (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the clienteles of that customer
// (clienteles.user_id) are found.
type ClienteleRepository struct {
	database db.Querier
}

// NewClienteleRepository creates a ClienteleRepository that runs its queries on database
// (a connection pool or a transaction).
func NewClienteleRepository(database db.Querier) *ClienteleRepository {
	return &ClienteleRepository{database: database}
}

// FindAll returns the clienteles of ownerID (every clientele when nil), by name.
func (clienteleRepository *ClienteleRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.Clientele, error) {
	query := clienteleSelect + `
		WHERE c.is_deleted = false AND ($1::integer IS NULL OR c.user_id = $1)
		ORDER BY LOWER(c.name), c.id`

	rows, err := db.Execute(ctx, clienteleRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all clienteles: %w", err)
	}
	clienteles, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Clientele])
	if err != nil {
		return nil, fmt.Errorf("find all clienteles: %w", err)
	}
	return clienteles, nil
}

// FindByID returns the clientele, or ErrNotFound when it is not a clientele of ownerID.
func (clienteleRepository *ClienteleRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.Clientele, error) {
	query := clienteleSelect + `
		WHERE c.id = $1 AND c.is_deleted = false AND ($2::integer IS NULL OR c.user_id = $2)`

	rows, err := db.Execute(ctx, clienteleRepository.database, query, id, ownerID)
	if err != nil {
		return model.Clientele{}, fmt.Errorf("find clientele: %w", err)
	}
	clientele, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Clientele])
	return clientele, mapClienteleError("find clientele", err)
}

// Create inserts clientele and returns the new id.
func (clienteleRepository *ClienteleRepository) Create(ctx context.Context, clientele model.Clientele) (int64, error) {
	query := `INSERT INTO clienteles (user_id, client_code, name, email, contact_number, address,
			country_name, state_name, city, status)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7, $8, $9, $10)
		RETURNING id`

	rows, err := db.Execute(ctx, clienteleRepository.database, query,
		clientele.UserID, clientele.ClientCode, clientele.Name, clientele.Email, clientele.ContactNumber,
		clientele.Address, clientele.CountryName, clientele.StateName, clientele.City, clientele.Status)
	if err != nil {
		return 0, fmt.Errorf("create clientele: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapClienteleError("create clientele", err)
}

// Update replaces the fields of clientele.ID. A blank Status keeps the saved status.
func (clienteleRepository *ClienteleRepository) Update(ctx context.Context, clientele model.Clientele) error {
	query := `UPDATE clienteles
		SET user_id = $2, client_code = $3, name = $4, email = NULLIF($5, ''), contact_number = $6,
			address = $7, country_name = $8, state_name = $9, city = $10,
			status = COALESCE(NULLIF($11, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, clienteleRepository.database, query,
		clientele.ID, clientele.UserID, clientele.ClientCode, clientele.Name, clientele.Email, clientele.ContactNumber,
		clientele.Address, clientele.CountryName, clientele.StateName, clientele.City, clientele.Status)
	if err != nil {
		return mapClienteleError("update clientele", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// IsUsedByParts reports whether any part refers to clientele id.
func (clienteleRepository *ClienteleRepository) IsUsedByParts(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM parts WHERE clientele_id = $1 AND is_deleted = false)`

	rows, err := db.Execute(ctx, clienteleRepository.database, query, id)
	if err != nil {
		return false, fmt.Errorf("check clientele parts: %w", err)
	}
	used, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check clientele parts: %w", err)
	}
	return used, nil
}

// IsUsedByPartsOfAnotherCustomer reports whether a part of a customer other than
// ownerID refers to clientele id.
func (clienteleRepository *ClienteleRepository) IsUsedByPartsOfAnotherCustomer(ctx context.Context, id, ownerID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM parts WHERE clientele_id = $1 AND user_id <> $2 AND is_deleted = false)`

	rows, err := db.Execute(ctx, clienteleRepository.database, query, id, ownerID)
	if err != nil {
		return false, fmt.Errorf("check clientele parts: %w", err)
	}
	used, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check clientele parts: %w", err)
	}
	return used, nil
}

// Delete soft-deletes the clientele of ownerID (any clientele when nil).
func (clienteleRepository *ClienteleRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE clienteles
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, clienteleRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete clientele: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapClienteleError turns driver errors into repository errors. A nil err stays nil.
func mapClienteleError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "clienteles_user_client_code_key": // unique_violation
			return ErrDuplicateCode
		case pgError.Code == "23503" && pgError.ConstraintName == "clienteles_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
