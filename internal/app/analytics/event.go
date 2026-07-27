package analytics

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type EventType string

const (
	EventTypeView = "view"
)

type Event struct {
	BookID    int64
	UserKey   string
	EventType EventType
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
		EventType: EventTypeView,
	}
}

type EventRepository interface {
	Insert(ctx context.Context, events []Event) error
}

type EventSink interface {
	SubmitEvent(ctx context.Context, event Event)
}

var eventModule = fx.Module("ol_analytics_events",
	fx.Provide(
		newEventRepository,
		newEventBackgroundService,
		fx.Private,
	),
	fx.Provide(newDedupedEventSink),

	fx.Invoke(
		func(*eventBackgroundService) {},
	),
)
