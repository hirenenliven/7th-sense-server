// Package tests holds integration test suites that run against the PostgreSQL
// database configured in the project's .env file.
package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
)

// testUserID is the existing users.id the profile suites create profiles for.
const testUserID int64 = 6

var testDBPool *pgxpool.Pool

// TestMain connects to the database once for every test in the package.
// go test runs with the package directory as working directory, so the
// project's .env is loaded from the parent directory.
func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env")
	appConfig := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	pool, err := db.Connect(ctx, appConfig.DB)
	cancel()
	if err != nil {
		println("tests: cannot connect to database:", err.Error())
		os.Exit(1)
	}
	testDBPool = pool

	exitCode := m.Run()
	pool.Close()
	os.Exit(exitCode)
}

// prepareProfileTable makes table ready for a CRUD run on testUserID.
//
// Delete is a soft delete and user_id is UNIQUE, so a row left behind by a
// previous run would block Create; such soft-deleted rows are removed. An
// active profile is real data, so the test is skipped rather than touching it.
// The row the test creates is NOT deleted when the test finishes (for debugging).
func prepareProfileTable(t *testing.T, table string) {
	t.Helper()
	ctx := context.Background()

	var userExists bool
	if err := testDBPool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`, testUserID).Scan(&userExists); err != nil {
		t.Fatalf("check user %d: %v", testUserID, err)
	}
	if !userExists {
		t.Skipf("user %d does not exist in the database", testUserID)
	}

	var activeExists bool
	if err := testDBPool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM `+table+` WHERE user_id = $1 AND is_deleted = false)`, testUserID).Scan(&activeExists); err != nil {
		t.Fatalf("check active %s row: %v", table, err)
	}
	if activeExists {
		t.Skipf("user %d already has an active row in %s; not modifying real data", testUserID, table)
	}

	// Clean up soft-deleted rows only (not active rows)
	hardDelete := func() {
		if _, err := testDBPool.Exec(context.Background(), `DELETE FROM `+table+` WHERE user_id = $1 AND is_deleted = true`, testUserID); err != nil {
			t.Errorf("clean up soft-deleted %s: %v", table, err)
		}
	}
	hardDelete()
	// Cleanup removed to allow inspection of test data in database
	// t.Cleanup(hardDelete)
}
