package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func timeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func timeDbToDomain(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func uuidV4() uuid.UUID {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return u
}

func uuidDbToDomain(v pgtype.UUID) uuid.UUID {
	return uuid.UUID(v.Bytes)
}
func uuidDomainToDb(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(v), Valid: true}
}

func rollbackTx(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		slog.Error("failed to rollback transaction", "err", err)
	}
}

func mapSlice[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

type Nullable[T any] struct {
	Value T
	Valid bool
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

func Value[T any](v T) Nullable[T] {
	return Nullable[T]{Value: v, Valid: true}
}
