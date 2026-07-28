package analytics

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type metricEntry struct {
	ValueSum float64
	Samples  int64
}

type metricBackgroundService struct {
	queueMx     sync.Mutex
	running     bool
	repo        MetricRepository
	done        chan struct{}
	cancel      context.CancelFunc
	log         *zap.SugaredLogger
	accumulator map[MetricType]map[int64]metricEntry
}

func (svc *metricBackgroundService) start(ctx context.Context) error {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	if svc.running {
		return errors.New("already running")
	}

	svc.running = true
	var stoppingCtx context.Context
	stoppingCtx, svc.cancel = context.WithCancel(context.Background())
	svc.done = make(chan struct{})
	svc.log.Debug("metricsBackgroundService: starting")
	go svc.run(stoppingCtx)

	return nil
}

func (svc *metricBackgroundService) stop(ctx context.Context) error {
	if !svc.running {
		return nil
	}

	svc.running = false
	svc.cancel()
	svc.log.Debug("metricsBackgroundService: requested stop, waiting for worker to finish")
	<-svc.done
	return nil
}

func (svc *metricBackgroundService) run(ctx context.Context) {

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Second * 5):
			svc.process(ctx)
		}
	}

	if svc.running {
		svc.log.Warn("unexpected behavior, context expired but running is false")
	}

	close(svc.done)
	svc.log.Debug("metricBackgroundService service stopped")
}

func (svc *metricBackgroundService) process(ctx context.Context) {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	count, err := svc.flush(ctx)
	if err != nil {
		svc.log.Errorw("failed to flush metrics", "err", err)
	} else if count > 0 {
		svc.log.Debugw("processed metrics", "count", count)
	}
	svc.clear()
}

func (svc *metricBackgroundService) flush(ctx context.Context) (int, error) {
	if len(svc.accumulator) == 0 {
		return 0, nil
	}

	if len(svc.accumulator) == 0 {
		return 0, nil
	}

	metrics := svc.constructMetrics()

	return len(metrics), svc.repo.AddMetrics(ctx, metrics)
}

func (svc *metricBackgroundService) clear() {
	svc.accumulator = make(map[MetricType]map[int64]metricEntry)
}

func (svc *metricBackgroundService) constructMetrics() []MetricRecord {
	metrics := make([]MetricRecord, 0)

	for metricType, bookMapping := range svc.accumulator {
		for bookID, entry := range bookMapping {
			metric := MetricRecord{
				BookID:  bookID,
				Value:   entry.ValueSum,
				Samples: entry.Samples,
				Type:    metricType,
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics

}

func (svc *metricBackgroundService) SubmitMetrics(ctx context.Context, metrics []MetricRecord) {
	if len(metrics) == 0 {
		svc.log.Debug("metricBackgroundService: SubmitMetrics - empty metrics array received")
		return
	}

	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	if svc.accumulator == nil {
		svc.accumulator = make(map[MetricType]map[int64]metricEntry)
	}

	var (
		ok      bool
		mapping map[int64]metricEntry
	)

	for _, metric := range metrics {
		if mapping, ok = svc.accumulator[metric.Type]; !ok {
			mapping = make(map[int64]metricEntry)
			svc.accumulator[metric.Type] = mapping
		}
		if entry, ok := mapping[metric.BookID]; ok {
			mapping[metric.BookID] = metricEntry{
				ValueSum: entry.ValueSum + metric.Value,
				Samples:  entry.Samples + metric.Samples,
			}
		} else {
			mapping[metric.BookID] = metricEntry{
				ValueSum: metric.Value,
				Samples:  metric.Samples,
			}
		}
	}

	svc.log.Debugw("metricBackgroundService: SubmitMetrics", "count", len(metrics))
}

func newMetricBackgroundService(
	lc fx.Lifecycle,
	log *zap.SugaredLogger,
	repo MetricRepository,
) *metricBackgroundService {
	svc := &metricBackgroundService{
		repo: repo,
		log:  log,
	}

	lc.Append(fx.Hook{
		OnStart: svc.start,
		OnStop:  svc.stop,
	})

	return svc
}

func newMetricSink(svc *metricBackgroundService) MetricSink {
	return svc
}
