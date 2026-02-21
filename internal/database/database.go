package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps pgx connection pool. Use for all database access.
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a connection pool from connStr (e.g. DATABASE_URL).
// Call once at startup; log.Fatal on failure.
func New(ctx context.Context, connStr string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	log.Println("database: connected")
	return &DB{Pool: pool}, nil
}

// Close closes the connection pool. Call on graceful shutdown.
func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
		log.Println("database: closed")
	}
}
