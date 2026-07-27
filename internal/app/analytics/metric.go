package analytics

import (
	"context"

	"go.uber.org/fx"
)

const (
	MetricViews        = "views"
	MetricSearchClicks = "searchclick"
	MetricImpressions  = "impression"
)

type Metric struct {
	Type   string
	Value  float64
	BookID int64
}

func NewMetric(bookID int64, name string) Metric {
	return Metric{
		BookID: bookID,
		Type:   MetricViews,
	}
}

type MetricSink interface {
	SubmitMetric(ctx context.Context, metric Metric)
}

var metricModule = fx.Module("ol_analytics",
	fx.Provide(
		newMetricRepository,
		newMetricBackgroundService,
		fx.Private,
	),

	fx.Provide(
		newMetricSink,
		newMetricService,
	),

	fx.Invoke(func(*metricBackgroundService) {}),
)
