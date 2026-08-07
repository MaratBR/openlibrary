package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createRootLogger(lc fx.Lifecycle) *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	log := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zap.DebugLevel,
		),
	)
	zap.ReplaceGlobals(log)
	lc.Append(fx.StopHook(func() error {
		_ = log.Sync()
		return nil
	}))
	return log
}
