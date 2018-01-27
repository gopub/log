package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var _pkgPath = func() string {
	s := os.Getenv("GOPATH")
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ""
	}
	s = strings.TrimSuffix(s, "/")
	//log.Println("GOPATH:", s)
	return s + "/src/"
}()

const _GoSrc = "/go/src/"

//defaultLogger is the default implementation of logger interface
type defaultLogger struct {
	level     Level
	flags     int
	logger    *log.Logger
	calldepth int
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) Level() Level {
	return l.level
}

func (l *defaultLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *defaultLogger) SetFlags(flags int) {
	l.flags = flags
	l.logger.SetFlags(flags)
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

func (l *defaultLogger) print(level Level, msg string) {
	if l.flags&(Llongfile|Lshortfile|Lfunction) != 0 {
		function, file, line, _ := runtime.Caller(l.calldepth)
		if l.flags&(Llongfile|Lshortfile) != 0 {
			if len(_pkgPath) > 0 {
				file = strings.TrimPrefix(file, _pkgPath)
			} else {
				start := strings.Index(file, _GoSrc)
				if start > 0 {
					start += len(_GoSrc)
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

		var funcName string
		if l.flags&Lfunction != 0 {
			funcName = runtime.FuncForPC(function).Name()
			if len(file) > 0 {
				i := strings.LastIndex(funcName, ".")
				if i >= 0 {
					funcName = funcName[i+1:]
				}
			}
		}

		if len(file) > 0 && len(funcName) > 0 {
			l.logger.Printf("[%s] %s(%s):%d %s", level.String(), file, funcName, line, msg)
		} else if len(file) > 0 {
			l.logger.Printf("[%s] %s:%d %s", level.String(), file, line, msg)
		} else if len(funcName) > 0 {
			l.logger.Printf("[%s] %s:%d %s", level.String(), funcName, line, msg)
		}
	} else {
		l.logger.Printf("[%s] %s", level.String(), msg)
	}
}
