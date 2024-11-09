package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
)

type txFactory interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DB interface {
	store.DBTX
	txFactory
}
