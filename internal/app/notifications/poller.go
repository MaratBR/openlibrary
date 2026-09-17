package notifications

import (
	"context"
	"time"
)

type Poller interface {
	Poll(ctx context.Context, lastPoll time.Time, maxItems int) ([]Notification, error)
}
