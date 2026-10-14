package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
)

// UserMachineAccessRepository runs the SQL queries for the user_machine_access
// table. The table has no soft delete: removing an assignment deletes its row.
type UserMachineAccessRepository struct {
	database db.Querier
}

// NewUserMachineAccessRepository creates a UserMachineAccessRepository that runs
// its queries on database (a connection pool or a transaction).
func NewUserMachineAccessRepository(database db.Querier) *UserMachineAccessRepository {
	return &UserMachineAccessRepository{database: database}
}

// FindMachineIDs returns the ids of the machines assigned to userID, by id.
// Soft-deleted machines are left out.
func (accessRepository *UserMachineAccessRepository) FindMachineIDs(ctx context.Context, userID int64) ([]int64, error) {
	query := `SELECT a.machine_id
		FROM user_machine_access a
		JOIN machines m ON m.id = a.machine_id AND m.is_deleted = false
		WHERE a.user_id = $1
		ORDER BY a.machine_id`

	rows, err := db.Execute(ctx, accessRepository.database, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find machine access: %w", err)
	}
	machineIDs, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return nil, fmt.Errorf("find machine access: %w", err)
	}
	return machineIDs, nil
}

// Replace makes machineIDs the machines assigned to userID: assignments not in
// the list are deleted and missing ones are added. An empty list removes all.
func (accessRepository *UserMachineAccessRepository) Replace(ctx context.Context, userID int64, machineIDs []int64) error {
	if machineIDs == nil {
		machineIDs = []int64{}
	}
	deleteQuery := `DELETE FROM user_machine_access
		WHERE user_id = $1 AND machine_id <> ALL($2::bigint[])`
	if err := execute(ctx, accessRepository.database, deleteQuery, userID, machineIDs); err != nil {
		return fmt.Errorf("remove machine access: %w", err)
	}

	insertQuery := `INSERT INTO user_machine_access (user_id, machine_id)
		SELECT $1, machine_id FROM unnest($2::bigint[]) AS machine_id
		ON CONFLICT (user_id, machine_id) DO NOTHING`
	if err := execute(ctx, accessRepository.database, insertQuery, userID, machineIDs); err != nil {
		return fmt.Errorf("add machine access: %w", err)
	}
	return nil
}

// execute runs a query whose rows are not needed.
func execute(ctx context.Context, database db.Querier, query string, args ...any) error {
	rows, err := db.Execute(ctx, database, query, args...)
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}
