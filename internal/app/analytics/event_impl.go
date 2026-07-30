package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
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
			CreatedAt: pgtype.Timestamptz{Valid: true, Time: event.CreatedAt},
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

func (r *eventRepository) GetEvents(ctx context.Context, maxCount int, cursor int64) (int64, []Event, error) {
	queries := store.New(r.db)
	rows, err := queries.Analytics_GetEvents(ctx, store.Analytics_GetEventsParams{
		Limit: int32(maxCount),
		ID:    cursor,
	})
	if err != nil {
		return cursor, nil, apperror.WrapUnexpectedDBError(err)
	}

	if len(rows) == 0 {
		return cursor, nil, nil
	}

	events := commonutil.MapSlice(rows, func(ev store.OlAnalyticsInteractionEvent) Event {
		return newEvent(ev.BookID, ev.UserKey, EventType(ev.EventType), ev.Value, ev.CreatedAt.Time)
	})

	newCursor := rows[len(rows)-1].ID

	return newCursor, events, nil
}
