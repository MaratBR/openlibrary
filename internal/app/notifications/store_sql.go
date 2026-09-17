package notifications

import (
	"context"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/dal"
)

type sqlStore struct {
	db dal.DB
}

// Add implements [Store].
func (s *sqlStore) Add(ctx context.Context, notification Notification) error {
	panic("unimplemented")
}

// DeleteAllByUserID implements [Store].
func (s *sqlStore) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

// DeleteByKey implements [Store].
func (s *sqlStore) DeleteByKey(ctx context.Context, key string) error {
	panic("unimplemented")
}

// Query implements [Store].
func (s *sqlStore) Query(ctx context.Context, userID uuid.UUID, cursor int64) error {
	panic("unimplemented")
}

func NewSqlStore(db dal.DB) Store {
	return &sqlStore{
		db: db,
	}
}
