package app

import (
	"context"
	"errors"

	"github.com/jasonlvhit/gocron"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type cronRunner struct {
	stopped chan bool
	log     *zap.SugaredLogger
}

func newCronRunner(log *zap.SugaredLogger, lc fx.Lifecycle) *cronRunner {
	r := &cronRunner{log: log}
	lc.Append(fx.Hook{
		OnStart: r.start,
		OnStop:  r.stop,
	})
	return r
}

func (r *cronRunner) start(ctx context.Context) error {
	if r.stopped != nil {
		return errors.New("already started")
	}
	r.stopped = gocron.Start()
	r.log.Debug("cronRunner has started")
	return nil
}

func (r *cronRunner) stop(ctx context.Context) error {
	if r.stopped == nil {
		return nil
	}
	gocron.Clear()
	r.stopped <- true
	r.log.Debug("cronRunner has stopped")
	return nil
}

var cronRunnerModule = fx.Module("cron_runner", fx.Provide(newCronRunner, fx.Private), fx.Invoke(func(*cronRunner) {}))
