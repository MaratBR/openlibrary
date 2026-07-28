package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

type metricService struct {
	queries *store.Queries
}

func (m *metricService) Get(ctx context.Context, metric MetricType, bookIds []int64) (map[int64]MetricValues, error) {
	rows, err := m.queries.Analytics_GetMetricValues(ctx, store.Analytics_GetMetricValuesParams{
		BookIds: bookIds,
		Metric:  string(metric),
	})
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}

	mapping := make(map[int64]MetricValues, len(bookIds))
	for _, bookID := range bookIds {
		mapping[bookID] = MetricValues{}
	}

	for _, row := range rows {
		var (
			values MetricValues
			ok     bool
		)
		if values, ok = mapping[row.BookID]; !ok {
			continue
		}

		switch row.BucketType {
		case store.OlAnalyticsBucketPeriodTypeAll:
			values.Total.Samples = row.SamplesCount
			values.Total.ValueSum = row.ValueSum
		case store.OlAnalyticsBucketPeriodTypeYear:
			values.Year.Samples = row.SamplesCount
			values.Year.ValueSum = row.ValueSum
		case store.OlAnalyticsBucketPeriodTypeMonth:
			values.Month.Samples = row.SamplesCount
			values.Month.ValueSum = row.ValueSum
		case store.OlAnalyticsBucketPeriodTypeWeek:
			values.Week.Samples = row.SamplesCount
			values.Week.ValueSum = row.ValueSum
		case store.OlAnalyticsBucketPeriodTypeDay:
			values.Day.Samples = row.SamplesCount
			values.Day.ValueSum = row.ValueSum
		}

		mapping[row.BookID] = values
	}

	return mapping, nil
}

func (m *metricService) GetMetricValue(ctx context.Context, metric MetricType, bookId int64, period store.OlAnalyticsBucketPeriodType) (MetricValue, error) {
	val, err := m.queries.Analytics_GetMetricValue(ctx, store.Analytics_GetMetricValueParams{
		Metric:     string(metric),
		BookID:     bookId,
		BucketType: period,
	})
	if err != nil {
		if err == store.ErrNoRows {
			return MetricValue{}, nil
		}
		return MetricValue{}, err
	}

	return MetricValue{
		ValueSum: val.ValueSum,
		Samples:  val.SamplesCount,
	}, nil
}

func (m *metricService) GetTopBooks(ctx context.Context, query GetTopBooksByMetricQuery) ([]BookEntry, error) {
	rows, err := m.queries.Analytics_GetTopBooksBySamplesCount(ctx, store.Analytics_GetTopBooksBySamplesCountParams{
		BucketType:  query.Period,
		BucketStart: pgtype.Timestamptz{Valid: true, Time: query.BucketStart},
		Metric:      string(query.Metric),
		Skip:        query.Offset,
		Limit:       query.Limit,
	})

	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}

	return commonutil.MapSlice(rows, func(row store.Analytics_GetTopBooksBySamplesCountRow) BookEntry {
		return BookEntry{
			BookID: row.BookID,
			Value:  NewMetricValue(row.SamplesCount, row.ValueSum),
		}
	}), nil
}

func newMetricService(db store.DBTX) MetricService {
	return &metricService{queries: store.New(db)}
}
