package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

type eventRepository struct {
	db  store.DBTX
	log *zap.SugaredLogger
}

func newEventRepository(db store.DBTX, log *zap.SugaredLogger) EventRepository {
	return &eventRepository{db: db, log: log}
}

func (r *eventRepository) Insert(ctx context.Context, events []Event) error {
	queries := store.New(r.db)
	rows := make([]store.Analytics_InsertEventParams, len(events))

	for i, event := range events {
		rows[i] = store.Analytics_InsertEventParams{
			EventType: string(event.EventType),
			UserKey:   event.UserKey,
			BookID:    event.BookID,
			Value:     event.Value,
		}
	}

	inserted, err := queries.Analytics_InsertEvent(ctx, rows)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	if inserted != int64(len(rows)) {
		r.log.Warnw("mismatch in number of inserted rows", "inserted", inserted, "rowsCount", len(rows))
	}

	return nil
}
