package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type BookBackgroundService interface {
	Start() error
	Stop()
	ScheduleBookRecalculation(bookID int64)
}

type dummyBookBackgroundService struct{}

func (d *dummyBookBackgroundService) ScheduleBookRecalculation(bookID int64) {}
func (d *dummyBookBackgroundService) Start() error                           { return nil }
func (d *dummyBookBackgroundService) Stop()                                  {}

func NewDummyBookBackgroundService() BookBackgroundService {
	return &dummyBookBackgroundService{}
}

type bookBackgroundService struct {
	db      DB
	mutex   sync.Mutex
	wg      sync.WaitGroup
	queue   chan int64
	running bool
}

func NewBookBackgroundService(db DB) BookBackgroundService {
	return &bookBackgroundService{db: db}
}

func (s *bookBackgroundService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return nil
	}

	s.queue = make(chan int64, 10000)
	s.running = true
	s.wg.Add(1)
	go s.worker(context.Background())
	return nil
}

func (s *bookBackgroundService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return
	}

	close(s.queue)
	s.running = false
}

func (s *bookBackgroundService) ScheduleBookRecalculation(bookID int64) {
	if !s.running {
		panic("cannot schedule book recalculation: background service is not running!")
	}
	s.queue <- bookID
}

func (c *bookBackgroundService) worker(baseCtx context.Context) {
	defer c.wg.Done()
	queries := store.New(c.db)
	interval := time.Second

	books := make(map[int64]struct{}, 256)
	batchDeadline := time.Now().Add(interval)

	for {
		bookID, ok := <-c.queue
		if !ok {
			break
		}
		books[bookID] = struct{}{}

		if len(books) == 256 || time.Now().After(batchDeadline) {
			batchDeadline = time.Now().Add(interval)
			if len(books) == 0 {
				continue
			}

			for bookID, _ := range books {
				ctx, cancel := context.WithTimeout(baseCtx, time.Second*5)
				defer cancel()

				err := queries.RecalculateBookRating(ctx, bookID)
				if err != nil {
					slog.Error("failed to recalculate favorites for book", "err", err.Error())
				}
			}
			clear(books)
		}
		time.Sleep(time.Millisecond)
	}
}
