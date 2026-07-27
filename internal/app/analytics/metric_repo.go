package analytics

import "context"

type MetricRepository interface {
	AddMetrics(ctx context.Context, metric []Metric) error
}
