package log

import (
	"fmt"
	"io"
)

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Name() string
	SetName(name string)
	Level() Level
	SetLevel(level Level)
	Flags() int
	SetFlags(flags int)
	SetOutput(output io.Writer)

	Log(level Level, callDepth int, args []interface{})
	Logf(level Level, callDepth int, format string, args []interface{})

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
	Derive(name string) Logger
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

func GetLogger(name string) Logger {
	return defaultLogger.Derive(name)
}

func WithFields(fields []*Field) Logger {
	return defaultLogger.WithFields(fields)
}

func With(keyValues ...interface{}) Logger {
	return defaultLogger.With(keyValues...)
}

func Trace(args ...interface{}) {
	defaultLogger.Log(TraceLevel, 2, args)
}

func Debug(args ...interface{}) {
	defaultLogger.Log(DebugLevel, 2, args)
}

func Info(args ...interface{}) {
	defaultLogger.Log(InfoLevel, 2, args)
}

func Warn(args ...interface{}) {
	defaultLogger.Log(WarnLevel, 2, args)
}

func Error(args ...interface{}) {
	defaultLogger.Log(ErrorLevel, 2, args)
}

func Fatal(args ...interface{}) {
	defaultLogger.Log(FatalLevel, 2, args)
}

func Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if l, ok := defaultLogger.(*logger); ok {
		e := newEntry(l.flags, PanicLevel, l.name, l.fields, msg, 2)
		panic(l.render.RenderString(e))
	} else {
		panic(msg)
	}
}

func Tracef(format string, args ...interface{}) {
	defaultLogger.Logf(TraceLevel, 2, format, args)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Logf(DebugLevel, 2, format, args)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Logf(InfoLevel, 2, format, args)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Logf(WarnLevel, 2, format, args)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Logf(ErrorLevel, 2, format, args)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Logf(FatalLevel, 2, format, args)
}

func Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l, ok := defaultLogger.(*logger); ok {
		e := newEntry(l.flags, PanicLevel, l.name, l.fields, msg, 2)
		panic(l.render.RenderString(e))
	} else {
		panic(msg)
	}
}
