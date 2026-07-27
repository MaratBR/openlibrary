package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
)

type metricRepository struct {
	db store.DBTX
}

// AddMetrics implements [MetricRepository].
func (m *metricRepository) AddMetrics(ctx context.Context, metric []Metric) error {
	queries := store.New(m.db)

	metrics := make([]string, len(metric))
	bookIds := make([]int64, len(metric))
	values := make([]float64, len(metric))

	err := queries.Analytics_UpdateMetrics(ctx, store.Analytics_UpdateMetricsParams{
		Metrics: metrics,
		BookIds: bookIds,
		Values:  values,
	})

	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func newMetricRepository(db store.DBTX) MetricRepository {
	return &metricRepository{db: db}
}
