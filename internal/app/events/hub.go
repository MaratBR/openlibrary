// Package events provides live, recipient-scoped application events. Features
// that need durable delivery (such as notifications) persist their own records.
package events

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"uuid"
)

// Event is the transport envelope. Type identifies the feature's event contract.
// Live events have no resume ID; durable feature sources supply IDs separately.
type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

var eventType = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]{0,99}$`)

// Hub delivers live events to active subscriptions for a user.
type Hub interface {
	Subscribe(userID uuid.UUID) (<-chan Event, func())
	Send(userID uuid.UUID, typ string, data any) error
}

type inMemoryHub struct {
	mu    sync.Mutex
	users map[uuid.UUID]map[chan Event]struct{}
}

func NewHub() Hub { return &inMemoryHub{users: make(map[uuid.UUID]map[chan Event]struct{})} }

// Subscribe creates an independent bounded subscription for a user's browser tab.
// A slow consumer is disconnected instead of blocking application work.
func (h *inMemoryHub) Subscribe(userID uuid.UUID) (<-chan Event, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan Event, 64)
	if h.users[userID] == nil {
		h.users[userID] = make(map[chan Event]struct{})
	}
	h.users[userID][ch] = struct{}{}
	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if _, ok := h.users[userID][ch]; ok {
			delete(h.users[userID], ch)
			close(ch)
		}
		if len(h.users[userID]) == 0 {
			delete(h.users, userID)
		}
	}
}

// Send serializes arbitrary feature data and delivers it only to this user's
// active connections. Live events are process-local and are not replayed.
func (h *inMemoryHub) Send(userID uuid.UUID, typ string, data any) error {
	if !eventType.MatchString(typ) {
		return fmt.Errorf("invalid event type")
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.users[userID] {
		select {
		case ch <- Event{Type: typ, Data: encoded}:
		default:
			delete(h.users[userID], ch)
			close(ch)
		}
	}
	if len(h.users[userID]) == 0 {
		delete(h.users, userID)
	}
	return nil
}
