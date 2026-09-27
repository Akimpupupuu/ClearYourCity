package core_logger

import (
	"context"
	"time"
)

type loggerKey struct{}

type Logger interface {
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Field struct {
	Key   string
	Value interface{}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func ToContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		panic("failed to get logger from context")
	}

	return log
}
