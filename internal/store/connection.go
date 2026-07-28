package store

import (
	"context"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	// Query parameters are deliberately excluded; they frequently contain user
	// input. SQL text and operation names retain enough detail to find hotspots.
	config.ConnConfig.Tracer = otelpgx.NewTracer()
	return pgxpool.NewWithConfig(ctx, config)
}

var ErrNoRows = pgx.ErrNoRows
