package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

//logger is the default implementation of logger interface
type logger struct {
	level  Level
	flags  int
	render *render
	fields []*Field
}

func NewLogger(output io.Writer, level Level, flags int) Logger {
	l := &logger{
		level:  level,
		flags:  flags,
		render: newRender(output),
	}
	return l
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) SetOutput(w io.Writer) {
	l.render.SetWriter(w)
}

func (l *logger) SetFlags(flags int) {
	l.flags = flags
}

func (l *logger) log(level Level, args []interface{}) {
	if l.level > level {
		return
	}
	msg := fmt.Sprint(args...)
	err := l.render.Render(newEntry(l.flags, TraceLevel, l.fields, msg, 3))
	if err != nil {
		log.Fatalf("Failed to write log: %v", err)
	}
}

func (l *logger) logf(level Level, format string, args []interface{}) {
	if l.level > level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	err := l.render.Render(newEntry(l.flags, TraceLevel, l.fields, msg, 3))
	if err != nil {
		log.Fatalf("Failed to write log: %v", err)
	}
}

func (l *logger) Trace(args ...interface{}) {
	l.log(TraceLevel, args)
}

func (l *logger) Debug(args ...interface{}) {
	l.log(DebugLevel, args)
}

func (l *logger) Info(args ...interface{}) {
	l.log(InfoLevel, args)
}

func (l *logger) Warn(args ...interface{}) {
	l.log(WarnLevel, args)
}

func (l *logger) Error(args ...interface{}) {
	l.log(ErrorLevel, args)
}

func (l *logger) Fatal(args ...interface{}) {
	l.log(FatalLevel, args)
}

func (l *logger) Panic(args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprint(args...)
	e := newEntry(l.flags, PanicLevel, l.fields, msg, 2)
	panic(l.render.RenderString(e))
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.logf(TraceLevel, format, args)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logf(FatalLevel, format, args)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	e := newEntry(l.flags, PanicLevel, l.fields, msg, 2)
	panic(l.render.RenderString(e))
}

func (l *logger) WithFields(fields []*Field) Logger {
	nl := &logger{
		level:  l.level,
		flags:  l.flags,
		render: l.render,
	}

	//in case of overlapping after multiple WithFields invokes
	nl.fields = make([]*Field, len(l.fields))
	copy(nl.fields, l.fields)
	nl.fields = append(nl.fields, fields...)
	return nl
}

func (l *logger) With(keyValues ...interface{}) Logger {
	return l.WithFields(makeFields(keyValues...))
}

func (l *logger) Derive() Logger {
	nl := &logger{
		level:  l.level,
		flags:  l.flags,
		render: l.render,
	}

	//in case of overlapping after multiple WithFields invokes
	nl.fields = make([]*Field, len(l.fields))
	copy(nl.fields, l.fields)
	return nl
}

func makeFields(keyValues ...interface{}) []*Field {
	n := len(keyValues)
	if n%2 != 0 {
		std.Panic("keyValues should be pairs of (string, interface{})", keyValues)
	}

	fields := make([]*Field, 0, n/2)
	for i := 0; i < n/2; i++ {
		if k, ok := keyValues[2*i].(string); !ok {
			std.Panicf("keyValues[%d] isn't convertible to string", i)
		} else if keyValues[2*i+1] == nil {
			std.Panicf("keyValues[%d] is nil", 2*i+1)
		} else {
			fields = append(fields, &Field{k, keyValues[2*i+1]})
		}
	}

	return fields
}
