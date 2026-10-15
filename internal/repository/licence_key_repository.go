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

// licenceKeySelect reads model.LicenceKey: the licence_keys row (k) plus its customer's
// display name. expires_on is read as "YYYY-MM-DD"; is_expired once it is before today.
const licenceKeySelect = `SELECT k.id, k.user_id, k.licence_key, k.duration, k.expires_on::text AS expires_on, k.status,
		k.is_deleted, k.deleted_at, k.created_at, k.updated_at,
		COALESCE(k.expires_on < CURRENT_DATE, false) AS is_expired,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name
	FROM licence_keys k
	LEFT JOIN users u ON u.id = k.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = k.user_id AND cp.is_deleted = false`

// LicenceKeyRepository runs the SQL queries for the licence_keys table. Soft-deleted rows
// (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the keys of that customer
// (licence_keys.user_id) are found.
type LicenceKeyRepository struct {
	database db.Querier
}

// NewLicenceKeyRepository creates a LicenceKeyRepository that runs its queries on database
// (a connection pool or a transaction).
func NewLicenceKeyRepository(database db.Querier) *LicenceKeyRepository {
	return &LicenceKeyRepository{database: database}
}

// FindAll returns the licence keys of ownerID (every customer's when nil), oldest first.
func (licenceKeyRepository *LicenceKeyRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.LicenceKey, error) {
	query := licenceKeySelect + `
		WHERE k.is_deleted = false AND ($1::integer IS NULL OR k.user_id = $1)
		ORDER BY k.id`

	rows, err := db.Execute(ctx, licenceKeyRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all licence keys: %w", err)
	}
	licenceKeys, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.LicenceKey])
	if err != nil {
		return nil, fmt.Errorf("find all licence keys: %w", err)
	}
	return licenceKeys, nil
}

// FindByIDs returns the licence keys ids (whoever owns them), oldest first.
func (licenceKeyRepository *LicenceKeyRepository) FindByIDs(ctx context.Context, ids []int64) ([]model.LicenceKey, error) {
	query := licenceKeySelect + `
		WHERE k.id = ANY($1) AND k.is_deleted = false
		ORDER BY k.id`

	rows, err := db.Execute(ctx, licenceKeyRepository.database, query, ids)
	if err != nil {
		return nil, fmt.Errorf("find licence keys: %w", err)
	}
	licenceKeys, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.LicenceKey])
	if err != nil {
		return nil, fmt.Errorf("find licence keys: %w", err)
	}
	return licenceKeys, nil
}

// FindByID returns the licence key, or ErrNotFound when it is not one of ownerID's.
func (licenceKeyRepository *LicenceKeyRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.LicenceKey, error) {
	query := licenceKeySelect + `
		WHERE k.id = $1 AND k.is_deleted = false AND ($2::integer IS NULL OR k.user_id = $2)`

	rows, err := db.Execute(ctx, licenceKeyRepository.database, query, id, ownerID)
	if err != nil {
		return model.LicenceKey{}, fmt.Errorf("find licence key: %w", err)
	}
	licenceKey, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.LicenceKey])
	return licenceKey, mapLicenceKeyError("find licence key", err)
}

// CreateMany inserts licenceKeys (UserID, LicenceKey, Duration, ExpiresOn, Status of each)
// in one statement (all or none) and returns the new ids.
func (licenceKeyRepository *LicenceKeyRepository) CreateMany(ctx context.Context, licenceKeys []model.LicenceKey) ([]int64, error) {
	userIDs := make([]int64, len(licenceKeys))
	keys := make([]string, len(licenceKeys))
	durations := make([]string, len(licenceKeys))
	expiresOn := make([]*string, len(licenceKeys))
	statuses := make([]string, len(licenceKeys))
	for index, licenceKey := range licenceKeys {
		userIDs[index], keys[index], durations[index] = licenceKey.UserID, licenceKey.LicenceKey, licenceKey.Duration
		expiresOn[index], statuses[index] = licenceKey.ExpiresOn, licenceKey.Status
	}

	query := `INSERT INTO licence_keys (user_id, licence_key, duration, expires_on, status)
		SELECT user_id, licence_key, duration, expires_on::date, status
		FROM UNNEST($1::integer[], $2::text[], $3::text[], $4::text[], $5::text[])
			AS new_keys (user_id, licence_key, duration, expires_on, status)
		RETURNING id`

	rows, err := db.Execute(ctx, licenceKeyRepository.database, query, userIDs, keys, durations, expiresOn, statuses)
	if err != nil {
		return nil, mapLicenceKeyError("create licence keys", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	return ids, mapLicenceKeyError("create licence keys", err)
}

// Update changes the customer, duration, end date and status of licenceKey.ID. A blank
// Status keeps the saved status.
func (licenceKeyRepository *LicenceKeyRepository) Update(ctx context.Context, licenceKey model.LicenceKey) error {
	query := `UPDATE licence_keys
		SET user_id = $2, duration = $3, expires_on = $4::date, status = COALESCE(NULLIF($5, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, licenceKeyRepository.database, query,
		licenceKey.ID, licenceKey.UserID, licenceKey.Duration, licenceKey.ExpiresOn, licenceKey.Status)
	if err != nil {
		return mapLicenceKeyError("update licence key", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the licence key.
func (licenceKeyRepository *LicenceKeyRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE licence_keys
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, licenceKeyRepository.database, query, id)
	if err != nil {
		return fmt.Errorf("delete licence key: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapLicenceKeyError turns driver errors into repository errors. A nil err stays nil.
func mapLicenceKeyError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "licence_keys_licence_key_key": // unique_violation
			return ErrDuplicateLicenceKey
		case pgError.Code == "23503" && pgError.ConstraintName == "licence_keys_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
