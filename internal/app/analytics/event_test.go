package analytics

import (
	"net"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

func TestEventToMetricPreservesFields(t *testing.T) {
	createdAt := time.Date(2026, time.March, 10, 23, 59, 0, 0, time.UTC)
	event := newEvent(42, "visitor", EventTypeView, 2.5, createdAt)

	metric := event.ToMetric()
	if metric.Type != MetricViews || metric.BookID != 42 || metric.Value != 2.5 || metric.Samples != 1 {
		t.Fatalf("ToMetric() = %#v", metric)
	}
	if !metric.OccurredAt.Equal(createdAt) {
		t.Errorf("OccurredAt = %s, want %s", metric.OccurredAt, createdAt)
	}
	if got, want := event.Key(), "view:42:visitor"; got != want {
		t.Errorf("Key() = %q, want %q", got, want)
	}
}

func TestEventMetadataUniqueID(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	tests := []struct {
		name string
		meta EventMetadata
		want string
	}{
		{
			name: "user takes precedence over IP",
			meta: EventMetadata{UserID: uuid.NullUUID{UUID: userID, Valid: true}, IP: net.ParseIP("192.0.2.1")},
			want: "U" + userID.String(),
		},
		{name: "anonymous IP", meta: EventMetadata{IP: net.ParseIP("192.0.2.1")}, want: "192.0.2.1"},
		{name: "unknown", meta: EventMetadata{}, want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.meta.UniqueID(); got != tt.want {
				t.Errorf("UniqueID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewBookViewEvent(t *testing.T) {
	before := time.Now()
	event := NewBookViewEvent(17, EventMetadata{IP: net.ParseIP("192.0.2.2")})
	after := time.Now()

	if event.BookID != 17 || event.EventType != EventTypeView || event.Value != 1 || event.UserKey != "192.0.2.2" {
		t.Fatalf("NewBookViewEvent() = %#v", event)
	}
	if event.CreatedAt.Before(before) || event.CreatedAt.After(after) {
		t.Errorf("CreatedAt = %s, want between %s and %s", event.CreatedAt, before, after)
	}
}
