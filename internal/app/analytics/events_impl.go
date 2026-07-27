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
			EventType: event.EventType,
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

	err = queries.Analytics_UpdatePopularity(ctx, store.Analytics_UpdatePopularityParams{
		BookViewScore:        1,
		SearchClickScore:     1,
		ChapterViewScore:     0.5,
		StartedReadingScore:  2,
		CompletedScore:       3,
		DroppedScore:         -2.5,
		FinishedChapterScore: 0.3,
		HalfLifeSeconds:      7 * 24 * 60 * 60,
	})
	if err != nil {
		r.log.Errorw("failed to call Analytics_UpdatePopularity", "err", err)
	}

	return nil
}
