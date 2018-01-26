package log

import (
	"io"
	"log"
	"os"
)

type Level int

const (
	TraceLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		panic("unknown level")
	}
}

const (
	Ldate         = 1 << iota                              // the date in the local time zone: 2009/01/23
	Ltime                                                  // the time in the local time zone: 01:23:23
	Lmicroseconds                                          // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                              // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                             // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                                   // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lfunction                                              // function name and line number: print:23. overrides Llongfile, Lshortfile
	LstdFlags     = Ldate | Ltime | Lshortfile | Lfunction // initial values for the standard logger
)

type Logger interface {
	SetLevel(level Level)
	Level() Level
	SetOutput(w io.Writer)
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
}

var Std Logger = NewLogger(os.Stderr, TraceLevel, LstdFlags, 3)

func NewLogger(w io.Writer, level Level, flags int, calldepth int) Logger {
	return &defaultLogger{
		level:     level,
		flags:     flags,
		logger:    log.New(w, "", flags&^Llongfile&^Lshortfile),
		calldepth: calldepth,
	}
}

func SetLevel(level Level) {
	Std.SetLevel(level)
}

func GetLevel() Level {
	return Std.Level()
}

func SetOutput(w io.Writer) {
	Std.SetOutput(w)
}

func SetFlags(flags int) {
	Std.SetFlags(flags)
}

func Trace(args ...interface{}) {
	Std.Trace(args...)
}

func Debug(args ...interface{}) {
	Std.Debug(args...)
}

func Info(args ...interface{}) {
	Std.Info(args...)
}

func Warn(args ...interface{}) {
	Std.Warn(args...)
}

func Error(args ...interface{}) {
	Std.Error(args...)
}

func Fatal(args ...interface{}) {
	Std.Fatal(args...)
}

func Panic(args ...interface{}) {
	Std.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	Std.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Std.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Std.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Std.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Std.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Std.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Std.Panicf(format, args...)
}
