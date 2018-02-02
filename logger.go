package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var PackagePath = func() string {
	s := os.Getenv("GOPATH")
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ""
	}
	s = strings.TrimSuffix(s, "/")
	//log.Println("GOPATH:", s)
	return s + "/src/"
}()

const GoSrc = "/go/src/"

type syncWriter struct {
	mu sync.Mutex // ensures atomic writes; protects the following fields
	w  io.Writer  // destination for syncWriter
}

func (o *syncWriter) WriteAll(b []byte) error {
	o.mu.Lock()
	n, err := o.w.Write(b)
	for n < len(b) && err == nil {
		b = b[n:]
		n, err = o.w.Write(b)
	}
	o.mu.Unlock()
	return err
}

func (o *syncWriter) Write(b []byte) (int, error) {
	o.mu.Lock()
	n, err := o.w.Write(b)
	o.mu.Unlock()
	return n, err
}

type Logger interface {
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

	PrintEntry(entry *Entry)
}

func NewLogger(output io.Writer, level Level, flags int, callDepth int) Logger {
	l := &logger{
		level:     level,
		flags:     flags,
		ew:        &EntryTextPrinter{},
		callDepth: callDepth,
	}
	l.output.w = output
	return l
}

//logger is the default implementation of logger interface
type logger struct {
	level     Level
	flags     int
	callDepth int
	ew        EntryPrinter

	output struct {
		mu sync.Mutex
		w  io.Writer
	}
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) Level() Level {
	return l.level
}

func (l *logger) SetOutput(w io.Writer) {
	l.output.w = w
}

func (l *logger) SetFlags(flags int) {
	l.flags = flags
}

func (l *logger) SetEntryWriter(w EntryPrinter) {
	l.ew = w
}

func (l *logger) Trace(args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, TraceLevel, nil, msg, 2))
}

func (l *logger) Debug(args ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, DebugLevel, nil, msg, 2))
}

func (l *logger) Info(args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, InfoLevel, nil, msg, 2))
}

func (l *logger) Warn(args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, WarnLevel, nil, msg, 2))
}

func (l *logger) Error(args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, ErrorLevel, nil, msg, 2))
}

func (l *logger) Fatal(args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, FatalLevel, nil, msg, 2))
}

func (l *logger) Panic(args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, PanicLevel, nil, msg, 2))
	panic(msg)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, TraceLevel, nil, msg, 2))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.level > DebugLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, DebugLevel, nil, msg, 2))
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, InfoLevel, nil, msg, 2))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, WarnLevel, nil, msg, 2))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, ErrorLevel, nil, msg, 2))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, FatalLevel, nil, msg, 2))
}

func (l *logger) Panicf(format string, args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, PanicLevel, nil, msg, 2))
	panic(msg)
}

func (l *logger) Write(p []byte) (int, error) {
	l.output.mu.Lock()
	defer l.output.mu.Unlock()
	return l.output.w.Write(p)
}

func (l *logger) PrintEntry(entry *Entry) {
	if l.level > entry.Level {
		return
	}
	l.ew.Print(entry, l)
}
