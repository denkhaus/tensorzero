// Package util provides common utility functions used across the TensorZero client.
// This includes pointer helpers, type conversions, and common operations.
package util

import (
	"encoding/json"
	"github.com/google/uuid"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// UUIDPtr returns a pointer to the given UUID
func UUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

// Float64Ptr returns a pointer to the given float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// ToJSONString safely converts any value to JSON string
func ToJSONString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsNilOrEmpty checks if a string pointer is nil or empty
func IsNilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}
