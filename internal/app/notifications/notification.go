package notifications

import (
	"maps"
	"regexp"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
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
	CreatedAt                             time.Time
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

func (n Notification) Category() string {
	category, _, _ := presentation(n.Data)
	return category
}

func (n Notification) Urgent() bool {
	_, urgent, _ := presentation(n.Data)
	return urgent
}

func (n Notification) Icon() *Icon {
	_, _, icon := presentation(n.Data)
	return icon
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

	if err := validateTitle(title); err != nil {
		return Notification{}, err
	}

	if data == nil {
		data = make(map[string]string)
	}

	return Notification{
		key:              key,
		notificationType: notificationType,
		Content:          SanitizeContent(content),
		Title:            title,
		Data:             maps.Clone(data),
		userID:           userID,
		CreatedAt:        time.Now(),
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

	if utf8.RuneCountInString(notificationType) > 100 {
		return InvalidNotificationType.New("type is too long")
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
