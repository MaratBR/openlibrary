package app

import (
	"context"
)

type BookReindexService interface {
	Reindex(ctx context.Context, id int64) error
	ScheduleReindex(ctx context.Context, id int64)
	ScheduleReindexAll() error
}
