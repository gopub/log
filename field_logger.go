package log

import "fmt"

type FieldLogger interface {
	Logger
	WithFields(fields []*Field) FieldLogger
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

func (l *fieldLogger) WithFields(fields []*Field) FieldLogger {
	return &fieldLogger{
		Logger: l.Logger,
		fields: append(l.fields, fields...),
	}
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
