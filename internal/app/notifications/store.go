package notifications

import (
	"context"
	"uuid"
)

type NotificationsResult struct {
	Notifications []Notification
	NextCursor    Cursor
}

type StoreQuery struct {
	Cursor Cursor
	Limit  int32
	Before bool
	UserID uuid.UUID
}

type Store interface {
	Query(ctx context.Context, query StoreQuery) (NotificationsResult, error)
	Add(ctx context.Context, notification Notification) error
	DeleteByKey(ctx context.Context, key string) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}
