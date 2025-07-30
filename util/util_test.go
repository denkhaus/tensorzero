package util

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStringPtr(t *testing.T) {
	s := "test"
	ptr := StringPtr(s)
	assert.NotNil(t, ptr)
	assert.Equal(t, s, *ptr)
}

func TestBoolPtr(t *testing.T) {
	b := true
	ptr := BoolPtr(b)
	assert.NotNil(t, ptr)
	assert.Equal(t, b, *ptr)
}

func TestIntPtr(t *testing.T) {
	i := 42
	ptr := IntPtr(i)
	assert.NotNil(t, ptr)
	assert.Equal(t, i, *ptr)
}

func TestUUIDPtr(t *testing.T) {
	u := uuid.New()
	ptr := UUIDPtr(u)
	assert.NotNil(t, ptr)
	assert.Equal(t, u, *ptr)
}

func TestFloat64Ptr(t *testing.T) {
	f := 3.14
	ptr := Float64Ptr(f)
	assert.NotNil(t, ptr)
	assert.Equal(t, f, *ptr)
}

func TestToJSONString(t *testing.T) {
	data := map[string]interface{}{
		"key": "value",
		"num": 42,
	}
	
	jsonStr, err := ToJSONString(data)
	assert.NoError(t, err)
	assert.Contains(t, jsonStr, "key")
	assert.Contains(t, jsonStr, "value")
}

func TestIsNilOrEmpty(t *testing.T) {
	// Test nil pointer
	assert.True(t, IsNilOrEmpty(nil))
	
	// Test empty string
	empty := ""
	assert.True(t, IsNilOrEmpty(&empty))
	
	// Test non-empty string
	nonEmpty := "test"
	assert.False(t, IsNilOrEmpty(&nonEmpty))
}