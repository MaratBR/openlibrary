package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

type metricRepository struct {
	db  store.DBTX
	log *zap.SugaredLogger
}

// AddMetrics implements [MetricRepository].
func (m *metricRepository) AddMetrics(ctx context.Context, metricRecords []MetricRecord) error {
	queries := store.New(m.db)

	metrics := make([]string, len(metricRecords))
	bookIds := make([]int64, len(metricRecords))
	samples := make([]int64, len(metricRecords))
	values := make([]float64, len(metricRecords))

	for i, metric := range metricRecords {
		metrics[i] = string(metric.Type)
		bookIds[i] = metric.BookID
		values[i] = metric.Value
		samples[i] = metric.Samples
	}

	params := store.Analytics_UpdateMetricsParams{
		Metrics: metrics,
		BookIds: bookIds,
		Values:  values,
		Samples: samples,
	}
	m.log.Debugw("metricRepository: Analytics_UpdateMetrics", "params", params)
	err := queries.Analytics_UpdateMetrics(ctx, params)

	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func newMetricRepository(db store.DBTX, log *zap.SugaredLogger) MetricRepository {
	return &metricRepository{db: db, log: log}
}
