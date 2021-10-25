package log

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var defaultLogger *Logger

func init() {
	dir := os.Getenv("LOG_DIR")
	if dir == "" {
		defaultLogger = NewLogger(NewZapLogger(os.Stderr))
		return
	}

	fw := NewFileWriter(filepath.Join(dir, os.Args[0], "log"))

	if s := os.Getenv("LOG_ROTATE_KEEP"); s != "" {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("Parse LOG_ROTATE_KEEP: %v\n", err)
		} else {
			fw.SetMaxBackups(int(n))
		}
	}

	if s := os.Getenv("LOG_ROTATE_SIZE"); s != "" {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("Parse LOG_ROTATE_SIZE: %v\n", err)
		} else {
			fw.SetMaxSize(int(n) << 20)
		}
	}

	defaultLogger = NewLogger(NewZapLogger(fw))
}

func Default() *Logger {
	return defaultLogger
}

func SetDefault(l *Logger) {
	defaultLogger = l
}

func GetLogger(name string) *Logger {
	return defaultLogger.Derive(name)
}

func With(keyValues ...interface{}) *Logger {
	return defaultLogger.With(keyValues...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func ErrorE(err error) {
	if err == nil {
		return
	}
	Error(err)
}

func FatalE(err error) {
	if err == nil {
		return
	}
	Fatal(err)
}

func PanicE(err error) {
	if err == nil {
		return
	}
	Panic(err)
}
