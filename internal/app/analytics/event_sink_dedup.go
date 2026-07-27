package analytics

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type DeduplicationEventSink struct {
	inner       EventSink
	dedupAtomic Atomic
	log         *zap.SugaredLogger
}

// SubmitEvent implements [EventSink].
func (d *DeduplicationEventSink) SubmitEvent(ctx context.Context, event Event) {
	err := d.submitEvent(ctx, event)
	if err != nil {
		d.log.Errorw("failed to submit event to deduplicating sink", "err", err)
	}

}

func (d *DeduplicationEventSink) submitEvent(ctx context.Context, event Event) error {

	acquired, err := d.dedupAtomic.TryAcquire(ctx, event.Key(), time.Hour*24)
	if err != nil {
		// TODO if atomic fails to we just submit it downstream or nah?
		return fmt.Errorf("acquire event deduplication key: %w", err)
	}

	d.log.Debugw("submitEvent", "acquired", acquired, "eventKey", event.Key())

	// if !acquired {
	// 	return nil
	// }

	d.inner.SubmitEvent(ctx, event)
	return nil
}

func newDedupedEventSink(inner *eventBackgroundService, atomic Atomic, log *zap.SugaredLogger) EventSink {
	log.Debug("calling NewDedupedEventSink")
	return &DeduplicationEventSink{
		inner:       inner,
		dedupAtomic: NewPrefixedAtomic(atomic, "evtdedup"),
		log:         log,
	}
}
