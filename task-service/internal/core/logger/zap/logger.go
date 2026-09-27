package core_zap_logger

import (
	"fmt"
	"time"

	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger() (core_logger.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encoderConfig,
	}

	log, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return &Logger{
		logger: log,
	}, nil
}

func (l *Logger) With(fields ...core_logger.Field) core_logger.Logger {
	zapFields := createSliceZapFields(fields...)

	log := l.logger.With(zapFields...)

	return &Logger{logger: log}
}

func (l *Logger) Error(msg string, fields ...core_logger.Field) {
	zapFields := createSliceZapFields(fields...)
	l.logger.Error(msg, zapFields...)
}

func (l *Logger) Warn(msg string, fields ...core_logger.Field) {
	zapFields := createSliceZapFields(fields...)
	l.logger.Warn(msg, zapFields...)
}

func (l *Logger) Debug(msg string, fields ...core_logger.Field) {
	zapFields := createSliceZapFields(fields...)
	l.logger.Debug(msg, zapFields...)
}

func (l *Logger) Fatal(msg string, fields ...core_logger.Field) {
	zapFields := createSliceZapFields(fields...)
	l.logger.Fatal(msg, zapFields...)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func createSliceZapFields(fields ...core_logger.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))

	for i := range fields {
		zapFields[i] = mapField(fields[i])
	}

	return zapFields
}

func mapField(f core_logger.Field) zap.Field {
	switch v := f.Value.(type) {
	case error:
		return zap.Error(v)
	case int:
		return zap.Int(f.Key, v)
	case string:
		return zap.String(f.Key, v)
	case time.Time:
		return zap.Time(f.Key, v)
	case time.Duration:
		return zap.Duration(f.Key, v)
	default:
		return zap.Any(f.Key, v)
	}
}
