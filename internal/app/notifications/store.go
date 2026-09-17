package notifications

import (
	"context"
	"uuid"
)

type NotificationsResult struct {
	Notifications []Notification
	NextCursor    int64
}

type Store interface {
	Query(ctx context.Context, userID uuid.UUID, cursor int64) error
	Add(ctx context.Context, notification Notification) error
	DeleteByKey(ctx context.Context, key string) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}
