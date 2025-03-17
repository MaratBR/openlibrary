package olresponse

import "encoding/json"

type NotificationType uint8

func (n NotificationType) MarshalJSON() ([]byte, error) {
	var v string

	switch n {
	case NotificationInfo:
		v = "info"
		break
	case NotificationError:
		v = "error"
		break
	}

	return json.Marshal(v)
}

const (
	NotificationInfo NotificationType = iota
	NotificationWarn
	NotificationError
)

type Notification struct {
	Text string           `json:"text"`
	Type NotificationType `json:"type"`
}

func NewNotification(text string, typ NotificationType) Notification {
	return Notification{
		Type: typ,
		Text: text,
	}
}
