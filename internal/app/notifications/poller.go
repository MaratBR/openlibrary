package notifications

import (
	"context"
	"time"
)

type PollerResult struct {
	PolledAt      time.Time
	Notifications []Notification
}

type Poller interface {
	Poll2(ctx context.Context, lastPoll time.Time, maxItems int) (PollerResult, error)
	Poll(ctx context.Context, maxItems int) (PollerResult, error)
}
