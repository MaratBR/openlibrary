package notifications

import (
	"context"
	"time"
)

type sqlPoller struct {
	store Store
}

// Poll implements [Poller].
func (s *sqlPoller) Poll(ctx context.Context, maxItems int) (PollerResult, error) {
	polledAt := time.Now()
	result, err := s.store.Query(ctx, StoreQuery{
		Limit: int32(maxItems),
	})
	if err != nil {
		return PollerResult{}, err
	}
	return PollerResult{
		Notifications: result.Notifications,
		PolledAt:      polledAt,
	}, nil
}

// Poll2 implements [Poller].
func (s *sqlPoller) Poll2(ctx context.Context, lastPoll time.Time, maxItems int) (PollerResult, error) {
	polledAt := time.Now()
	result, err := s.store.Query(ctx, StoreQuery{
		Limit:  int32(maxItems),
		Before: true,
		Cursor: getCursorFromTime(lastPoll),
	})
	if err != nil {
		return PollerResult{}, err
	}
	return PollerResult{
		Notifications: result.Notifications,
		PolledAt:      polledAt,
	}, nil
}

func NewSQLPoller(store Store) Poller {
	return &sqlPoller{store: store}
}
