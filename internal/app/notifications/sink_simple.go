package notifications

import "context"

type simpleSink struct {
	store Store
}

// Push implements [Sink].
func (s *simpleSink) Push(ctx context.Context, notification Notification) error {
	return s.store.Add(ctx, notification)
}

func NewSimpleSink(store Store) Sink {
	return &simpleSink{store: store}
}
