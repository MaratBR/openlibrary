package app

import (
	"context"
	"net"
	"time"

	"github.com/gofrs/uuid"
)

type AnalyticsPeriod int32

const AnalyticsPeriodTotal AnalyticsPeriod = 0

type AnalyticsPeriods struct {
	Hour  AnalyticsPeriod
	Day   AnalyticsPeriod
	Week  AnalyticsPeriod
	Month AnalyticsPeriod
	Year  AnalyticsPeriod
}

type Views struct {
	Total int64
	Year  int64
	Month int64
	Week  int64
	Day   int64
	Hour  int64
}

func CurrentAnalyticsPeriods(now time.Time) AnalyticsPeriods {
	var (
		year      int
		month     int
		dayInYear int
		week      int
	)

	year = now.Year()
	month = int(now.Month())
	dayInYear = now.Day()
	_, week = now.ISOWeek()

	return AnalyticsPeriods{
		Year:  AnalyticsPeriod(50_000 + year),
		Month: AnalyticsPeriod(4_000_000 + year*100 + month),
		Day:   AnalyticsPeriod(30_000_000 + year*1_000 + dayInYear),
		Week:  AnalyticsPeriod(2_000_000 + year*100 + week),
		Hour:  AnalyticsPeriod(1_000_000_000 + year*100_000 + dayInYear*100 + now.Hour()),
	}
}

type AnalyticsViewsService interface {
	IncrBookView(ctx context.Context, bookID int64, userID Nullable[uuid.UUID], ip net.IP) error
	GetBookViews(ctx context.Context, bookID int64) (Views, error)
}
