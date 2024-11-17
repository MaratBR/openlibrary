package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type FavoriteRecalculationBackgroundService struct {
	db                     DB
	scheduledRecalculation map[int64]struct{}
	mutex                  sync.Mutex
	ctx                    context.Context
	cancel                 func()
}

func NewFavoriteRecalculationBackgroundService(db DB) *FavoriteRecalculationBackgroundService {
	return &FavoriteRecalculationBackgroundService{db: db, scheduledRecalculation: map[int64]struct{}{}}
}

func (s *FavoriteRecalculationBackgroundService) ScheduleRecalculation(bookID int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.scheduledRecalculation[bookID] = struct{}{}
}

func (s *FavoriteRecalculationBackgroundService) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.run()
}

func (s *FavoriteRecalculationBackgroundService) Stop() {
	if s.cancel == nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cancel()
}

func (s *FavoriteRecalculationBackgroundService) run() {
	for s.ctx.Err() == nil {
		time.Sleep(time.Second)
		if len(s.scheduledRecalculation) > 0 {
			slog.Info("recalculating favorites for books", "count", len(s.scheduledRecalculation))

			s.mutex.Lock()
			books := s.scheduledRecalculation
			s.scheduledRecalculation = map[int64]struct{}{}
			s.mutex.Unlock()

			queries := store.New(s.db)

			for bookID, _ := range books {
				err := queries.RecalculateBookFavorites(s.ctx, bookID)
				if err != nil {
					slog.Error("failed to recalculate favorites for book", "err", bookID, "err", err.Error())
					if err == context.Canceled {
						break
					}
				}
			}

		}
	}

	if s.ctx.Err() == context.Canceled {
		slog.Info("FavoriteRecalculationBackgroundService: background service stopped")
	} else {
		slog.Error("FavoriteRecalculationBackgroundService: background service failed", "err", s.ctx.Err().Error())
	}
}
