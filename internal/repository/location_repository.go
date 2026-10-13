package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
)

// LocationRepository runs the lookup queries for the countries and states tables.
// Names are matched case-insensitively.
type LocationRepository struct {
	database db.Querier
}

// NewLocationRepository creates a LocationRepository that runs its queries on database
// (a connection pool or a transaction).
func NewLocationRepository(database db.Querier) *LocationRepository {
	return &LocationRepository{database: database}
}

// FindCountryIDByName returns the id of the country called name.
func (locationRepository *LocationRepository) FindCountryIDByName(ctx context.Context, name string) (int64, error) {
	query := `SELECT id FROM countries WHERE lower(name) = lower($1) ORDER BY id LIMIT 1`

	rows, err := db.Execute(ctx, locationRepository.database, query, name)
	if err != nil {
		return 0, fmt.Errorf("find country: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapLocationError("find country", err)
}

// FindStateIDByName returns the id of the state called name in countryID.
func (locationRepository *LocationRepository) FindStateIDByName(ctx context.Context, countryID int64, name string) (int64, error) {
	query := `SELECT id FROM states WHERE country_id = $1 AND lower(name) = lower($2) ORDER BY id LIMIT 1`

	rows, err := db.Execute(ctx, locationRepository.database, query, countryID, name)
	if err != nil {
		return 0, fmt.Errorf("find state: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapLocationError("find state", err)
}

// mapLocationError turns driver errors into repository errors. A nil err stays nil.
func mapLocationError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
