package app

import (
	"time"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/value"
	"github.com/jackc/pgx/v5/pgtype"
	"uuid"
)

func timeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return dal.TimeToDB(t)
}

func timeDbToDomain(t pgtype.Timestamptz) time.Time {
	return dal.TimeFromDB(t)
}

func timeNullableDbToDomain(t pgtype.Timestamptz) Nullable[time.Time] {
	if !t.Valid {
		return Null[time.Time]()
	}
	return Value(t.Time)
}

func uuidV4() uuid.UUID {
	return uuid.NewV4()
}

func uuidDbToDomain(v pgtype.UUID) uuid.UUID {
	return uuid.UUID(v.Bytes)
}

func uuidDomainToDb(v uuid.UUID) pgtype.UUID {
	return dal.UUIDToDB(v)
}

func uuidNullableDomainToDb(v Nullable[uuid.UUID]) pgtype.UUID {
	if v.Valid {
		return uuidDomainToDb(v.Value)
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

type Nullable[T any] = value.Nullable[T]

func Null[T any]() Nullable[T]                  { return value.Null[T]() }
func Value[T any](v T) Nullable[T]              { return value.Value(v) }
func NullableFromPtr[T any](ptr *T) Nullable[T] { return value.NullableFromPtr(ptr) }

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
