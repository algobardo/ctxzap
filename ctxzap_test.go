package ctxzap

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithFields(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() context.Context
		fields   []zap.Field
		expected []zap.Field
	}{
		{
			name:  "add fields to empty context",
			setup: context.Background,
			fields: []zap.Field{
				zap.String("key1", "value1"),
				zap.Int("key2", 42),
			},
			expected: []zap.Field{
				zap.String("key1", "value1"),
				zap.Int("key2", 42),
			},
		},
		{
			name: "add fields to context with existing fields",
			setup: func() context.Context {
				ctx := context.Background()
				return WithFields(ctx, zap.String("existing", "value"))
			},
			fields: []zap.Field{
				zap.String("new", "field"),
			},
			expected: []zap.Field{
				zap.String("existing", "value"),
				zap.String("new", "field"),
			},
		},
		{
			name: "override existing field",
			setup: func() context.Context {
				ctx := context.Background()
				return WithFields(ctx, zap.String("key", "old"))
			},
			fields: []zap.Field{
				zap.String("key", "new"),
			},
			expected: []zap.Field{
				zap.String("key", "new"),
			},
		},
		{
			name: "empty fields slice",
			setup: func() context.Context {
				ctx := context.Background()
				return WithFields(ctx, zap.String("existing", "value"))
			},
			fields: []zap.Field{},
			expected: []zap.Field{
				zap.String("existing", "value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			ctx = WithFields(ctx, tt.fields...)

			got := FieldsFromContext(ctx)

			if len(got) != len(tt.expected) {
				t.Errorf("expected %d fields, got %d", len(tt.expected), len(got))
				return
			}

			// Create maps for comparison
			gotMap := make(map[string]interface{})
			expectedMap := make(map[string]interface{})

			for _, f := range got {
				gotMap[f.Key] = f.Interface
			}

			for _, f := range tt.expected {
				expectedMap[f.Key] = f.Interface
			}

			for k, v := range expectedMap {
				if gotMap[k] != v {
					t.Errorf("field %s: expected %v, got %v", k, v, gotMap[k])
				}
			}
		})
	}
}

func TestFieldsFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected []zap.Field
	}{
		{
			name:     "nil context",
			ctx:      nil,
			expected: nil,
		},
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: nil,
		},
		{
			name: "context with fields",
			ctx: WithFields(context.Background(),
				zap.String("key1", "value1"),
				zap.Int("key2", 42),
			),
			expected: []zap.Field{
				zap.String("key1", "value1"),
				zap.Int("key2", 42),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldsFromContext(tt.ctx)

			if tt.expected == nil && got != nil {
				t.Errorf("expected nil, got %v", got)
				return
			}

			if len(got) != len(tt.expected) {
				t.Errorf("expected %d fields, got %d", len(tt.expected), len(got))
			}
		})
	}
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name             string
		setupContext     func() context.Context
		logLevel         zapcore.Level
		message          string
		additionalFields []zap.Field
		expectedFields   map[string]interface{}
	}{
		{
			name: "info with context fields",
			setupContext: func() context.Context {
				ctx := context.Background()
				return WithFields(ctx,
					zap.String("request_id", "123"),
					zap.String("user_id", "456"),
				)
			},
			logLevel: zapcore.InfoLevel,
			message:  "test message",
			additionalFields: []zap.Field{
				zap.String("action", "test"),
			},
			expectedFields: map[string]interface{}{
				"request_id": "123",
				"user_id":    "456",
				"action":     "test",
			},
		},
		{
			name: "error with overridden field",
			setupContext: func() context.Context {
				ctx := context.Background()
				return WithFields(ctx,
					zap.String("key", "context_value"),
				)
			},
			logLevel: zapcore.ErrorLevel,
			message:  "error message",
			additionalFields: []zap.Field{
				zap.String("key", "override_value"),
			},
			expectedFields: map[string]interface{}{
				"key": "override_value",
			},
		},
		{
			name:         "debug with empty context",
			setupContext: context.Background,
			logLevel:     zapcore.DebugLevel,
			message:      "debug message",
			additionalFields: []zap.Field{
				zap.Int("count", 42),
			},
			expectedFields: map[string]interface{}{
				"count": int64(42),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an observer to capture log output
			core, observed := observer.New(tt.logLevel)
			zapLogger := zap.New(core)
			logger := New(zapLogger)

			ctx := tt.setupContext()

			// Log based on level
			switch tt.logLevel {
			case zapcore.DebugLevel:
				logger.Debug(ctx, tt.message, tt.additionalFields...)
			case zapcore.InfoLevel:
				logger.Info(ctx, tt.message, tt.additionalFields...)
			case zapcore.ErrorLevel:
				logger.Error(ctx, tt.message, tt.additionalFields...)
			}

			// Check logged entries
			entries := observed.All()
			if len(entries) != 1 {
				t.Fatalf("expected 1 log entry, got %d", len(entries))
			}

			entry := entries[0]
			if entry.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, entry.Message)
			}

			// Check fields
			contextMap := entry.ContextMap()
			for key, expectedValue := range tt.expectedFields {
				if value, ok := contextMap[key]; !ok {
					t.Errorf("expected field %q not found", key)
				} else if value != expectedValue {
					t.Errorf("field %q: expected %v, got %v", key, expectedValue, value)
				}
			}

			if len(contextMap) != len(tt.expectedFields) {
				t.Errorf("expected %d fields, got %d", len(tt.expectedFields), len(contextMap))
			}
		})
	}
}

