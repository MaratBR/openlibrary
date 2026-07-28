package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/commonutil"
)

type EventProcessor interface {
	Process(ctx context.Context, events []Event) error
}

type simpleEventProcessor struct {
	metricSink MetricSink
}

// Process implements [EventProcessor].
func (s *simpleEventProcessor) Process(ctx context.Context, events []Event) error {
	metrics := extractMetrics(events)
	s.metricSink.SubmitMetrics(ctx, metrics)
	return nil
}

func extractMetrics(events []Event) []MetricRecord {
	metrics := commonutil.MapSlice(events, func(ev Event) MetricRecord {
		return MetricRecord{
			Type:    MetricType(ev.EventType),
			Value:   ev.Value,
			Samples: 1,
			BookID:  ev.BookID,
		}
	})

	return metrics

}

func newSimpleEventProcessor(metricSink MetricSink) EventProcessor {
	return &simpleEventProcessor{metricSink: metricSink}
}
