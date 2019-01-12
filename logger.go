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
	//Log.Println("GOPATH:", s)
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

func (l *logger) Level() Level {
	return l.level
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func SetLevel(level Level) Level {
	defaultLogger.SetLevel(level)
	return level
}

func GetLevel() Level {
	return defaultLogger.Level()
}

func (l *logger) Flags() int {
	return l.flags
}

func (l *logger) SetFlags(flags int) {
	l.flags = flags
}

func (l *logger) SetOutput(w io.Writer) {
	l.render.SetWriter(w)
}

func (l *logger) Log(level Level, callDepth int, args []interface{}) {
	if l.level > level {
		return
	}
	msg := fmt.Sprint(args...)
	err := l.render.Render(newEntry(l.flags, level, l.fields, msg, callDepth+1))
	if err != nil {
		log.Fatalf("Failed to write Log: %v", err)
	}
}

func (l *logger) Logf(level Level, callDepth int, format string, args []interface{}) {
	if l.level > level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	err := l.render.Render(newEntry(l.flags, level, l.fields, msg, callDepth+1))
	if err != nil {
		log.Fatalf("Failed to write Log: %v", err)
	}
}

func (l *logger) Trace(args ...interface{}) {
	l.Log(TraceLevel, 2, args)
}

func (l *logger) Debug(args ...interface{}) {
	l.Log(DebugLevel, 2, args)
}

func (l *logger) Info(args ...interface{}) {
	l.Log(InfoLevel, 2, args)
}

func (l *logger) Warn(args ...interface{}) {
	l.Log(WarnLevel, 2, args)
}

func (l *logger) Error(args ...interface{}) {
	l.Log(ErrorLevel, 2, args)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Log(FatalLevel, 2, args)
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
	l.Logf(TraceLevel, 2, format, args)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLevel, 2, format, args)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Logf(InfoLevel, 2, format, args)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logf(WarnLevel, 2, format, args)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLevel, 2, format, args)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logf(FatalLevel, 2, format, args)
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
		defaultLogger.Panic("keyValues should be pairs of (string, interface{})", keyValues)
	}

	fields := make([]*Field, 0, n/2)
	for i := 0; i < n/2; i++ {
		if k, ok := keyValues[2*i].(string); !ok {
			defaultLogger.Panicf("keyValues[%d] isn't convertible to string", i)
		} else if keyValues[2*i+1] == nil {
			defaultLogger.Panicf("keyValues[%d] is nil", 2*i+1)
		} else {
			fields = append(fields, &Field{k, keyValues[2*i+1]})
		}
	}

	return fields
}
