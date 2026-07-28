package analytics

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gofrs/uuid"
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

func newEvent(bookID int64, userKey string, eventType EventType, value float64, createdAt time.Time) Event {
	return Event{
		BookID:    bookID,
		UserKey:   userKey,
		EventType: eventType,
		Value:     value,
		CreatedAt: createdAt,
	}
}

func (ev Event) Key() string {
	return fmt.Sprintf("%s:%d:%s", ev.EventType, ev.BookID, ev.UserKey)
}

func NewBookViewEvent(
	bookID int64,
	meta EventMetadata,
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
	GetEvents(ctx context.Context, maxCount int, cursor int64) (int64, []Event, error)
}

type EventSink interface {
	SubmitEvent(ctx context.Context, event Event)
}

type EventMetadata struct {
	UserID uuid.NullUUID
	IP     net.IP
}

func (m EventMetadata) UniqueID() string {
	if m.UserID.Valid {
		return "U" + m.UserID.UUID.String()
	}

	if m.IP != nil {
		return m.IP.String()
	}

	return "unknown"
}

var eventModule = fx.Module("ol_analytics_events",
	fx.Provide(
		newEventRepository,
		newEventBackgroundService,
		newSimpleEventProcessor,
		newEventProcessorWorkerState,
		newEventProcessorWorker,
		fx.Private,
	),
	fx.Provide(newDedupedEventSink, newEventProcessorWorkerHandle),

	fx.Invoke(
		func(*eventBackgroundService, *eventProcessorWorker) {},
	),
)
