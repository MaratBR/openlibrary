package events

import (
	"context"
	"uuid"
)

// Source adapts a feature's persisted records into recipient-scoped events.
// Each source has its own cursor namespace; live Hub events need no cursor.
type Source interface {
	Name() string
	Poll(context.Context, uuid.UUID, int64) ([]StoredEvent, error)
}

type StoredEvent struct {
	Cursor int64
	Event  Event
}
