package analytics

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type metricRepository struct {
	db  dal.DB
	log *zap.SugaredLogger
}

// AddMetrics implements [MetricRepository].
func (m *metricRepository) AddMetrics(ctx context.Context, tx store.DBTX, metricRecords []MetricRecord) error {

	groupedByDay := commonutil.GroupBy(metricRecords, func(m MetricRecord) time.Time {
		return newBucketStartTime(m.OccurredAt).Day()
	})

	queries := store.New(tx)

	for dayStart, records := range groupedByDay {
		metrics := make([]string, len(records))
		bookIds := make([]int64, len(records))
		samples := make([]int64, len(records))
		values := make([]float64, len(records))

		for i, metric := range records {
			metrics[i] = string(metric.Type)
			bookIds[i] = metric.BookID
			values[i] = metric.Value
			samples[i] = metric.Samples
		}

		params := store.Analytics_UpdateMetricsParams{
			Metrics:  metrics,
			BookIds:  bookIds,
			Values:   values,
			Samples:  samples,
			DayStart: pgtype.Timestamptz{Valid: true, Time: dayStart},
		}
		m.log.Debugw("metricRepository: Analytics_UpdateMetrics", "params", params)
		err := queries.Analytics_UpdateMetrics(ctx, params)

		if err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}

	}

	return nil
}

func newMetricRepository(db dal.DB, log *zap.SugaredLogger) MetricRepository {
	return &metricRepository{db: db, log: log}
}
