package log

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	TraceLevel Level = iota
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
		return ""
	}
}

const (
	Ldate         = 1 << iota                                      // the date in the local time zone: 2009/01/23
	Ltime                                                          // the time in the local time zone: 01:23:23
	Lmicroseconds                                                  // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                                      // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                                     // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                                           // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lfunction                                                      // function name and line number: print:23. overrides Llongfile, Lshortfile
	LstdFlags     = Ldate | Lmicroseconds | Lshortfile | Lfunction // initial values for the standard logger
)

var globals = struct {
	flags        int
	level        Level
	output       io.Writer
	entryPrinter EntryPrinter
}{
	flags:        LstdFlags,
	level:        TraceLevel,
	output:       os.Stderr,
	entryPrinter: &EntryTextPrinter{},
}

type LevelSettable interface {
	SetLevel(level Level)
}

type FlagsSettable interface {
	SetFlags(flags int)
}

type OutputSettable interface {
	SetOutput(output io.Writer)
}

type EntryWriterSettable interface {
	SetEntryPrinter(w EntryPrinter)
}

var std Logger = NewLogger(globals.output, globals.level, globals.flags)

func GetStd() Logger {
	return std
}

func SetStd(l Logger) {
	std = l
}

func SetLevel(level Level) {
	globals.level = level
	if s, ok := std.(LevelSettable); ok {
		s.SetLevel(level)
	}
}

func GetLevel() Level {
	return globals.level
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

func SetFlags(flags int) {
	globals.flags = flags
	if s, ok := std.(FlagsSettable); ok {
		s.SetFlags(flags)
	}
}

func GetFlags() int {
	return globals.flags
}

func SetEntryPrinter(w EntryPrinter) {
	globals.entryPrinter = w
	if s, ok := std.(EntryWriterSettable); ok {
		s.SetEntryPrinter(w)
	}
}

func GetEntryPrinter() EntryPrinter {
	return globals.entryPrinter
}

func WithFields(fields []*Field) FieldLogger {
	return NewFieldLogger(std, globals.level, globals.flags, fields)
}

func With(keyValues ...interface{}) FieldLogger {
	return NewFieldLogger(std, globals.level, globals.flags, makeFields(keyValues...))
}

func Trace(args ...interface{}) {
	if globals.level > TraceLevel {
		return
	}
	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, TraceLevel, nil, msg, 2))
}

func Debug(args ...interface{}) {
	if globals.level > DebugLevel {
		return
	}
	msg := fmt.Sprintln(args...)

	std.PrintEntry(MakeEntry(globals.flags, DebugLevel, nil, msg, 2))
}

func Info(args ...interface{}) {
	if globals.level > InfoLevel {
		return
	}

	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, InfoLevel, nil, msg, 2))
}

func Warn(args ...interface{}) {
	if globals.level > WarnLevel {
		return
	}

	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, WarnLevel, nil, msg, 2))
}

func Error(args ...interface{}) {
	if globals.level > ErrorLevel {
		return
	}

	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, ErrorLevel, nil, msg, 2))
}

func Fatal(args ...interface{}) {
	if globals.level > FatalLevel {
		return
	}

	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, FatalLevel, nil, msg, 2))
}

func Panic(args ...interface{}) {
	if globals.level > PanicLevel {
		return
	}
	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	std.PrintEntry(MakeEntry(globals.flags, PanicLevel, nil, msg, 2))
	panic(msg)
}

func Tracef(format string, args ...interface{}) {
	if globals.level > TraceLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, TraceLevel, nil, msg, 2))
}

func Debugf(format string, args ...interface{}) {
	if globals.level > DebugLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, DebugLevel, nil, msg, 2))
}

func Infof(format string, args ...interface{}) {
	if globals.level > InfoLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, InfoLevel, nil, msg, 2))
}

func Warnf(format string, args ...interface{}) {
	if globals.level > WarnLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, WarnLevel, nil, msg, 2))
}

func Errorf(format string, args ...interface{}) {
	if globals.level > ErrorLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, ErrorLevel, nil, msg, 2))
}

func Fatalf(format string, args ...interface{}) {
	if globals.level > FatalLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, FatalLevel, nil, msg, 2))
}

func Panicf(format string, args ...interface{}) {
	if globals.level > PanicLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	std.PrintEntry(MakeEntry(globals.flags, PanicLevel, nil, msg, 2))
	panic(msg)
}
