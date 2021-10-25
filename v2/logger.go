package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	l  *zap.Logger
	sl *zap.SugaredLogger
}

func NewZapLogger(w io.Writer) *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(w),
		zap.DebugLevel,
	)
	return zap.New(core).WithOptions(zap.AddCaller())
}

func NewLogger(l *zap.Logger) *Logger {
	l = l.WithOptions(zap.AddCallerSkip(2))
	return &Logger{
		l:  l,
		sl: l.Sugar(),
	}
}

func (l *Logger) Sync() error {
	err1 := l.l.Sync()
	err2 := l.sl.Sync()
	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return nil
}

func (l *Logger) Debug(args ...interface{}) {
	l.sl.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sl.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sl.Error(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sl.Panic(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sl.Fatal(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sl.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sl.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sl.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sl.Errorf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sl.Panicf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sl.Fatalf(format, args...)
}

func (l *Logger) With(keyValues ...interface{}) *Logger {
	sl := l.sl.With(keyValues...)
	return &Logger{
		l:  sl.Desugar(),
		sl: sl,
	}
}

func (l *Logger) Named(name string) *Logger {
	nl := l.l.Named(name)
	return &Logger{
		l:  nl,
		sl: nl.Sugar(),
	}
}

func (l *Logger) Derive(name string) *Logger {
	return l.Named(name)
}

func (l *Logger) WithLevel(level Level) *Logger {
	nl := l.l.WithOptions(zap.IncreaseLevel(level))
	return &Logger{
		l:  nl,
		sl: nl.Sugar(),
	}
}
