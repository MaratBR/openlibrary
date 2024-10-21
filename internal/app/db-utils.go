package app

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func timeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
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
