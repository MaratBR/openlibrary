package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var ErrNoRows = pgx.ErrNoRows
