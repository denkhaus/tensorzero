//go:build unit

package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTensorZeroError(t *testing.T) {
	t.Run("Error with status code", func(t *testing.T) {
		err := NewTensorZeroError("API error", http.StatusBadRequest)
		assert.Equal(t, "TensorZero API error (status 400): API error", err.Error())
		assert.Equal(t, "API error", err.Message)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Error without status code", func(t *testing.T) {
		err := &TensorZeroError{Message: "Generic error"}
		assert.Equal(t, "TensorZero error: Generic error", err.Error())
	})

	t.Run("Error with request ID", func(t *testing.T) {
		err := &TensorZeroError{
			Message:   "API error",
			RequestID: "req-123",
		}
		assert.Equal(t, "TensorZero error: API error", err.Error())
		assert.Equal(t, "req-123", err.RequestID)
	})
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("email", "invalid format")
	assert.Equal(t, "validation error for field 'email': invalid format", err.Error())
	assert.Equal(t, "email", err.Field)
	assert.Equal(t, "invalid format", err.Message)
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "500 error is retryable",
			err:      NewTensorZeroError("Server error", http.StatusInternalServerError),
			expected: true,
		},
		{
			name:     "429 error is retryable",
			err:      NewTensorZeroError("Rate limited", http.StatusTooManyRequests),
			expected: true,
		},
		{
			name:     "400 error is not retryable",
			err:      NewTensorZeroError("Bad request", http.StatusBadRequest),
			expected: false,
		},
		{
			name:     "Non-TensorZero error is not retryable",
			err:      NewValidationError("field", "error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsRetryable(tt.err))
		})
	}
}