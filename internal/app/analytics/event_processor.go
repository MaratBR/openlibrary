package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/commonutil"
	"go.opentelemetry.io/otel/attribute"
)

type EventProcessor interface {
	Process(ctx context.Context, events []Event) error
}

type simpleEventProcessor struct {
	metricSink MetricSink
}

// Process implements [EventProcessor].
func (s *simpleEventProcessor) Process(ctx context.Context, events []Event) error {
	ctx, span := tracer.Start(ctx, "simpleEventProcessor.Process")
	span.SetAttributes(attribute.Int("analytics.event_count", len(events)))
	defer span.End()
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
