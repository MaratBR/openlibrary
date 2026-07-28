package analytics

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type GetTopBooksByMetricQuery struct {
	Metric      MetricType
	Period      store.OlAnalyticsBucketPeriodType
	BucketStart time.Time
	Limit       int32
	Offset      int32
}

type BookEntry struct {
	BookID int64
	Value  MetricValue
}

type MetricService interface {
	Get(ctx context.Context, metric MetricType, bookIds []int64) (map[int64]MetricValues, error)
	GetMetricValue(ctx context.Context, metric MetricType, bookId int64, period store.OlAnalyticsBucketPeriodType) (MetricValue, error)
	GetTopBooks(ctx context.Context, query GetTopBooksByMetricQuery) ([]BookEntry, error)
}

type dummyMetricService struct{}

// Get implements [MetricService].
func (d *dummyMetricService) Get(ctx context.Context, metric MetricType, bookIds []int64) (map[int64]MetricValues, error) {
	return make(map[int64]MetricValues), nil
}

// GetMetricValue implements [MetricService].
func (d *dummyMetricService) GetMetricValue(ctx context.Context, metric MetricType, bookId int64, period store.OlAnalyticsBucketPeriodType) (MetricValue, error) {
	return MetricValue{}, nil
}

func (d *dummyMetricService) GetTopBooks(ctx context.Context, query GetTopBooksByMetricQuery) ([]BookEntry, error) {
	return make([]BookEntry, 0), nil
}

func NewDummyMetricService() MetricService {
	return &dummyMetricService{}
}
