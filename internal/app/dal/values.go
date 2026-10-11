package dal

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"uuid"
)

func TimeToDB(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TimeFromDB(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		panic("received null timestamptz from database")
	}
	return t.Time
}

func UUIDToDB(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(v), Valid: true}
}
