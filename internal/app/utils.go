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
	if !t.Valid {
		panic("received null timestamptz from database")
	}
	return t.Time
}

func timeNullableDbToDomain(t pgtype.Timestamptz) Nullable[time.Time] {
	if !t.Valid {
		return Null[time.Time]()
	}
	return Value(t.Time)
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

func uuidNullableDomainToDb(v uuid.NullUUID) pgtype.UUID {
	if v.Valid {
		return uuidDomainToDb(v.UUID)
	}
	return pgtype.UUID{Valid: false}
}

func arrUuidDomainToDb(v []uuid.UUID) []pgtype.UUID {
	uuids := make([]pgtype.UUID, len(v))
	for i := range v {
		uuids[i] = uuidDomainToDb(v[i])
	}
	return uuids
}

func arrUuidDomainToString(v []uuid.UUID) []string {
	uuids := make([]string, len(v))
	for i := range v {
		uuids[i] = v[i].String()
	}
	return uuids
}

func arrInt64ToInt64String(v []int64) []Int64String {
	ints := make([]Int64String, len(v))
	for i := range v {
		ints[i] = Int64String(v[i])
	}
	return ints
}

func ArrInt64StringToInt64(v []Int64String) []int64 {
	ints := make([]int64, len(v))
	for i := range v {
		ints[i] = int64(v[i])
	}
	return ints
}

func int64NullableDomainToDb(v Nullable[int64]) pgtype.Int8 {
	return pgtype.Int8{Valid: v.Valid, Int64: v.Value}
}

func rollbackTx(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		slog.Error("failed to rollback transaction", "err", err)
	}
}

func MapSlice[T, U any](ts []T, f func(T) U) []U {
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

func NullableFromPtr[T any](ptr *T) Nullable[T] {
	if ptr == nil {
		return Null[T]()
	}
	return Value(*ptr)
}

func (v Nullable[T]) Or(vv T) T {
	if v.Valid {
		return v.Value
	}
	return vv
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

func (v Nullable[T]) Ptr() *T {
	if !v.Valid {
		return nil
	}
	return &v.Value
}

func int64ToNullable(v pgtype.Int8) Nullable[Int64String] {
	if v.Valid {
		return Value(Int64String(v.Int64))
	} else {
		return Null[Int64String]()
	}
}

func float64ToNullable(v pgtype.Float8) Nullable[float64] {
	if v.Valid {
		return Value(v.Float64)
	} else {
		return Null[float64]()
	}
}

func moveArrEl[T any](arr []T, oldIndex, newIndex int) {
	if oldIndex < 0 {
		panic("oldIndex < 0")
	}
	if oldIndex >= len(arr) {
		panic("oldIndex >= len(arr)")
	}
	if newIndex < 0 {
		panic("newIndex < 0")
	}
	if newIndex >= len(arr) {
		panic("newIndex >= len(arr)")
	}

	if oldIndex == newIndex {
		return
	}

	el := arr[oldIndex]
	if oldIndex > newIndex {
		// 1 2 3 4 5 6 7 8
		//           ^ - oldIndex
		//   ^ - newIndex
		// move all elements before by one
		for i := oldIndex; i > newIndex; i-- {
			arr[i] = arr[i-1]
		}
	} else {
		// 1 2 3 4 5 6 7 8
		//   ^ - oldIndex
		//           ^ - newIndex
		for i := oldIndex; i < newIndex; i++ {
			arr[i] = arr[i+1]
		}
	}
	arr[newIndex] = el
}
