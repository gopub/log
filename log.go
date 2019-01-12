package log

import (
	"io"
	"os"
)

var globals = struct {
	flags  int
	level  Level
	output io.Writer
}{
	flags:  LstdFlags,
	level:  TraceLevel,
	output: os.Stderr,
}

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	SetLevel(level Level)
	SetOutput(output io.Writer)
	SetFlags(flags int)

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

type OutputSettable interface {
	SetOutput(output io.Writer)
}

var std Logger = NewLogger(globals.output, globals.level, globals.flags)

func GetStd() Logger {
	return std
}

func SetStd(l Logger) {
	std = l
}

func SetOutput(w io.Writer) {
	globals.output = w
	if s, ok := std.(OutputSettable); ok {
		s.SetOutput(w)
	}
}

func GetOutput() io.Writer {
	return globals.output
}

func WithFields(fields []*Field) Logger {
	return std.WithFields(fields)
}

func With(keyValues ...interface{}) Logger {
	return std.With(keyValues...)
}

func Trace(args ...interface{}) {
	std.Trace(args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

func Panic(args ...interface{}) {
	std.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}
