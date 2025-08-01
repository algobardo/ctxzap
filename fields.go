package ctxzap

import "go.uber.org/zap"

// MergeFields merges two slices of zap fields. If fields with the same key
// exist in both slices, the fields from the second slice take precedence.
// This ensures that newer fields can override older ones.
func MergeFields(existingFields, newFields []zap.Field) []zap.Field {
	if len(existingFields) == 0 {
		return newFields
	}
	if len(newFields) == 0 {
		return existingFields
	}

	// Create a map to track field keys for deduplication
	fieldMap := make(map[string]zap.Field, len(existingFields)+len(newFields))

	// Add existing fields first
	for _, field := range existingFields {
		fieldMap[field.Key] = field
	}

	// Add new fields, overwriting any existing ones with the same key
	for _, field := range newFields {
		fieldMap[field.Key] = field
	}

	// Convert map back to slice
	result := make([]zap.Field, 0, len(fieldMap))
	// First, add fields from existingFields that weren't overridden
	for _, field := range existingFields {
		if f, exists := fieldMap[field.Key]; exists && f.Equals(field) {
			result = append(result, field)
			delete(fieldMap, field.Key)
		}
	}
	// Then add new fields
	for _, field := range newFields {
		if _, exists := fieldMap[field.Key]; exists {
			result = append(result, field)
			delete(fieldMap, field.Key)
		}
	}

	return result
}
