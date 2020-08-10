package logger

import "go.uber.org/zap"

type ILogger interface {
	Debug(msg string, keyValues ...interface{})
	Info(msg string, keyValues ...interface{})
	Warn(msg string, keyValues ...interface{})
	Error(msg string, keyValues ...interface{})
	Fatal(msg string, keyValues ...interface{})
}

func NewLogger() (ILogger, error) {
	conf := zap.NewProductionConfig()
	conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	conf.DisableCaller = true
	conf.DisableStacktrace = true
	zapLogger, err := conf.Build(zap.AddCaller(), zap.AddCallerSkip(1))

	if err != nil {
		return nil, err
	}

	logger := zpLg{
		lg: zapLogger.Sugar(),
	}

	return &logger, nil
}
