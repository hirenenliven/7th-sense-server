package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// AddressRepository runs the SQL queries for the addresses table. Soft-deleted
// addresses (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.Address, so nullable
// columns are read with COALESCE and given their column name with AS.
type AddressRepository struct {
	dbPool *pgxpool.Pool
}

// NewAddressRepository creates an AddressRepository that uses dbPool.
func NewAddressRepository(dbPool *pgxpool.Pool) *AddressRepository {
	return &AddressRepository{dbPool: dbPool}
}

func (addressRepository *AddressRepository) FindByUserID(ctx context.Context, userID int64) ([]model.Address, error) {
	query := `SELECT id, user_id,
			COALESCE(address_type, '') AS address_type,
			COALESCE(address_line1, '') AS address_line1,
			COALESCE(address_line2, '') AS address_line2,
			COALESCE(city, '') AS city,
			state_id, country_id,
			COALESCE(pin_code, '') AS pin_code,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			is_default,
			is_deleted, deleted_at, created_at, updated_at
		FROM addresses
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY id`

	rows, err := db.Execute(ctx, addressRepository.dbPool, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find addresses by user: %w", err)
	}
	addresses, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Address])
	if err != nil {
		return nil, fmt.Errorf("find addresses by user: %w", err)
	}
	return addresses, nil
}

func (addressRepository *AddressRepository) Create(ctx context.Context, address model.Address) (model.Address, error) {
	query := `INSERT INTO addresses (user_id, address_type, address_line1, address_line2, city, state_id, country_id, pin_code, contact_person, contact_phone, is_default)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6, $7, NULLIF($8, ''), NULLIF($9, ''), NULLIF($10, ''), $11)
		RETURNING id, user_id,
			COALESCE(address_type, '') AS address_type,
			COALESCE(address_line1, '') AS address_line1,
			COALESCE(address_line2, '') AS address_line2,
			COALESCE(city, '') AS city,
			state_id, country_id,
			COALESCE(pin_code, '') AS pin_code,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			is_default,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, addressRepository.dbPool, query,
		address.UserID, address.AddressType, address.AddressLine1, address.AddressLine2, address.City, address.StateID, address.CountryID, address.PinCode, address.ContactPerson, address.ContactPhone, address.IsDefault,
	)
	if err != nil {
		return model.Address{}, fmt.Errorf("create address: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Address])
	return created, mapAddressError("create address", err)
}

func (addressRepository *AddressRepository) Update(ctx context.Context, address model.Address) (model.Address, error) {
	query := `UPDATE addresses
		SET address_type = NULLIF($2, ''), address_line1 = NULLIF($3, ''), address_line2 = NULLIF($4, ''),
			city = NULLIF($5, ''), state_id = $6, country_id = $7, pin_code = NULLIF($8, ''),
			contact_person = NULLIF($9, ''), contact_phone = NULLIF($10, ''), is_default = $11
		WHERE id = $1 AND is_deleted = false
		RETURNING id, user_id,
			COALESCE(address_type, '') AS address_type,
			COALESCE(address_line1, '') AS address_line1,
			COALESCE(address_line2, '') AS address_line2,
			COALESCE(city, '') AS city,
			state_id, country_id,
			COALESCE(pin_code, '') AS pin_code,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			is_default,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, addressRepository.dbPool, query,
		address.ID, address.AddressType, address.AddressLine1, address.AddressLine2, address.City, address.StateID, address.CountryID, address.PinCode, address.ContactPerson, address.ContactPhone, address.IsDefault,
	)
	if err != nil {
		return model.Address{}, fmt.Errorf("update address: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Address])
	return updated, mapAddressError("update address", err)
}

// Delete soft-deletes the address: the row is kept with is_deleted = true.
func (addressRepository *AddressRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE addresses
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, addressRepository.dbPool, query, id)
	if err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapAddressError turns driver errors into repository errors. A nil err stays nil.
func mapAddressError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
