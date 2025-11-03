package app

import (
	"context"
	"log/slog"
	"net"
	"runtime/debug"
	"sync"
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
	IncrBookView(ctx context.Context, bookID int64, userID uuid.NullUUID, ip net.IP) error
	GetBookViews(ctx context.Context, bookID int64) (Views, error)

	ApplyPendingViews(ctx context.Context)
}

type AnalyticsBackgroundService struct {
	started        bool
	mx             sync.Mutex
	analytics      AnalyticsViewsService
	stopWg         sync.WaitGroup
	stopRequested  bool
	nextLaunchTime time.Time
	parentCtx      context.Context
}

func NewAnalyticsBackgroundService(analytics AnalyticsViewsService) *AnalyticsBackgroundService {
	srv := &AnalyticsBackgroundService{analytics: analytics}
	return srv
}

func (s *AnalyticsBackgroundService) Start() {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.started {
		return
	}

	s.started = true
	go s.start()
}

func (s *AnalyticsBackgroundService) start() {
	s.stopWg.Add(1)
	defer s.stopWg.Done()

	var cancel context.CancelFunc
	s.parentCtx, cancel = context.WithCancel(context.Background())
	defer cancel()

	for !s.stopRequested {
		time.Sleep(time.Minute)

		if s.itIsTime() {
			s.nextLaunchTime = s.calculateNextLaunchTime()
			s.process()
		}
	}
}

func (s *AnalyticsBackgroundService) process() {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("AnalyticsBackgroundService panicked")
			debug.PrintStack()
		}
	}()

	s.analytics.ApplyPendingViews(s.parentCtx)
}

func (s *AnalyticsBackgroundService) itIsTime() bool {
	if (s.nextLaunchTime == time.Time{}) {
		return true
	}

	now := time.Now()
	return s.nextLaunchTime.Before(now)
}

func (s *AnalyticsBackgroundService) calculateNextLaunchTime() time.Time {
	now := time.Now().UTC()
	nextHourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC).Add(time.Hour)
	return nextHourStart
}
