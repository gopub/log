package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
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

//defaultLogger is the default implementation of logger interface
type defaultLogger struct {
	level     Level
	flags     int
	output    *syncWriter
	calldepth int
	fields    []*Field
	ew        EntryWriter
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) Level() Level {
	return l.level
}

func (l *defaultLogger) SetOutput(w io.Writer) {
	l.output.w = w
}

func (l *defaultLogger) SetFlags(flags int) {
	l.flags = flags
}

func (l *defaultLogger) SetEntryWriter(w EntryWriter) {
	l.ew = w
}

func (l *defaultLogger) Trace(args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	l.print(TraceLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Debug(args ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	l.print(DebugLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Info(args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	l.print(InfoLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Warn(args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	l.print(WarnLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Error(args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	l.print(ErrorLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Fatal(args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	l.print(FatalLevel, fmt.Sprint(args...))
}

func (l *defaultLogger) Panic(args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.print(PanicLevel, msg)
	panic(msg)
}

func (l *defaultLogger) Tracef(format string, args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	l.print(TraceLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	l.print(DebugLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	l.print(InfoLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	l.print(WarnLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	l.print(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Fatalf(format string, args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	l.print(FatalLevel, fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Panicf(format string, args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.print(PanicLevel, msg)
	panic(msg)
}

func (l *defaultLogger) WithFields(fields []*Field) FieldLogger {
	nl := &defaultLogger{}
	*nl = *l
	nl.fields = append(nl.fields, fields...)
	return nl
}

func (l *defaultLogger) print(level Level, msg string) {
	var entry Entry
	if l.flags&(Ltime|Ldate|Lmicroseconds) != 0 {
		entry.Time = time.Now()
		if l.flags&LUTC != 0 {
			entry.Time = entry.Time.UTC()
		}
	}

	if l.flags&(Llongfile|Lshortfile|Lfunction) != 0 {
		function, file, line, _ := runtime.Caller(l.calldepth)
		entry.Line = line
		if l.flags&(Llongfile|Lshortfile) != 0 {
			if len(PackagePath) > 0 {
				file = strings.TrimPrefix(file, PackagePath)
			} else {
				start := strings.Index(file, GoSrc)
				if start > 0 {
					start += len(GoSrc)
				}
				file = file[start:]
			}

			if l.flags&Lshortfile != 0 {
				names := strings.Split(file, "/")
				for i := 1; i < len(names)-1; i++ {
					names[i] = names[i][0:1]
				}
				file = strings.Join(names, "/")
			}
		} else {
			file = ""
		}

		entry.File = file
		if l.flags&Lfunction != 0 {
			entry.Function = runtime.FuncForPC(function).Name()
			if len(file) > 0 {
				i := strings.LastIndex(entry.Function, ".")
				if i >= 0 {
					entry.Function = entry.Function[i+1:]
				}
			}
		}
	}

	entry.Flags = l.flags
	entry.Message = msg
	l.ew.Write(&entry, l.output)
}
