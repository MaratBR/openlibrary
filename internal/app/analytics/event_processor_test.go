package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

type recordingMetricRepository struct {
	err     error
	records []MetricRecord
	tx      store.DBTX
}

func (r *recordingMetricRepository) AddMetrics(_ context.Context, tx store.DBTX, records []MetricRecord) error {
	r.tx = tx
	r.records = append([]MetricRecord(nil), records...)
	return r.err
}

func TestSimpleEventProcessorConvertsEvents(t *testing.T) {
	repo := &recordingMetricRepository{}
	processor := newSimpleEventProcessor(repo, zap.NewNop().Sugar())
	createdAt := time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)
	events := []Event{
		newEvent(10, "a", EventTypeView, 1, createdAt),
		newEvent(20, "b", EventType(MetricSearchClicks), 3, createdAt.Add(time.Minute)),
	}

	if err := processor.Process(context.Background(), nil, events); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if len(repo.records) != 2 {
		t.Fatalf("Process() stored %d metrics, want 2", len(repo.records))
	}
	for i, event := range events {
		metric := repo.records[i]
		if metric.BookID != event.BookID || metric.Type != MetricType(event.EventType) || metric.Value != event.Value || metric.Samples != 1 || !metric.OccurredAt.Equal(event.CreatedAt) {
			t.Errorf("metric[%d] = %#v, event = %#v", i, metric, event)
		}
	}
}

func TestSimpleEventProcessorPropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repo := &recordingMetricRepository{err: wantErr}
	processor := newSimpleEventProcessor(repo, zap.NewNop().Sugar())

	err := processor.Process(context.Background(), nil, []Event{newEvent(1, "a", EventTypeView, 1, time.Now())})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Process() error = %v, want %v", err, wantErr)
	}
}
