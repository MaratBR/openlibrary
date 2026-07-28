package analytics

import "context"

type MetricRepository interface {
	AddMetrics(ctx context.Context, metric []MetricRecord) error
}