func TestLoggerWith(t *testing.T) {
	core, observed := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := New(zapLogger)

	// Create child logger with additional fields
	childLogger := logger.With(zap.String("service", "test"))

	// Log with parent
	ctx := WithFields(context.Background(), zap.String("request_id", "123"))
	logger.Info(ctx, "parent log")

	// Log with child
	childLogger.Info(ctx, "child log")

	entries := observed.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(entries))
	}

	// Check parent log
	parentEntry := entries[0]
	parentContext := parentEntry.ContextMap()
	if parentContext["request_id"] != "123" {
		t.Errorf("parent: expected request_id=123, got %v", parentContext["request_id"])
	}
	if _, ok := parentContext["service"]; ok {
		t.Error("parent should not have service field")
	}

	// Check child log
	childEntry := entries[1]
	childContext := childEntry.ContextMap()
	if childContext["request_id"] != "123" {
		t.Errorf("child: expected request_id=123, got %v", childContext["request_id"])
	}
	if childContext["service"] != "test" {
		t.Errorf("child: expected service=test, got %v", childContext["service"])
	}
}

func TestMergeFields(t *testing.T) {
	tests := []struct {
		name     string
		existing []zap.Field
		new      []zap.Field
		expected []zap.Field
	}{
		{
			name:     "merge empty slices",
			existing: []zap.Field{},
			new:      []zap.Field{},
			expected: []zap.Field{},
		},
		{
			name:     "merge with empty existing",
			existing: []zap.Field{},
			new: []zap.Field{
				zap.String("key", "value"),
			},
			expected: []zap.Field{
				zap.String("key", "value"),
			},
		},
		{
			name: "merge with empty new",
			existing: []zap.Field{
				zap.String("key", "value"),
			},
			new: []zap.Field{},
			expected: []zap.Field{
				zap.String("key", "value"),
			},
		},
		{
			name: "merge without conflicts",
			existing: []zap.Field{
				zap.String("key1", "value1"),
			},
			new: []zap.Field{
				zap.String("key2", "value2"),
			},
			expected: []zap.Field{
				zap.String("key1", "value1"),
				zap.String("key2", "value2"),
			},
		},
		{
			name: "merge with override",
			existing: []zap.Field{
				zap.String("key", "old"),
				zap.String("other", "value"),
			},
			new: []zap.Field{
				zap.String("key", "new"),
			},
			expected: []zap.Field{
				zap.String("other", "value"),
				zap.String("key", "new"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeFields(tt.existing, tt.new)

			if len(got) != len(tt.expected) {
				t.Errorf("expected %d fields, got %d", len(tt.expected), len(got))
				return
			}

			// Create maps for comparison
			gotMap := make(map[string]interface{})
			expectedMap := make(map[string]interface{})

			for _, f := range got {
				gotMap[f.Key] = f.Interface
			}

			for _, f := range tt.expected {
				expectedMap[f.Key] = f.Interface
			}

			for k, v := range expectedMap {
				if gotMap[k] != v {
					t.Errorf("field %s: expected %v, got %v", k, v, gotMap[k])
				}
			}
		})
	}
}
