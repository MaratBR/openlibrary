package notifications

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"unicode/utf8"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/htmlbag"
	"github.com/MaratBR/openlibrary/internal/olhttp/webcomponents"
	"github.com/a-h/templ"
)

const (
	ContentMaxLength = 1000
	TitleMaxLength   = 100
	DefaultType      = "info"
	dataKeyIcon      = "icon"
)

var (
	notificationErrors         = apperror.AppErrors.NewSubNamespace("notif")
	InvalidNotificationType    = notificationErrors.NewType("invalid_type")
	InvalidNotificationKey     = notificationErrors.NewType("invalid_key")
	invalidNotificationContent = notificationErrors.NewType("invalid_content")
	ErrContentTooLong          = invalidNotificationContent.New("content is too long")
	ErrTitleTooLong            = invalidNotificationContent.New("title is too long")
)

type Notification struct {
	userID                                uuid.UUID
	key, notificationType, Content, Title string
	Data                                  map[string]string
}

func (n Notification) Type() string {
	return n.notificationType
}

func (n Notification) Key() string {
	return n.key
}

func (n Notification) UserID() uuid.UUID {
	return n.userID
}

func newNotification(key, notificationType, content, title string, userID uuid.UUID, data map[string]string) (Notification, error) {
	if err := validateKey(key); err != nil {
		return Notification{}, err
	}

	if err := validateType(notificationType); err != nil {
		return Notification{}, err
	}

	if err := validateContent(content); err != nil {
		return Notification{}, err
	}

	if err := validateTitle(content); err != nil {
		return Notification{}, err
	}

	return Notification{
		key:              key,
		notificationType: notificationType,
		Content:          content,
		Title:            title,
		Data:             data,
		userID:           userID,
	}, nil
}

var (
	regexAlphanumeric = regexp.MustCompile("[a-zA-Z0-9]+")
)

func validateKey(key string) error {
	if key == "" {
		return InvalidNotificationKey.New("key cannot be empty")
	}

	if !regexAlphanumeric.Match([]byte(key)) {
		return InvalidNotificationKey.New("%s", "key must match "+regexAlphanumeric.String())
	}

	return nil
}

func validateType(notificationType string) error {
	if notificationType == "" {
		return InvalidNotificationType.New("type cannot be empty")
	}

	if !regexAlphanumeric.Match([]byte(notificationType)) {
		return InvalidNotificationType.New("%s", "type must match "+regexAlphanumeric.String())
	}

	return nil
}

func validateContent(content string) error {
	if len(content) > ContentMaxLength {
		return ErrContentTooLong
	}

	return nil
}

func validateTitle(title string) error {
	if utf8.RuneCountInString(title) > TitleMaxLength {
		return ErrTitleTooLong
	}

	return nil
}

func (n Notification) Render(ctx context.Context, w io.Writer) error {
	return webcomponents.Notification(
		n.Title, n.Content,
		n.icon(),
	).Render(ctx, w)
}

func (n Notification) icon() templ.Component {

	if n.Data == nil {
		return nil
	}

	icon, ok := n.Data[dataKeyIcon]
	if !ok {
		return nil
	}

	var img htmlbag.Image
	err := json.Unmarshal([]byte(icon), &img)
	if err != nil {
		return nil
	}

	return img
}
