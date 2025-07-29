package types

import (
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
