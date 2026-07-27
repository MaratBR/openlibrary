package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type Event struct {
	BookID    int64
	UserKey   string
	EventType store.OlAnalyticsInteractionEventType
	Value     float64
	CreatedAt time.Time
}

func (ev Event) Key() string {
	return fmt.Sprintf("%s:%d:%s", ev.EventType, ev.BookID, ev.UserKey)
}

func NewBookViewEvent(
	bookID int64,
	meta ViewMetadata,
) Event {
	return Event{
		BookID:    bookID,
		UserKey:   meta.UniqueID(),
		CreatedAt: time.Now(),
		EventType: store.OlAnalyticsInteractionEventTypeBookView,
	}
}

type EventRepository interface {
	Insert(ctx context.Context, events []Event) error
}

type EventSink interface {
	SubmitEvent(ctx context.Context, event Event)
}
