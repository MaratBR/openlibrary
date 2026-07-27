package analytics

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type eventBackgroundService struct {
	log       *zap.SugaredLogger
	repo      EventRepository
	wg        sync.WaitGroup
	cancelCtx context.CancelFunc
	running   bool
	queue     []Event
	queueMx   sync.Mutex
}

func (svc *eventBackgroundService) SubmitEvent(ctx context.Context, event Event) {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()
	svc.queue = append(svc.queue, event)
}

func (svc *eventBackgroundService) start(ctx context.Context) error {
	if svc.running {
		return errors.New("eventBackgroundService is already running")
	}
	svc.running = true
	var stopingCtx context.Context
	stopingCtx, svc.cancelCtx = context.WithCancel(context.Background())
	go svc.run(stopingCtx)
	return nil
}

func (svc *eventBackgroundService) run(stopingCtx context.Context) {
	if !svc.running {
		svc.log.Warn("eventBackgroundService is stopping immidiately - running is false, is this a bug?")
		return
	}

	svc.wg.Add(1)

loop:
	for {

		select {
		case <-stopingCtx.Done():
			break loop
		case <-time.After(time.Second * 5):
			err := svc.flushEvents(stopingCtx)
			if err != nil {
				svc.log.Errorw("failed to flush events in eventBackgroundService", "err", err)
			} else if len(svc.queue) > 0 {
				svc.log.Debugw("flushed events", "count", len(svc.queue))
			}
			svc.queue = svc.queue[:0]
		}

	}

	svc.wg.Done()
}

func (svc *eventBackgroundService) flushEvents(ctx context.Context) error {
	if len(svc.queue) == 0 {
		return nil
	}

	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	err := svc.repo.Insert(ctx, svc.queue)
	if err != nil {
		return err
	}

	return nil
}

func (svc *eventBackgroundService) stop(ctx context.Context) error {
	if !svc.running {
		return nil
	}
	svc.cancelCtx()
	svc.wg.Wait() // TODO timeout?
	return nil
}

func newEventBackgroundService(
	log *zap.SugaredLogger,
	repo EventRepository,
	lc fx.Lifecycle,
) *eventBackgroundService {
	svc := &eventBackgroundService{
		log:  log,
		repo: repo,
	}

	lc.Append(fx.Hook{
		OnStart: svc.start,
		OnStop:  svc.stop,
	})

	return svc
}
