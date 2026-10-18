package notifications

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"uuid"
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

func (b *builder) Type(value string) *builder {
	b.notificationType = value
	return b
}

func (b *builder) Category(value string) *builder  { return b.Data("category", value) }
func (b *builder) Urgent(value bool) *builder      { return b.Data("urgent", strconv.FormatBool(value)) }
func (b *builder) IconEmoji(emoji string) *builder { return b.icon(IconEmoji, emoji) }

func (b *builder) IconImg(url string) *builder {
	return b.icon(IconImage, url)
}

func (b *builder) IconFontAwesome(classes string) *builder {
	return b.icon(IconFontAwesome, classes)
}

func (b *builder) icon(kind byte, payload string) *builder {
	data, _ := json.Marshal(Icon{Type: kind, Payload: payload})
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

	var icon *Icon
	if raw, ok := b.data[dataKeyIcon]; ok {
		icon = &Icon{}
		if err := json.Unmarshal([]byte(raw), icon); err != nil {
			return Notification{}, err
		}
	}
	if _, err := (Options{Category: b.data["category"], Icon: icon}).metadata(); err != nil {
		return Notification{}, fmt.Errorf("%w: %v", ErrInvalidSend, err)
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
