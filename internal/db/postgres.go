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
func Connect(ctx context.Context, dbConfig config.DBConfig) (*pgxpool.Pool, error) {
	connectionURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbConfig.User, dbConfig.Password),
		Host:     net.JoinHostPort(dbConfig.Host, dbConfig.Port),
		Path:     dbConfig.Name,
		RawQuery: "sslmode=" + url.QueryEscape(dbConfig.SSLMode),
	}

	pool, err := pgxpool.New(ctx, connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping %s@%s/%s: %w", dbConfig.User, connectionURL.Host, dbConfig.Name, err)
	}

	return pool, nil
}

// Querier runs SQL queries. Both *pgxpool.Pool and pgx.Tx implement it, so a
// repository created with a transaction runs its queries inside that transaction.
type Querier interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

// Execute runs any SQL query with parameters and returns its rows.
//   - SELECT / ... RETURNING: read the rows with pgx.CollectRows or pgx.CollectOneRow.
//   - INSERT / UPDATE / DELETE without RETURNING: call rows.Close(), check rows.Err(),
//     then use rows.CommandTag().RowsAffected().
func Execute(ctx context.Context, querier Querier, query string, args ...any) (pgx.Rows, error) {
	return querier.Query(ctx, query, args...)
}

// WithTransaction runs fn inside a transaction. The transaction is committed
// when fn returns nil and rolled back when it returns an error.
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, pool, fn)
}

// Listen runs LISTEN on channels on one connection taken from pool and calls notify for
// every notification until ctx is cancelled or the connection fails (the error is returned;
// the caller reconnects). ready is called once the LISTENs are active.
func Listen(ctx context.Context, pool *pgxpool.Pool, channels []string, ready func(), notify func(channel, payload string)) error {
	pooled, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db: acquire listen connection: %w", err)
	}
	// Taken out of the pool for good: it holds LISTENs and may be mid-wait when ctx ends.
	conn := pooled.Hijack()
	defer conn.Close(context.Background())

	for _, channel := range channels {
		if _, err := conn.Exec(ctx, "LISTEN "+pgx.Identifier{channel}.Sanitize()); err != nil {
			return fmt.Errorf("db: listen %s: %w", channel, err)
		}
	}
	ready()

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return fmt.Errorf("db: wait for notification: %w", err)
		}
		notify(notification.Channel, notification.Payload)
	}
}
