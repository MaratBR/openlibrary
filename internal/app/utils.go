package app

import (
	"context"
	"encoding/json"
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

func arrUuidDomainToDb(v []uuid.UUID) []pgtype.UUID {
	uuids := make([]pgtype.UUID, len(v))
	for i := range v {
		uuids[i] = uuidDomainToDb(v[i])
	}
	return uuids
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

func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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

func (v Nullable[T]) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Value)
}

func (v *Nullable[T]) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		*v = Null[T]()
		return nil
	} else {
		err := json.Unmarshal(b, &v.Value)
		if err != nil {
			return err
		}
		v.Valid = true
		return nil
	}
}
