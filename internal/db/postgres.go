package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/config"
)

// Connect opens a PostgreSQL connection pool and verifies it with a ping.
// The caller must Close the pool on shutdown.
func Connect(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     cfg.Name,
		RawQuery: "sslmode=" + url.QueryEscape(cfg.SSLMode),
	}

	pool, err := pgxpool.New(ctx, dsn.String())
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping %s@%s/%s: %w", cfg.User, dsn.Host, cfg.Name, err)
	}

	return pool, nil
}

// Execute executes any type of SQL query with parameters.
// For SELECT queries, use QueryRow for single row or Query for multiple rows.
// For INSERT/UPDATE/DELETE, use Exec to get rows affected.
func Execute(ctx context.Context, pool *pgxpool.Pool, query string, args ...any) (pgx.Rows, error) {
	return pool.Query(ctx, query, args...)
}
