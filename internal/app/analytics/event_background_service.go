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
	log                        *zap.SugaredLogger
	repo                       EventRepository
	done                       chan struct{}
	cancelCtx                  context.CancelFunc
	running                    bool
	queue                      []Event
	queueMx                    sync.Mutex
	eventProcessorWorkerHandle EventProcessorWorkerHandle
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
	svc.done = make(chan struct{})
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

loop:
	for {

		select {
		case <-stopingCtx.Done():
			break loop
		case <-time.After(time.Second * 5):
			svc.log.Debug("eventBackgroundService: woke up")
			svc.process(stopingCtx)
		}
	}

	close(svc.done)
	svc.log.Debug("eventBackground service stopped")
}

func (svc *eventBackgroundService) process(ctx context.Context) {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()
	err := svc.flush(ctx)
	if err != nil {
		svc.log.Errorw("eventBackgroundService: failed to flush events", "err", err)
	} else if len(svc.queue) > 0 {
		svc.log.Debugw("eventBackgroundService: flushed events", "count", len(svc.queue))
		svc.eventProcessorWorkerHandle.WakeUp()
		svc.clear()
	}
}

func (svc *eventBackgroundService) flush(ctx context.Context) error {
	if len(svc.queue) == 0 {
		return nil
	}

	err := svc.repo.Insert(ctx, svc.queue)
	if err != nil {
		return err
	}

	return nil
}

func (svc *eventBackgroundService) clear() {
	svc.queue = svc.queue[:0]

}

func (svc *eventBackgroundService) stop(ctx context.Context) error {
	if !svc.running {
		return nil
	}
	svc.cancelCtx()
	<-svc.done
	return nil
}

func newEventBackgroundService(
	log *zap.SugaredLogger,
	repo EventRepository,
	eventProcessorWorkerHandle EventProcessorWorkerHandle,
	lc fx.Lifecycle,
) *eventBackgroundService {
	svc := &eventBackgroundService{
		log:                        log,
		repo:                       repo,
		eventProcessorWorkerHandle: eventProcessorWorkerHandle,
	}

	lc.Append(fx.Hook{
		OnStart: svc.start,
		OnStop:  svc.stop,
	})

	return svc
}
