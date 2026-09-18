package notifications

import "context"

type Sink interface {
	Push(ctx context.Context, notification Notification) error
}
