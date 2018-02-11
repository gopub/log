package log

import (
	"fmt"
)

type FieldLogger interface {
	Logger

	WithFields([]*Field) FieldLogger

	//With return a new FieldLogger with appending fields
	//keyValues is key1, value1, key2, value2, ...
	//key must be convertible to string
	With(keyValues ...interface{}) FieldLogger
}

type fieldLogger struct {
	Logger
	fields []*Field
	level  Level
	flags  int
}

func NewFieldLogger(l Logger, level Level, flags int, fields []*Field) FieldLogger {
	return &fieldLogger{
		Logger: l,
		fields: fields,
		level:  level,
		flags:  flags,
	}
}

func makeFields(keyValues ...interface{}) []*Field {
	n := len(keyValues)
	if n%2 != 0 {
		std.Panic("keyValues should be pairs of (string, interface{})")
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

func (l *fieldLogger) WithFields(fields []*Field) FieldLogger {
	lo := &fieldLogger{
		Logger: l.Logger,
		level:  l.level,
		flags:  l.flags,
	}

	//in case of overlapping after multiple WithFields invokes
	lo.fields = make([]*Field, len(l.fields))
	copy(lo.fields, l.fields)
	lo.fields = append(lo.fields, fields...)
	return lo
}

func (l *fieldLogger) With(keyValues ...interface{}) FieldLogger {
	return l.WithFields(makeFields(keyValues...))
}

func (l *fieldLogger) Trace(args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, TraceLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Debug(args ...interface{}) {
	if l.level > DebugLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, DebugLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Info(args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, InfoLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Warn(args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, WarnLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Error(args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, ErrorLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Fatal(args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, FatalLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Panic(args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprint(args...)
	l.PrintEntry(MakeEntry(l.flags, PanicLevel, l.fields, msg, 2))
	panic(msg)
}

func (l *fieldLogger) Tracef(format string, args ...interface{}) {
	if l.level > TraceLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, TraceLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Debugf(format string, args ...interface{}) {
	if l.level > DebugLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, DebugLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Infof(format string, args ...interface{}) {
	if l.level > InfoLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, InfoLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Warnf(format string, args ...interface{}) {
	if l.level > WarnLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, WarnLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Errorf(format string, args ...interface{}) {
	if l.level > ErrorLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, ErrorLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Fatalf(format string, args ...interface{}) {
	if l.level > FatalLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, FatalLevel, l.fields, msg, 2))
}

func (l *fieldLogger) Panicf(format string, args ...interface{}) {
	if l.level > PanicLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.PrintEntry(MakeEntry(l.flags, PanicLevel, l.fields, msg, 2))
	panic(msg)
}
