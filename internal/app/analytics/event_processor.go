package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type EventProcessor interface {
	Process(ctx context.Context, tx store.DBTX, events []Event) error
}

type simpleEventProcessor struct {
	metricRepo MetricRepository
	log        *zap.SugaredLogger
}

// Process implements [EventProcessor].
func (s *simpleEventProcessor) Process(ctx context.Context, tx store.DBTX, events []Event) error {
	ctx, span := tracer.Start(ctx, "simpleEventProcessor.Process")
	span.SetAttributes(attribute.Int("analytics.event_count", len(events)))
	defer span.End()
	metrics := extractMetrics(events)

	err := s.metricRepo.AddMetrics(ctx, tx, metrics)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func extractMetrics(events []Event) []MetricRecord {
	metrics := commonutil.MapSlice(events, func(ev Event) MetricRecord {
		return ev.ToMetric()
	})

	return metrics

}

func newSimpleEventProcessor(metricRepo MetricRepository, log *zap.SugaredLogger) EventProcessor {
	return &simpleEventProcessor{metricRepo: metricRepo, log: log}
}
