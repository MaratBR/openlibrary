package notifications

import (
	"strings"
	"uuid"

	"encoding/json"
)

type builder struct {
	key, notificationType, title, content string

	data   map[string]string
	userID uuid.UUID
}

func NewBuilder(title string, userID uuid.UUID) *builder {
	return &builder{title: title, userID: userID}
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
	return b.icon(0, url)
}

func (b *builder) IconFontAwesome(classes string) *builder {
	return b.icon(1, classes)
}

func (b *builder) icon(kind byte, payload string) *builder {
	data, _ := json.Marshal(struct {
		Type    byte   `json:"type"`
		Payload string `json:"payload"`
	}{kind, payload})
	return b.Data(dataKeyIcon, string(data))
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
