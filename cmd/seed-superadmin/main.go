// Command seed-superadmin adds the superadmin user (role 1). It reads the database
// settings from .env like the server does, and does nothing if an active user with
// the same username or email already exists, so it is safe to run more than once.
//
//	go run ./cmd/seed-superadmin
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
)

const (
	superadminUsername = "superadmin"
	superadminEmail    = "superadmin@example.com"
	superadminPassword = "password123"
	superadminRoleID   = 1
	superadminStatus   = "active"
)

func main() {
	if err := run(); err != nil {
		slog.Error("seed superadmin failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	appConfig := config.Load()

	dbPool, err := db.Connect(ctx, appConfig.DB)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superadminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	query := `INSERT INTO users (username, email, password, role_id, status)
		SELECT $1::text, $2::text, $3::text, $4::integer, $5::text
		WHERE NOT EXISTS (
			SELECT 1 FROM users
			WHERE (lower(username) = lower($1::text) OR lower(email) = lower($2::text)) AND is_deleted = false
		)
		RETURNING id`

	rows, err := db.Execute(ctx, dbPool, query,
		superadminUsername, superadminEmail, string(hashedPassword), superadminRoleID, superadminStatus,
	)
	if err != nil {
		return fmt.Errorf("insert superadmin: %w", err)
	}
	userID, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Info("superadmin already exists, nothing to do", "username", superadminUsername, "email", superadminEmail)
		return nil
	}
	if err != nil {
		return fmt.Errorf("insert superadmin: %w", err)
	}

	slog.Info("superadmin created", "id", userID, "username", superadminUsername, "email", superadminEmail)
	return nil
}
