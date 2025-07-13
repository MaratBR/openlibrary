package flash

import "context"

type FlashBackend interface {
	PullAll(ctx context.Context, id string) ([]Message, error)
	Append(ctx context.Context, id string, message Message)
}
