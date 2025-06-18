package app

import (
	"context"
)

type BookReindexService interface {
	Reindex(ctx context.Context, id int64) error
	ScheduleReindexAll() error
}
