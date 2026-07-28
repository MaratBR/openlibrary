package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type EventProcessorWorkerHandle interface {
	WakeUp()
}

type eventProcessorWorker struct {
	log                       *zap.SugaredLogger
	eventProcessor            EventProcessor
	eventRepo                 EventRepository
	eventProcessorWorkerState eventProcessorWorkerState
	currentCursor             int64

	running bool
	done    chan struct{}
	wakeUp  chan struct{}
	cancel  context.CancelFunc
}

func (p *eventProcessorWorker) WakeUp() {
	if p.wakeUp != nil && len(p.wakeUp) == 0 {
		p.wakeUp <- struct{}{}
	}
}

func (p *eventProcessorWorker) start(ctx context.Context) error {
	if p.running {
		return errors.New("already running")
	}

	var stoppingCtx context.Context
	stoppingCtx, p.cancel = context.WithCancel(context.Background())
	p.running = true
	p.done = make(chan struct{})
	p.wakeUp = make(chan struct{})
	go p.run(stoppingCtx)
	return nil
}

func (p *eventProcessorWorker) run(ctx context.Context) {
	p.init(ctx)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-p.wakeUp:
			p.log.Debug("eventProcessor: received wake-up signal")
			p.process(ctx)
		case <-time.After(time.Second * 120):
			p.process(ctx)
		}
	}

	close(p.done)
}

func (p *eventProcessorWorker) init(ctx context.Context) {
	err := retry.New(
		retry.Attempts(6),
		retry.DelayType(retry.FullJitterBackoffDelay),
		retry.MaxDelay(time.Second*10),
		retry.MaxJitter(time.Second),
	).Do(func() error {
		var err error
		p.currentCursor, err = p.eventProcessorWorkerState.GetCurrentCursor(ctx)
		return err
	})
	if err != nil {
		p.log.Errorw("eventProcessorWorker initialization failed", "err", err)
		panic(err)
	}
}

func (p *eventProcessorWorker) process(ctx context.Context) {
	newCursor, events, err := p.eventRepo.GetEvents(ctx, 1000, p.currentCursor)
	if err != nil {
		p.log.Error("failed to fetch events", "err", err)
		return
	}

	err = p.eventProcessor.Process(ctx, events)

	if err != nil {
		p.log.Error("eventProcessorWorker: failed to process events with event processor")
	} else {
		p.log.Debugw("eventProcessorWorker: finished processing", "count", len(events), "newCursor", newCursor)
	}
	p.currentCursor = newCursor
}

func (p *eventProcessorWorker) stop(ctx context.Context) error {
	if !p.running {
		return nil
	}

	p.cancel()
	<-p.done

	return nil
}

func newEventProcessorWorker(log *zap.SugaredLogger, eventProcessor EventProcessor, eventRepo EventRepository, eventProcessorWorkerState eventProcessorWorkerState, lc fx.Lifecycle) *eventProcessorWorker {
	svc := &eventProcessorWorker{
		log:                       log,
		eventProcessor:            eventProcessor,
		eventRepo:                 eventRepo,
		eventProcessorWorkerState: eventProcessorWorkerState,
	}

	lc.Append(fx.Hook{
		OnStart: svc.start,
		OnStop:  svc.stop,
	})

	return svc
}

func newEventProcessorWorkerHandle(worker *eventProcessorWorker) EventProcessorWorkerHandle {
	return worker
}
