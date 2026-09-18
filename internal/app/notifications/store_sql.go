package notifications

import (
	"context"
	"encoding/json"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type sqlStore struct {
	db  dal.DB
	log *zap.SugaredLogger
}

// Add implements [Store].
func (s *sqlStore) Add(ctx context.Context, notification Notification) error {
	queries := store.New(s.db)
	_, err := queries.Notification_Upsert(ctx, store.Notification_UpsertParams{
		Key:       notification.Key(),
		Type:      notification.Type(),
		Content:   notification.Content,
		Title:     notification.Title,
		CreatedAt: pgtype.Timestamptz{Valid: true, Time: notification.CreatedAt},
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

// DeleteAllByUserID implements [Store].
func (s *sqlStore) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	queries := store.New(s.db)
	err := queries.Notification_DeleteByUserID(ctx, store.PgtypeUUID(userID))
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

// DeleteByKey implements [Store].
func (s *sqlStore) DeleteByKey(ctx context.Context, key string) error {
	queries := store.New(s.db)
	err := queries.Notification_DeleteByKey(ctx, key)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

// Query implements [Store].
func (s *sqlStore) Query(ctx context.Context, query StoreQuery) (NotificationsResult, error) {
	queries := store.New(s.db)

	var (
		rows []store.Notification
		err  error
	)

	if query.Cursor == 0 {
		rows, err = queries.Notification_Query(ctx, store.Notification_QueryParams{
			Limit:  query.Limit,
			UserID: store.PgtypeUUID(query.UserID),
		})
	} else if query.Before {
		rows, err = queries.Notification_QueryBefore(ctx, store.Notification_QueryBeforeParams{
			Limit:     query.Limit,
			UserID:    store.PgtypeUUID(query.UserID),
			CreatedAt: store.PgtypeTimestamptz(getTimeFromCursor(query.Cursor)),
		})
	} else {
		rows, err = queries.Notification_QueryAfter(ctx, store.Notification_QueryAfterParams{
			Limit:     query.Limit,
			UserID:    store.PgtypeUUID(query.UserID),
			CreatedAt: store.PgtypeTimestamptz(getTimeFromCursor(query.Cursor)),
		})
	}

	if err != nil {
		return NotificationsResult{}, apperror.WrapUnexpectedDBError(err)
	}

	if len(rows) == 0 {
		return NotificationsResult{NextCursor: 0}, nil
	}

	var nextCursor Cursor

	if query.Before {
		nextCursor = getCursorFromTime(store.PgtypeTimestamptzToTime(rows[0].CreatedAt))
	} else {
		nextCursor = getCursorFromTime(store.PgtypeTimestamptzToTime(rows[len(rows)-1].CreatedAt))
	}

	notifications := make([]Notification, 0, len(rows))

	for _, row := range rows {
		notification, err := s.newNotificationFromRow(row)
		if err != nil {
			s.log.Warnw("failed to instantiate a notification during Query call", "err", err, "notification_id", row.ID)
		} else {
			notifications = append(notifications, notification)
		}
	}

	return NotificationsResult{
		NextCursor:    nextCursor,
		Notifications: notifications,
	}, nil

}

func NewSQLStore(db dal.DB) Store {
	return &sqlStore{
		db: db,
	}
}

func (s *sqlStore) parseMetadata(id int64, data []byte) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal(data, &m)
	if err != nil {

		m2 := make(map[string]any)

		err = json.Unmarshal(data, &m2)
		if err != nil {
			return nil, err
		}

		finalMap := make(map[string]string)

		for key, value := range m2 {
			if str, ok := value.(string); ok {
				finalMap[key] = str
			} else {
				s.log.Debugw("ignoring key in notifiaction metadata - not a string", "notification_id", id, "meta_key", key)
			}
		}

		return finalMap, nil
	} else {
		return m, nil
	}
}

func (s *sqlStore) newNotificationFromRow(row store.Notification) (Notification, error) {
	data, err := s.parseMetadata(row.ID, row.Metadata)
	if err != nil {
		s.log.Warnw("failed to parse metadata of notification", "notification_id", row.ID, "metadata", string(row.Metadata))
	}

	return newNotification(
		row.Key,
		row.Type,
		row.Content,
		row.Title,
		store.PgtypeUUIDToDomain(row.UserID),
		data,
	)
}
