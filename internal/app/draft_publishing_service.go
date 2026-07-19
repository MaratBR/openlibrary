package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DraftPublishingService struct {
	cancel context.CancelFunc
	done   chan struct{}
}

func NewDraftPublishingService(
	db DB,
	reindex BookReindexService,
	lifecycle fx.Lifecycle,
	log *zap.SugaredLogger,
) *DraftPublishingService {
	service := &DraftPublishingService{done: make(chan struct{})}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
			service.cancel = cancel
			go service.run(ctx, store.New(db), reindex, log)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if service.cancel != nil {
				service.cancel()
			}
			select {
			case <-service.done:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
	return service
}

func (s *DraftPublishingService) run(ctx context.Context, queries *store.Queries, reindex BookReindexService, log *zap.SugaredLogger) {
	defer close(s.done)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		s.publishDue(ctx, queries, reindex, log)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *DraftPublishingService) publishDue(ctx context.Context, queries *store.Queries, reindex BookReindexService, log *zap.SugaredLogger) {
	bookIDs, err := queries.Draft_PublishDue(ctx)
	if err != nil {
		if ctx.Err() == nil {
			log.Errorw("failed to publish scheduled drafts", "error", err)
		}
		return
	}

	for _, bookID := range bookIDs {
		if err := queries.RecalculateBookStats(ctx, bookID); err != nil {
			log.Errorw("failed to recalculate a scheduled book", "bookID", bookID, "error", err)
			continue
		}
		reindex.ScheduleReindex(ctx, bookID)
	}
}
