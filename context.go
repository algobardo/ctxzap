package ctxzap

import (
	"context"

	"go.uber.org/zap"
)

// contextKey is used as a key for storing fields in context
type contextKey struct{}

var fieldsKey = contextKey{}

// WithFields adds zap fields to the context. Multiple calls to WithFields
// will accumulate fields. If a field with the same key already exists,
// it will be overwritten by the new value.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}

	existingFields := FieldsFromContext(ctx)
	if len(existingFields) == 0 {
		return context.WithValue(ctx, fieldsKey, fields)
	}

	// Merge fields with existing ones
	mergedFields := MergeFields(existingFields, fields)
	return context.WithValue(ctx, fieldsKey, mergedFields)
}

// FieldsFromContext extracts all zap fields stored in the context.
// Returns an empty slice if no fields are found.
func FieldsFromContext(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(fieldsKey).([]zap.Field)
	if !ok {
		return nil
	}

	// Return a copy to prevent external modifications
	result := make([]zap.Field, len(fields))
	copy(result, fields)
	return result
}
