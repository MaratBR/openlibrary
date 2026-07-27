package analytics

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type metricBackgroundService struct {
	queueMx     sync.Mutex
	running     bool
	repo        MetricRepository
	wg          sync.WaitGroup
	cancel      context.CancelFunc
	log         *zap.SugaredLogger
	accumulator map[string]map[int64]float64
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
	svc.log.Debug("metricsBackgroundService: starting, waiting for worker to initialize")
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
	return nil
}

func (svc *metricBackgroundService) run(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Second * 5):
			err := svc.flush(ctx)
			if err != nil {
				svc.log.Errorw("failed to flush metrics", "err", err)
			}
			svc.clear()
		}
	}

	if svc.running {
		svc.log.Warn("unexpected behavior, context expired but running is false")
	}

	svc.wg.Done()
}

func (svc *metricBackgroundService) flush(ctx context.Context) error {
	if len(svc.accumulator) == 0 {
		return nil
	}

	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	if len(svc.accumulator) == 0 {
		return nil
	}

	metrics := svc.constructMetrics()

	return svc.repo.AddMetrics(ctx, metrics)
}

func (svc *metricBackgroundService) constructMetrics() []Metric {
	metrics := make([]Metric, 0)

	for metricType, bookMapping := range svc.accumulator {
		for bookID, value := range bookMapping {
			if value == 0 {
				continue
			}

			metric := Metric{
				BookID: bookID,
				Value:  value,
				Type:   metricType,
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics

}

func (svc *metricBackgroundService) clear() {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()
	svc.accumulator = make(map[string]map[int64]float64)
}

func (svc *metricBackgroundService) SubmitMetric(ctx context.Context, metric Metric) {
	svc.queueMx.Lock()
	defer svc.queueMx.Unlock()

	var (
		ok      bool
		mapping map[int64]float64
	)

	if mapping, ok = svc.accumulator[metric.Type]; !ok {
		mapping = make(map[int64]float64)
		svc.accumulator[metric.Type] = mapping
	}
	if value, ok := mapping[metric.BookID]; ok {
		mapping[metric.BookID] = value + metric.Value
	} else {
		mapping[metric.BookID] = metric.Value
	}
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
