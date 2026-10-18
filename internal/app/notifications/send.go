package notifications

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidSend = errors.New("invalid notification")

// Send persists a notification for a user or a snapshot of all users with a role.
// One INSERT statement makes role broadcasts atomic.
func (s *StreamReader) Send(ctx context.Context, userID, role, title, content string) (int64, error) {
	if (userID == "") == (role == "") {
		return 0, fmt.Errorf("%w: choose exactly one recipient or role", ErrInvalidSend)
	}
	if err := validateTitle(title); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidSend, err)
	}
	if title == "" {
		return 0, fmt.Errorf("%w: title is required", ErrInvalidSend)
	}
	if err := validateContent(content); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidSend, err)
	}
	params := store.Notification_SendParams{Key: newGenericNotificationKey(), Title: title, Content: content}
	if userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return 0, fmt.Errorf("%w: invalid recipient ID", ErrInvalidSend)
		}
		params.RecipientID = store.PgtypeUUID(id)
	} else {
		switch role {
		case "user", "admin", "moderator", "system":
		default:
			return 0, fmt.Errorf("%w: invalid role", ErrInvalidSend)
		}
		params.RecipientRole = pgtype.Text{String: role, Valid: true}
	}
	return store.New(s.db).Notification_Send(ctx, params)
}
