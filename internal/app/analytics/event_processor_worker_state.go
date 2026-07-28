package analytics

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	eventProcessorStateName = "event_processor"
)

type eventProcessorWorkerState interface {
	GetCurrentCursor(ctx context.Context) (int64, error)
	SaveCurrentCursor(ctx context.Context, cursor int64) error
}

type eventProcessorWorkerStateImpl struct {
	queries *store.Queries
}

// GetCurrentCursor implements [eventProcessorWorkerState].
func (e *eventProcessorWorkerStateImpl) GetCurrentCursor(ctx context.Context) (int64, error) {
	state, err := e.queries.Analytics_GetWorkerState(ctx, eventProcessorStateName)
	if err != nil {
		if err == store.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return state.LastCursor, nil

}

// SaveCurrentCursor implements [eventProcessorWorkerState].
func (e *eventProcessorWorkerStateImpl) SaveCurrentCursor(ctx context.Context, cursor int64) error {
	err := e.queries.Analytics_SetWorkerState(ctx, store.Analytics_SetWorkerStateParams{
		WorkerName: eventProcessorStateName,
		LastLaunch: pgtype.Timestamptz{Valid: true, Time: time.Now()},
		LastCursor: cursor,
	})
	return err
}

func newEventProcessorWorkerState(db store.DBTX) eventProcessorWorkerState {
	return &eventProcessorWorkerStateImpl{queries: store.New(db)}
}
