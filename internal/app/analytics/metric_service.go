package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
)

type MetricService interface {
	Get(ctx context.Context, metric string, bookIds []int64) (map[int64]MetricValues, error)
	GetMetricValue(ctx context.Context, metric string, bookId int64, period store.OlAnalyticsBucketPeriodType) (MetricValue, error)
}

type dummyMetricService struct{}

// Get implements [MetricService].
func (d *dummyMetricService) Get(ctx context.Context, metric string, bookIds []int64) (map[int64]MetricValues, error) {
	return make(map[int64]MetricValues), nil
}

// GetMetricValue implements [MetricService].
func (d *dummyMetricService) GetMetricValue(ctx context.Context, metric string, bookId int64, period store.OlAnalyticsBucketPeriodType) (MetricValue, error) {
	return MetricValue{}, nil
}

func NewDummyMetricService() MetricService {
	return &dummyMetricService{}
}
