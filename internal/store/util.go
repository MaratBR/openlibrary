package store

import (
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeUUID(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Valid: true, Bytes: uuid}
}

func PgtypeUUIDToDomain(pgTypeUUID pgtype.UUID) uuid.UUID {
	return uuid.UUID(pgTypeUUID.Bytes)
}

func PgtypeTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Valid: true, Time: t}
}

func PgtypeTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	return t.Time
}
