package ctxzap

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkCtxZapWithContext(b *testing.B) {
	// Create a no-op logger for benchmarking
	logger := New(zap.NewNop())

	ctx := context.Background()
	ctx = WithFields(ctx,
		zap.String("request_id", "123"),
		zap.String("user_id", "456"),
		zap.String("service", "api"),
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "Processing request",
			zap.Int("items", 5),
			zap.String("action", "update"),
		)
	}
}

func BenchmarkCtxZapEmptyContext(b *testing.B) {
	// Create a no-op logger for benchmarking
	logger := New(zap.NewNop())
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "Processing request",
			zap.Int("items", 5),
			zap.String("action", "update"),
		)
	}
}

func BenchmarkZapWithFields(b *testing.B) {
	// Baseline: standard zap logger with fields
	logger := zap.NewNop()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.With(
			zap.String("request_id", "123"),
			zap.String("user_id", "456"),
			zap.String("service", "api"),
		).Info("Processing request",
			zap.Int("items", 5),
			zap.String("action", "update"),
		)
	}
}

func BenchmarkWithFields(b *testing.B) {
	ctx := context.Background()
	fields := []zap.Field{
		zap.String("request_id", "123"),
		zap.String("user_id", "456"),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = WithFields(ctx, fields...)
	}
}

func BenchmarkFieldsFromContext(b *testing.B) {
	ctx := context.Background()
	ctx = WithFields(ctx,
		zap.String("request_id", "123"),
		zap.String("user_id", "456"),
		zap.String("service", "api"),
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = FieldsFromContext(ctx)
	}
}

func BenchmarkMergeFields(b *testing.B) {
	existing := []zap.Field{
		zap.String("request_id", "123"),
		zap.String("user_id", "456"),
		zap.String("service", "api"),
	}

	newFields := []zap.Field{
		zap.Int("items", 5),
		zap.String("action", "update"),
		zap.String("service", "api-v2"), // Override
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = MergeFields(existing, newFields)
	}
}

func BenchmarkLoggerWithManyFields(b *testing.B) {
	logger := New(zap.NewNop())

	ctx := context.Background()
	// Add many fields to context
	fields := make([]zap.Field, 20)
	for i := 0; i < 20; i++ {
		fields[i] = zap.String(string(rune('a'+i)), "value")
	}
	ctx = WithFields(ctx, fields...)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "Processing request",
			zap.Int("items", 5),
		)
	}
}

// Benchmark to compare with production logger
func BenchmarkCtxZapProduction(b *testing.B) {
	// Create a production logger that writes to discard
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel) // Only log errors to reduce output

	zapLogger, _ := config.Build()
	logger := New(zapLogger)

	ctx := context.Background()
	ctx = WithFields(ctx,
		zap.String("request_id", "123"),
		zap.String("user_id", "456"),
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Error(ctx, "Processing failed",
			zap.Int("items", 5),
			zap.Error(nil),
		)
	}
}
