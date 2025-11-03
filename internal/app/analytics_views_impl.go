package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type analyticsViewsService struct {
	counter AnalyticsCounters
	db      DB
}

// GetBookViews implements AnalyticsViewsService.
func (a *analyticsViewsService) GetBookViews(ctx context.Context, bookID int64) (Views, error) {
	periods := CurrentAnalyticsPeriods(time.Now())
	queries := store.New(a.db)
	var views Views
	rows, _ := queries.Analytics_GetViews(ctx, store.Analytics_GetViewsParams{
		BookID:      bookID,
		YearPeriod:  int32(periods.Year),
		MonthPeriod: int32(periods.Month),
		WeekPeriod:  int32(periods.Week),
		DayPeriod:   int32(periods.Day),
		HourPeriod:  int32(periods.Hour),
	})

	for _, row := range rows {
		if row.Period == int32(periods.Year) {
			views.Year = row.Count
		} else if row.Period == int32(periods.Month) {
			views.Month = row.Count
		} else if row.Period == int32(periods.Week) {
			views.Week = row.Count
		} else if row.Period == int32(periods.Day) {
			views.Day = row.Count
		} else if row.Period == int32(periods.Hour) {
			views.Hour = row.Count
		} else if row.Period == int32(AnalyticsPeriodTotal) {
			views.Total = row.Count
		}
	}

	hourCount, err := a.counter.Get(ctx, fmt.Sprintf("%d", bookID))
	if err != nil {
		slog.Error("failed to get views count from redis", "err", err, "bookID", bookID)
	} else if hourCount > 0 {
		views.Hour += hourCount
		views.Day += hourCount
		views.Week += hourCount
		views.Month += hourCount
		views.Year += hourCount
		views.Total += hourCount
	}

	return views, nil
}

// IncrBookView implements AnalyticsViewsService.
func (a *analyticsViewsService) IncrBookView(ctx context.Context, bookID int64, userID uuid.NullUUID, ip net.IP) error {
	var uniqueId string
	if userID.Valid {
		uniqueId = userID.UUID.String()
	} else {
		uniqueId = ip.String()
	}
	err := a.counter.Incr(ctx, fmt.Sprintf("%d", bookID), uniqueId, 1, time.Hour)
	return err
}

func (a *analyticsViewsService) ApplyPendingViews(ctx context.Context) {
	queries := store.New(a.db)

	periods := CurrentAnalyticsPeriods(time.Now())

	views, err := a.counter.GetPendingCounters(ctx)
	if err != nil {
		slog.Error("could not find pending counters", "err", err)
		return
	}

	slog.Warn("found pending counters", "count", len(views))
	for key, count := range views {
		if count <= 0 {
			continue
		}
		id, err := strconv.ParseInt(key, 10, 62)
		if err != nil {
			continue
		}

		updateCounter(queries, ctx, periods.Hour, id, count)
		updateCounter(queries, ctx, periods.Day, id, count)
		updateCounter(queries, ctx, periods.Week, id, count)
		updateCounter(queries, ctx, periods.Month, id, count)
		updateCounter(queries, ctx, periods.Year, id, count)
		updateCounter(queries, ctx, AnalyticsPeriodTotal, id, count)
	}

}

func updateCounter(queries *store.Queries, ctx context.Context, period AnalyticsPeriod, bookID int64, incr int64) {
	slog.Debug("incrementing views counter", "period", period, "bookID", bookID, "incrBy", incr)
	err := queries.Analytics_IncrView(ctx, store.Analytics_IncrViewParams{
		BookID: bookID,
		Count:  incr,
		Period: int32(period),
	})
	if err != nil {
		slog.Error("failed to update counter for period", "period", period, "err", err)
	}
}

func NewAnalyticsViewsService(db DB, counters AnalyticsCounters) AnalyticsViewsService {
	return &analyticsViewsService{
		counter: counters,
		db:      db,
	}
}
