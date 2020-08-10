package logger

import (
	"go.uber.org/zap"
	"runtime/debug"
)

type zpLg struct {
	lg *zap.SugaredLogger
}

func (l zpLg) Debug(msg string, keyValues ...interface{}) {
	l.lg.Debugw(msg, keyValues...)
}

func (l zpLg) Info(msg string, keyValues ...interface{}) {
	l.lg.Infow(msg, keyValues...)
}

func (l zpLg) Warn(msg string, keyValues ...interface{}) {
	l.lg.Warnw(msg, keyValues...)
}

func (l zpLg) Error(msg string, keyValues ...interface{}) {
	l.lg.With("stacktrace", string(debug.Stack())).Errorw(msg, keyValues...)
}

func (l zpLg) Fatal(msg string, keyValues ...interface{}) {
	l.lg.Fatalw(msg, keyValues...)
}
