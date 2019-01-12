package log

import (
	"io"
)

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Level() Level
	SetLevel(level Level)
	Flags() int
	SetFlags(flags int)
	SetOutput(output io.Writer)

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	WithFields([]*Field) Logger

	//With return a new Logger with appending fields
	//keyValues is key1, value1, key2, value2, ...
	//key must be convertible to string
	With(keyValues ...interface{}) Logger
	Derive() Logger
}

var defaultLogger Logger

func Default() Logger {
	return defaultLogger
}

func SetDefault(l Logger) {
	defaultLogger = l
}

func SetFlags(flags int) {
	defaultLogger.SetFlags(flags)
}

func Flags() int {
	return defaultLogger.Flags()
}

func WithFields(fields []*Field) Logger {
	return defaultLogger.WithFields(fields)
}

func With(keyValues ...interface{}) Logger {
	return defaultLogger.With(keyValues...)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}
