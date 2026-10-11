package value

import "encoding/json"

// Nullable represents an optional value and encodes as that value or JSON null.
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
