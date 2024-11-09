package app

import (
	"encoding/json"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type Int64String int64

func (i Int64String) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

func (i *Int64String) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64String(v)
	return nil
}

func int64ToNullable(v pgtype.Int8) Nullable[Int64String] {
	if v.Valid {
		return Value(Int64String(v.Int64))
	} else {
		return Null[Int64String]()
	}
}
