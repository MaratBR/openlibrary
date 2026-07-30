package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/avast/retry-go/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

	db      dal.DB
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
	err := p.retry(func() error {
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
	ctx, span := tracer.Start(ctx, "eventProcessorWorker.process")
	defer span.End()

	newCursor, events, err := p.eventRepo.GetEvents(ctx, 1000, p.currentCursor)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.log.Error("failed to fetch events", "err", err)
		return
	}
	span.SetAttributes(attribute.Int("analytics.event_count", len(events)))

	if len(events) == 0 {
		return
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.log.Errorw("failed to open transaction in eventProcessorWorker.process", "err", err)
		return
	}

	err = p.eventProcessor.Process(ctx, tx, events)

	if err != nil {
		dal.RollbackTx(ctx, tx)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.log.Error("eventProcessorWorker: failed to process events with event processor")
		return
	} else {
		p.log.Debugw("eventProcessorWorker: finished processing", "count", len(events), "newCursor", newCursor)
	}

	err = p.retry(func() error {
		return p.eventProcessorWorkerState.SaveCurrentCursor(ctx, tx, newCursor)
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.log.Errorw("failed to update worker's state")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.log.Errorw("failed to commit transaction")
		return
	}

	p.currentCursor = newCursor
}

func (p *eventProcessorWorker) retry(fn func() error) error {
	return retry.New(
		retry.Attempts(5),
		retry.MaxJitter(time.Millisecond*200),
		retry.MaxDelay(time.Second*5),
		retry.DelayType(retry.FullJitterBackoffDelay),
	).Do(fn)
}

func (p *eventProcessorWorker) stop(ctx context.Context) error {
	if !p.running {
		return nil
	}

	p.cancel()
	<-p.done

	return nil
}

func newEventProcessorWorker(log *zap.SugaredLogger, eventProcessor EventProcessor, eventRepo EventRepository, db dal.DB, eventProcessorWorkerState eventProcessorWorkerState, lc fx.Lifecycle) *eventProcessorWorker {
	svc := &eventProcessorWorker{
		log:                       log,
		eventProcessor:            eventProcessor,
		eventRepo:                 eventRepo,
		eventProcessorWorkerState: eventProcessorWorkerState,
		db:                        db,
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
