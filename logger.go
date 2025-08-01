package ctxzap

import (
	"context"

	"go.uber.org/zap"
)

// Logger wraps a zap.Logger to provide context-aware logging methods.
type Logger struct {
	*zap.Logger
}

// New creates a new context-aware logger from an existing zap.Logger.
func New(zapLogger *zap.Logger) *Logger {
	return &Logger{Logger: zapLogger}
}

// Debug logs a message at DebugLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Debug(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Debug(msg, allFields...)
}

// Info logs a message at InfoLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Info(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Info(msg, allFields...)
}

// Warn logs a message at WarnLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Warn(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Warn(msg, allFields...)
}

// Error logs a message at ErrorLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Error(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Error(msg, allFields...)
}

// DPanic logs a message at DPanicLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.DPanic(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.DPanic(msg, allFields...)
}

// Panic logs a message at PanicLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Panic(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Panic(msg, allFields...)
}

// Fatal logs a message at FatalLevel. The message includes fields from
// both the context and any additional fields provided.
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		l.Logger.Fatal(msg, fields...)
		return
	}

	allFields := MergeFields(contextFields, fields)
	l.Logger.Fatal(msg, allFields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithOptions clones the current Logger, applies the supplied Options,
// and returns the resulting Logger. It's safe to use concurrently.
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	return &Logger{Logger: l.Logger.WithOptions(opts...)}
}
