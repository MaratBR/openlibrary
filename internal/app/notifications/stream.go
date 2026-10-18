package notifications

import (
	"context"
	"encoding/json"
	"github.com/MaratBR/openlibrary/internal/app/events"
	"time"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/store"
)

// Event is a persisted notification addressed to a single recipient.
type Event struct {
	ID        int64     `json:"id,string"`
	Key       string    `json:"key"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type StreamReader struct{ db store.DBTX }

func NewStreamReader(db store.DBTX) *StreamReader { return &StreamReader{db: db} }

func (s *StreamReader) Read(ctx context.Context, userID uuid.UUID, after int64) ([]Event, error) {
	rows, err := store.New(s.db).Notification_Stream(ctx, store.Notification_StreamParams{UserID: store.PgtypeUUID(userID), ID: after})
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, Event{ID: row.ID, Key: row.Key, Type: row.Type, Title: row.Title, Content: row.Content, CreatedAt: row.CreatedAt.Time})
	}
	return events, nil
}

func (s *StreamReader) History(ctx context.Context, userID uuid.UUID, before int64) ([]Event, error) {
	rows, err := store.New(s.db).Notification_History(ctx, store.Notification_HistoryParams{UserID: store.PgtypeUUID(userID), BeforeID: before})
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, Event{ID: row.ID, Key: row.Key, Type: row.Type, Title: row.Title, Content: row.Content, CreatedAt: row.CreatedAt.Time})
	}
	return events, nil
}

// Name and Poll adapt notification persistence to the shared user event stream.
func (s *StreamReader) Name() string { return "notifications" }

func (s *StreamReader) Poll(ctx context.Context, userID uuid.UUID, after int64) ([]events.StoredEvent, error) {
	records, err := s.Read(ctx, userID, after)
	if err != nil {
		return nil, err
	}
	result := make([]events.StoredEvent, 0, len(records))
	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}
		result = append(result, events.StoredEvent{Cursor: record.ID, Event: events.Event{Type: "notification", Data: data}})
	}
	return result, nil
}
