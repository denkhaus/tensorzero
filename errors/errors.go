// Package errors provides structured error types for the TensorZero client.
// This includes API errors, validation errors, and client-side errors.
package errors

import (
	"fmt"
	"net/http"
)

// TensorZeroError represents a TensorZero API error
type TensorZeroError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
}

func (e *TensorZeroError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("TensorZero API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("TensorZero error: %s", e.Message)
}

// NewTensorZeroError creates a new TensorZero error
func NewTensorZeroError(message string, statusCode int) *TensorZeroError {
	return &TensorZeroError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// ValidationError represents a client-side validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsRetryable determines if an error is retryable
func IsRetryable(err error) bool {
	if tzErr, ok := err.(*TensorZeroError); ok {
		return tzErr.StatusCode >= 500 || tzErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}