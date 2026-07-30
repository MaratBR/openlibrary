package analytics

import (
	"context"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var popularityModule = fx.Module(
	"ol_analytics_popularity",
	fx.Provide(
		newPopularityWorker,
		fx.Private,
	),
	fx.Invoke(func(*popularityWorker) {}),
)

type popularityWorker struct {
	db     dal.DB
	mx     sync.Mutex
	cancel context.CancelFunc
	log    *zap.SugaredLogger
}

func (w *popularityWorker) start(ctx context.Context) error {
	var ctx2 context.Context
	ctx2, w.cancel = context.WithCancel(context.Background())
	now := time.Now().Add(time.Second)
	return gocron.Every(1).Hour().From(&now).Do(w.run, ctx2)
}

func (w *popularityWorker) run(ctx context.Context) {

	w.mx.Lock()
	defer w.mx.Unlock()

	w.recalculate(ctx)
}

func (w *popularityWorker) recalculate(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "popularityWorker.recalculate")
	defer span.End()

	w.log.Debug("initiating popularity recalculation")

	tx, err := w.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		w.log.Errorw("popularityWorker.recalculate: failed to start new transaction", "err", err)
		return
	}
	queries := store.New(tx)

	weights := DefaultPopularityWeightsValues
	bucketStartTimes := newBucketStartTime(time.Now())
	startedAt := time.Now()

	// recalculate buckets for each period
	// we also recalculate previous buckets within 2 hours mark
	for _, item := range bucketStartTimes.BucketsWithLookback(time.Hour * 2) {
		err = queries.Analytics_RecalculateBookPopularity(ctx, store.Analytics_RecalculateBookPopularityParams{
			BucketStart: pgtype.Timestamptz{Valid: true, Time: item.Start},
			BucketType:  item.Type,
			Weights:     weights.JSON(),
		})
		if err != nil {
			dal.RollbackTx(ctx, tx)
			span.RecordError(err)
			w.log.Errorw("daily book popularity recalculation failed", "err", err)
			return
		}
	}

	took := time.Since(startedAt).Microseconds()
	w.log.Infow("daily book popularity recalculate finished", "took_us", took, "weights", weights.weights, "bucket_starts", bucketStartTimes)

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		w.log.Errorw("daily book popularity recalculation failed", "err", err)
	}
}

func (w *popularityWorker) stop(ctx context.Context) error {
	gocron.Remove(w.run)
	w.cancel()

	// hacky way to wait for job to finish
	w.mx.Lock()
	defer w.mx.Unlock()
	return nil
}

func newPopularityWorker(db dal.DB, log *zap.SugaredLogger, lc fx.Lifecycle) *popularityWorker {
	svc := &popularityWorker{db: db, log: log}

	lc.Append(fx.Hook{
		OnStart: svc.start,
		OnStop:  svc.stop,
	})

	return svc
}
