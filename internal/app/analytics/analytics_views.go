package analytics

import (
	"context"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

type ViewsService interface {
	IncrBookView(ctx context.Context, bookID int64, meta ViewMetadata) error
	IncrChapterView(ctx context.Context, bookID, chapterID int64, meta ViewMetadata) error
	GetBookViews(ctx context.Context, bookID int64) (Views, error)
	GetMostViewedBooks(ctx context.Context, period AnalyticsPeriod) ([]BookViewEntry, error)
	CommitPendingViewsToDB(ctx context.Context)
}

type AnalyticsBackgroundService struct {
	started        bool
	mx             sync.Mutex
	analytics      ViewsService
	stopWg         sync.WaitGroup
	stopRequested  bool
	stopCh         chan struct{}
	nextLaunchTime time.Time
	parentCtx      context.Context
	log            *zap.SugaredLogger
}

func NewAnalyticsBackgroundService(analytics ViewsService, log *zap.SugaredLogger, lc fx.Lifecycle) *AnalyticsBackgroundService {
	srv := &AnalyticsBackgroundService{analytics: analytics, stopCh: make(chan struct{}), log: log}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			srv.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Stop()
			return nil
		},
	})
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
	s.log.Info("AnalyticsBackgroundService started")
}

func (s *AnalyticsBackgroundService) Stop() {
	s.mx.Lock()
	defer s.mx.Unlock()

	if !s.started {
		return
	}

	s.stopRequested = true
	s.stopCh <- struct{}{}
	s.stopWg.Wait()
}

func (s *AnalyticsBackgroundService) start() {
	s.stopWg.Add(1)
	defer s.stopWg.Done()

	var cancel context.CancelFunc
	s.parentCtx, cancel = context.WithCancel(context.Background())
	defer cancel()

forLoop:
	for !s.stopRequested {
		if s.itIsTime() {
			s.nextLaunchTime = s.calculateNextLaunchTime()
			s.process()
		}

		select {
		case <-time.After(time.Minute):
		case <-s.stopCh:
			break forLoop
		}
	}

	s.log.Info("AnalyticsBackgroundService stopped")
}

func (s *AnalyticsBackgroundService) process() {
	defer func() {
		if rec := recover(); rec != nil {
			s.log.Error("AnalyticsBackgroundService panicked")
			debug.PrintStack()
		}
	}()

	s.log.Debug("AnalyticsBackgroundService.process")
	s.analytics.CommitPendingViewsToDB(s.parentCtx)
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
