package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// UserMachineFilterRepository runs the SQL queries for the user_machine_filters table
// (one row per user). The table has no soft delete: saving replaces the row.
type UserMachineFilterRepository struct {
	database db.Querier
}

// NewUserMachineFilterRepository creates a UserMachineFilterRepository that runs its
// queries on database (a connection pool or a transaction).
func NewUserMachineFilterRepository(database db.Querier) *UserMachineFilterRepository {
	return &UserMachineFilterRepository{database: database}
}

// FindByUserID returns the machine filter of userID, or ErrNotFound when it has none.
func (filterRepository *UserMachineFilterRepository) FindByUserID(ctx context.Context, userID int64) (model.UserMachineFilter, error) {
	query := `SELECT id, user_id, user_country_id, user_state_id, user_city_id, branch_id,
			user_location_id, user_shop_id, user_line_id, all_machines, created_at, updated_at
		FROM user_machine_filters
		WHERE user_id = $1`

	rows, err := db.Execute(ctx, filterRepository.database, query, userID)
	if err != nil {
		return model.UserMachineFilter{}, fmt.Errorf("find machine filter: %w", err)
	}
	filter, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserMachineFilter])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.UserMachineFilter{}, ErrNotFound
	}
	if err != nil {
		return model.UserMachineFilter{}, fmt.Errorf("find machine filter: %w", err)
	}
	return filter, nil
}

// Replace makes filter the machine filter of userID; nil removes it.
func (filterRepository *UserMachineFilterRepository) Replace(ctx context.Context, userID int64, filter *model.UserMachineFilter) error {
	if filter == nil {
		if err := execute(ctx, filterRepository.database, `DELETE FROM user_machine_filters WHERE user_id = $1`, userID); err != nil {
			return fmt.Errorf("remove machine filter: %w", err)
		}
		return nil
	}

	query := `INSERT INTO user_machine_filters (user_id, user_country_id, user_state_id, user_city_id, branch_id,
			user_location_id, user_shop_id, user_line_id, all_machines)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			user_country_id = EXCLUDED.user_country_id, user_state_id = EXCLUDED.user_state_id,
			user_city_id = EXCLUDED.user_city_id, branch_id = EXCLUDED.branch_id,
			user_location_id = EXCLUDED.user_location_id, user_shop_id = EXCLUDED.user_shop_id,
			user_line_id = EXCLUDED.user_line_id, all_machines = EXCLUDED.all_machines,
			updated_at = CURRENT_TIMESTAMP`
	if err := execute(ctx, filterRepository.database, query, userID, filter.UserCountryID, filter.UserStateID,
		filter.UserCityID, filter.BranchID, filter.UserLocationID, filter.UserShopID, filter.UserLineID, filter.AllMachines); err != nil {
		return fmt.Errorf("save machine filter: %w", err)
	}
	return nil
}

// levelInMachineFilter reports whether a user who is not deleted has a machine filter that names
// the organization level id in column (a column of user_machine_filters).
func levelInMachineFilter(ctx context.Context, database db.Querier, column string, id int64) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM user_machine_filters f
			JOIN users u ON u.id = f.user_id AND u.is_deleted = false
			WHERE f.` + column + ` = $1)`

	rows, err := db.Execute(ctx, database, query, id)
	if err != nil {
		return false, fmt.Errorf("check machine filter level: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check machine filter level: %w", err)
	}
	return found, nil
}

// IsInOrganization reports whether the levels of filter are one chain of companyID's
// organization (organization_paths), e.g. the state is under the country and the line
// under the shop.
func (filterRepository *UserMachineFilterRepository) IsInOrganization(ctx context.Context, companyID int64, filter model.UserMachineFilter) (bool, error) {
	query := `SELECT EXISTS (
			SELECT 1 FROM organization_paths p
			WHERE p.user_id = $1
				AND ($2::integer IS NULL OR p.user_country_id = $2)
				AND ($3::integer IS NULL OR p.user_state_id = $3)
				AND ($4::integer IS NULL OR p.user_city_id = $4)
				AND ($5::integer IS NULL OR p.branch_id = $5)
				AND ($6::integer IS NULL OR p.user_location_id = $6)
				AND ($7::integer IS NULL OR p.user_shop_id = $7)
				AND ($8::integer IS NULL OR p.user_line_id = $8))`

	rows, err := db.Execute(ctx, filterRepository.database, query, companyID, filter.UserCountryID, filter.UserStateID,
		filter.UserCityID, filter.BranchID, filter.UserLocationID, filter.UserShopID, filter.UserLineID)
	if err != nil {
		return false, fmt.Errorf("check machine filter: %w", err)
	}
	found, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, fmt.Errorf("check machine filter: %w", err)
	}
	return found, nil
}
