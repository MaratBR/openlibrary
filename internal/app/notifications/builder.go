package notifications

import (
	"strings"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/htmlbag"
)

type builder struct {
	key, notificationType, title, content string

	data   map[string]string
	userID uuid.UUID
}

func NewBuilder(title string, userID uuid.UUID) *builder {
	return &builder{title: title}
}

func (b *builder) Content(content string) *builder {
	b.content = strings.Trim(content, "\t\n\r ")
	return b
}

func (b *builder) Data(key, value string) *builder {
	if b.data == nil {
		b.data = map[string]string{
			key: value,
		}
	} else {
		b.data[key] = value
	}
	return b
}

func (b *builder) IconImg(url string) *builder {
	return b.icon(htmlbag.NewURLImage(url))
}

func (b *builder) IconFontAwesome(classes string) *builder {
	return b.icon(htmlbag.NewFontAwesomeImage(classes))
}

func (b *builder) icon(icon htmlbag.Image) *builder {
	text := icon.MarshalText()
	return b.Data(dataKeyIcon, text)
}

func (b builder) Build() (Notification, error) {
	key := b.key
	if key == "" {
		key = newGenericNotificationKey()
	}

	notificationType := b.notificationType
	if notificationType == "" {
		notificationType = DefaultType
	}

	return newNotification(
		key,
		notificationType,
		b.content,
		b.title,
		b.userID,
		b.data,
	)
}

func newGenericNotificationKey() string {
	id := uuid.NewV7()
	return "g-" + id.String()
}
