package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/zap"
)

type analyticsViewsService struct {
	counterBooks    CountersNamespace
	counterChapters CountersNamespace
	db              store.DBTX
	log             *zap.SugaredLogger
	commitMX        sync.Mutex
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
			views.Year = row.ViewCount
		} else if row.Period == int32(periods.Month) {
			views.Month = row.ViewCount
		} else if row.Period == int32(periods.Week) {
			views.Week = row.ViewCount
		} else if row.Period == int32(periods.Day) {
			views.Day = row.ViewCount
		} else if row.Period == int32(periods.Hour) {
			views.Hour = row.ViewCount
		} else if row.Period == int32(ANALYTICS_PERIOD_TOTAL) {
			views.Total = row.ViewCount
		}
	}

	pendingViews, err := a.counterBooks.Get(ctx, fmt.Sprintf("%d", bookID))
	if err != nil {
		slog.Error("failed to get views count from redis", "err", err, "bookID", bookID)
	} else if pendingViews > 0 {
		views.Hour += pendingViews
		views.Day += pendingViews
		views.Week += pendingViews
		views.Month += pendingViews
		views.Year += pendingViews
		views.Total += pendingViews
	}

	return views, nil
}

// IncrBookView implements AnalyticsViewsService.
func (a *analyticsViewsService) IncrBookView(ctx context.Context, bookID int64, meta ViewMetadata) error {
	err := a.counterBooks.Incr(ctx, fmt.Sprintf("%d", bookID), meta.UniqueID(), 1, time.Hour)
	return err
}

// IncrBookView implements AnalyticsViewsService.
func (a *analyticsViewsService) IncrChapterView(ctx context.Context, bookID, chapterID int64, meta ViewMetadata) error {
	err := a.counterBooks.Incr(ctx, fmt.Sprintf("%d,%d", bookID, chapterID), meta.UniqueID(), 1, time.Hour)
	return err
}

func (a *analyticsViewsService) CommitPendingViewsToDB(ctx context.Context) {
	a.commitMX.Lock()
	defer a.commitMX.Unlock()

	queries := store.New(a.db)
	periods := CurrentAnalyticsPeriods(time.Now())

	views, err := a.counterBooks.PullPendingCounters(ctx, true)
	if err != nil {
		slog.Error("could not find pending counters", "err", err)
		return
	}

	if len(views) == 0 {
		return
	}

	slog.Warn("found pending counters", "count", len(views))

	for key, count := range views {
		if count <= 0 {
			continue
		}
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			a.log.Errorw("failed to parse book id as int64", "v", key)
			continue
		}

		updateCounter(ctx, queries, periods.Hour, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
		updateCounter(ctx, queries, periods.Day, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
		updateCounter(ctx, queries, periods.Week, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
		updateCounter(ctx, queries, periods.Month, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
		updateCounter(ctx, queries, periods.Year, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
		updateCounter(ctx, queries, ANALYTICS_PERIOD_TOTAL, id, store.ANALYTICS_VIEW_COUNTER_TYPE_BOOK, id, count, a.log)
	}
}

func (a *analyticsViewsService) GetMostViewedBooks(ctx context.Context, period AnalyticsPeriod) ([]BookViewEntry, error) {
	queries := store.New(a.db)
	// TODO take into account chapter views too?
	rows, err := queries.Analytics_GetMostViewedBooksByBookViewsOnly(ctx, store.Analytics_GetMostViewedBooksByBookViewsOnlyParams{
		Period: int32(period),
		Limit:  50,
	})
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}

	entries := make([]BookViewEntry, len(rows))

	for i := range rows {
		entries[i] = BookViewEntry{
			BookID: rows[i].BookID,
			Views:  rows[i].ViewCount,
		}
	}

	return entries, err

}

func updateCounter(
	ctx context.Context,
	queries *store.Queries,
	period AnalyticsPeriod,
	bookID int64,
	entityType int16,
	entityID int64,
	incr int64,
	log *zap.SugaredLogger,
) {
	log.Debug("incrementing views counter", "period", period, "bookID", bookID, "incrBy", incr)
	err := queries.Analytics_IncrView(ctx, store.Analytics_IncrViewParams{
		BookID:     bookID,
		EntityType: entityType,
		EntityID:   entityID,
		IncrBy:     incr,
		Period:     int32(period),
	})
	if err != nil {
		log.Error("failed to update counter for period", "period", period, "err", err)
	}
}

func NewAnalyticsViewsService(db store.DBTX, counters Counters, log *zap.SugaredLogger) ViewsService {
	return &analyticsViewsService{
		counterBooks:    counters.Namespace("views"),
		counterChapters: counters.Namespace("views_chapters"),
		db:              db,
		log:             log,
	}
}
