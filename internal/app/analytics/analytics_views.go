package analytics

import (
	"context"
	"net"
	"time"

	"github.com/gofrs/uuid"
)

type AnalyticsPeriod int32

const ANALYTICS_PERIOD_TOTAL AnalyticsPeriod = 0

type AnalyticsPeriods struct {
	Hour  AnalyticsPeriod
	Day   AnalyticsPeriod
	Week  AnalyticsPeriod
	Month AnalyticsPeriod
	Year  AnalyticsPeriod
}

// go2tsdef:generate
type MetricValue struct {
	Samples  int64   `json:"samples"`
	ValueSum float64 `json:"valueSum"`
}

// go2tsdef:generate
type MetricValues struct {
	Total MetricValue `json:"total"`
	Year  MetricValue `json:"year"`
	Month MetricValue `json:"month"`
	Week  MetricValue `json:"week"`
	Day   MetricValue `json:"day"`
}

func CurrentAnalyticsPeriods(now time.Time) AnalyticsPeriods {
	var (
		year      int
		month     int
		dayInYear int
		week      int
		weekYear  int
	)

	year = now.Year()
	month = int(now.Month())
	dayInYear = now.YearDay()
	weekYear, week = now.ISOWeek()

	return AnalyticsPeriods{
		Year:  AnalyticsPeriod(50_000 + year),
		Month: AnalyticsPeriod(4_000_000 + year*100 + month),
		Day:   AnalyticsPeriod(30_000_000 + year*1_000 + dayInYear),
		Week:  AnalyticsPeriod(2_000_000 + weekYear*100 + week),
		Hour:  AnalyticsPeriod(1_000_000_000 + year*100_000 + dayInYear*100 + now.Hour()),
	}
}

type BookViewEntry struct {
	BookID int64
	Views  int64
}

type ViewMetadata struct {
	UserID uuid.NullUUID
	IP     net.IP
}

func (m ViewMetadata) UniqueID() string {
	if m.UserID.Valid {
		return "U" + m.UserID.UUID.String()
	}

	if m.IP != nil {
		return m.IP.String()
	}

	return "unknown"
}

// Deprecated: use EventSink instead
type ViewsService interface {
	GetMostViewedBooks(ctx context.Context, period AnalyticsPeriod) ([]BookViewEntry, error)
}
