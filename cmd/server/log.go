package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createRootLogger(lc fx.Lifecycle) *zap.Logger {
	log := newRootLogger()
	zap.ReplaceGlobals(log)
	lc.Append(fx.StopHook(func() error {
		_ = log.Sync()
		return nil
	}))
	return log
}

func newRootLogger() *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zap.DebugLevel,
		),
	)
}
