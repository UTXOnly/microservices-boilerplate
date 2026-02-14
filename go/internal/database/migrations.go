package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSchema creates the items table if it doesn't exist.
// Run on startup so the service can initialize its own schema.
func CreateSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
		CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at DESC);
	`)
	return err
}
