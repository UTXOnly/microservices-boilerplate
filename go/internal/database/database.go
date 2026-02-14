package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool wraps pgxpool.Pool for dependency injection.
type Pool struct {
	*pgxpool.Pool
}

// NewPool creates a new connection pool for async/concurrent Postgres access.
func NewPool(ctx context.Context, connString string) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Pool sizing for concurrent requests
	cfg.MaxConns = 25
	cfg.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Pool{Pool: pool}, nil
}

// Close closes the connection pool.
func (p *Pool) Close() {
	p.Pool.Close()
}
