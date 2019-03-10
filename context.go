package log

import "context"

const keyLogger = "logger"

func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func ContextLogger(ctx context.Context) Logger {
	l, ok := ctx.Value(keyLogger).(Logger)
	if ok {
		return l
	}
	return defaultLogger
}
