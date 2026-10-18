package events

import (
	"encoding/json"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

func TestSendIsRecipientScopedAndSupportsAnyPayload(t *testing.T) {
	hub := NewHub()
	alice, bob := uuid.NewV7(), uuid.NewV7()
	a, stopA := hub.Subscribe(alice)
	defer stopA()
	a2, stopA2 := hub.Subscribe(alice)
	defer stopA2()
	b, stopB := hub.Subscribe(bob)
	defer stopB()
	require.NoError(t, hub.Send(alice, "book.progress", map[string]any{"bookId": 123, "progress": 0.5}))
	for _, ch := range []<-chan Event{a, a2} {
		event := <-ch
		require.Equal(t, "book.progress", event.Type)
		require.JSONEq(t, `{"bookId":123,"progress":0.5}`, string(event.Data))
	}
	select {
	case <-b:
		t.Fatal("event leaked to another user")
	default:
	}
}
func TestSlowConsumersAndUnsubscribeReleaseSubscriptions(t *testing.T) {
	hub := NewHub()
	id := uuid.NewV7()
	ch, stop := hub.Subscribe(id)
	for i := 0; i < 65; i++ {
		require.NoError(t, hub.Send(id, "test", i))
	}
	count := 0
	for range ch {
		count++
	}
	require.Equal(t, 64, count)
	stop()
	stop()
	require.Empty(t, hub.(*inMemoryHub).users)
}
func TestSendRejectsInvalidEventContracts(t *testing.T) {
	hub := NewHub()
	id := uuid.NewV7()
	for _, typ := range []string{"", "x\ndata: injection", "x y"} {
		require.Error(t, hub.Send(id, typ, nil))
	}
	require.Error(t, hub.Send(id, "test", json.RawMessage("invalid")))
}
