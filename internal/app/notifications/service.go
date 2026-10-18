package notifications

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/store"
)

// Service owns notification delivery, history and recipient read state. Callers
// never access the notification user_data key directly.
type Service interface {
	History(context.Context, uuid.UUID, int64) ([]Event, error)
	Status(context.Context, uuid.UUID) (Status, error)
	MarkRead(context.Context, uuid.UUID, time.Time) (Status, error)
	Send(context.Context, string, string, string, string, ...Options) (int64, error)
}

var _ Service = (*StreamReader)(nil)
var ErrInvalidReadTime = errors.New("invalid notification read timestamp")

type Status struct {
	UnreadCount int64      `json:"unreadCount"`
	LastReadAt  *time.Time `json:"lastReadAt"`
	ObservedAt  time.Time  `json:"observedAt"`
}

func (s *StreamReader) Status(ctx context.Context, userID uuid.UUID) (Status, error) {
	row, err := store.New(s.db).Notification_Status(ctx, store.PgtypeUUID(userID))
	if err != nil {
		return Status{}, err
	}
	result := Status{UnreadCount: row.UnreadCount, ObservedAt: row.ObservedAt.Time}
	if row.LastReadAt.Valid {
		result.LastReadAt = &row.LastReadAt.Time
	}
	return result, nil
}

func (s *StreamReader) MarkRead(ctx context.Context, userID uuid.UUID, through time.Time) (Status, error) {
	if through.IsZero() || through.Before(time.Unix(0, 0)) {
		return Status{}, ErrInvalidReadTime
	}
	err := store.New(s.db).Notification_MarkRead(ctx, store.Notification_MarkReadParams{
		UserID: store.PgtypeUUID(userID), ReadThrough: store.PgtypeTimestamptz(through),
	})
	if err != nil {
		return Status{}, err
	}
	return s.Status(ctx, userID)
}
