package log

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Field struct {
	Key   string
	Value interface{}
}

var defaultLogger *Logger

func init() {
	dir := os.Getenv("LOG_DIR")
	if dir == "" {
		defaultLogger = NewLogger(os.Stderr)
	} else {
		fw := NewFileWriter(dir)
		if s := os.Getenv("LOG_ROTATE_KEEP"); s != "" {
			n, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				log.Printf("Parse LOG_ROTATE_KEEP: %v", err)
			} else {
				fw.SetRotateKeep(int(n))
			}
		}
		if s := os.Getenv("LOG_ROTATE_SIZE"); s != "" {
			n, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				log.Printf("Parse LOG_ROTATE_SIZE: %v", err)
			} else {
				fw.SetRotateSize(int(n) << 20)
			}
		}
		defaultLogger = NewLogger(fw)
	}
}

func Default() *Logger {
	return defaultLogger
}

func SetDefault(l *Logger) {
	defaultLogger = l
}

var _level = AllLevel
var _flags = LstdFlags

func SetLevel(level Level) Level {
	_level = level
	return level
}

func GetLevel() Level {
	return _level
}

func SetFlags(flags int) {
	_flags = flags
}

func Flags() int {
	return _flags
}

func GetLogger(name string) *Logger {
	return defaultLogger.Derive(name)
}

func WithFields(fields []*Field) *Logger {
	return defaultLogger.WithFields(fields)
}

func With(keyValues ...interface{}) *Logger {
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
	msg := fmt.Sprintln(args...)
	msg = msg[0 : len(msg)-1]
	l := defaultLogger
	e := newEntry(l.Flags(), PanicLevel, l.name, l.fields, msg, 2)
	panic(l.render.RenderString(e))
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
	l := defaultLogger
	e := newEntry(l.Flags(), PanicLevel, l.name, l.fields, msg, 2)
	panic(l.render.RenderString(e))
}
