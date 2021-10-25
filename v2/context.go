package log

import "context"

const keyLogger = "_logger"

func BuildContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(keyLogger).(*Logger)
	if ok {
		return l
	}
	return defaultLogger
}
