//go:build unit

package filter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatMetricFilter(t *testing.T) {
	// Test NewFloatMetricFilter
	filter := NewFloatMetricFilter("accuracy", 0.95, ">=")
	assert.NotNil(t, filter)
	assert.Equal(t, "accuracy", filter.MetricName)
	assert.Equal(t, 0.95, filter.Value)
	assert.Equal(t, ">=", filter.ComparisonOperator)
	assert.Equal(t, "float_metric", filter.Type)
	assert.Equal(t, "float_metric", filter.GetType())

	// Test JSON unmarshalling
	jsonStr := `{"metric_name": "latency", "value": 100.5, "comparison_operator": "<", "type": "float_metric"}`
	var fmf FloatMetricFilter
	err := json.Unmarshal([]byte(jsonStr), &fmf)
	assert.NoError(t, err)
	assert.Equal(t, "latency", fmf.MetricName)
	assert.Equal(t, 100.5, fmf.Value)
	assert.Equal(t, "<", fmf.ComparisonOperator)
	assert.Equal(t, "float_metric", fmf.Type)
}

func TestBooleanMetricFilter(t *testing.T) {
	// Test NewBooleanMetricFilter
	filter := NewBooleanMetricFilter("success", true)
	assert.NotNil(t, filter)
	assert.Equal(t, "success", filter.MetricName)
	assert.Equal(t, true, filter.Value)
	assert.Equal(t, "boolean_metric", filter.Type)
	assert.Equal(t, "boolean_metric", filter.GetType())

	// Test JSON unmarshalling
	jsonStr := `{"metric_name": "passed", "value": false, "type": "boolean_metric"}`
	var bmf BooleanMetricFilter
	err := json.Unmarshal([]byte(jsonStr), &bmf)
	assert.NoError(t, err)
	assert.Equal(t, "passed", bmf.MetricName)
	assert.Equal(t, false, bmf.Value)
	assert.Equal(t, "boolean_metric", bmf.Type)
}

func TestTagFilter(t *testing.T) {
	// Test NewTagFilter
	filter := NewTagFilter("env", "production", "=")
	assert.NotNil(t, filter)
	assert.Equal(t, "env", filter.Key)
	assert.Equal(t, "production", filter.Value)
	assert.Equal(t, "=", filter.ComparisonOperator)
	assert.Equal(t, "tag", filter.Type)
	assert.Equal(t, "tag", filter.GetType())

	// Test JSON unmarshalling
	jsonStr := `{"key": "user_id", "value": "123", "comparison_operator": "!=", "type": "tag"}`
	var tf TagFilter
	err := json.Unmarshal([]byte(jsonStr), &tf)
	assert.NoError(t, err)
	assert.Equal(t, "user_id", tf.Key)
	assert.Equal(t, "123", tf.Value)
	assert.Equal(t, "!=", tf.ComparisonOperator)
	assert.Equal(t, "tag", tf.Type)
}

func TestTimeFilter(t *testing.T) {
	// Test NewTimeFilter
	filter := NewTimeFilter("2023-01-01T00:00:00Z", ">")
	assert.NotNil(t, filter)
	assert.Equal(t, "2023-01-01T00:00:00Z", filter.Time)
	assert.Equal(t, ">", filter.ComparisonOperator)
	assert.Equal(t, "time", filter.Type)
	assert.Equal(t, "time", filter.GetType())

	// Test JSON unmarshalling
	jsonStr := `{"time": "2023-01-02T12:00:00Z", "comparison_operator": "<=", "type": "time"}`
	var tmf TimeFilter
	err := json.Unmarshal([]byte(jsonStr), &tmf)
	assert.NoError(t, err)
	assert.Equal(t, "2023-01-02T12:00:00Z", tmf.Time)
	assert.Equal(t, "<=", tmf.ComparisonOperator)
	assert.Equal(t, "time", tmf.Type)
}

// To test AndFilter, OrFilter, and NotFilter with JSON unmarshalling, we need to
// implement a custom unmarshaller for InferenceFilterTreeNode, as it's an interface.
// For simplicity in this test, we'll only test the constructor and GetType method.
// A full test would require a custom UnmarshalJSON for InferenceFilterTreeNode
// that can correctly determine the concrete type based on the "type" field.

func TestAndFilter(t *testing.T) {
	floatFilter := NewFloatMetricFilter("score", 0.8, ">")
	tagFilter := NewTagFilter("source", "web", "=")
	filter := NewAndFilter(floatFilter, tagFilter)
	assert.NotNil(t, filter)
	assert.Len(t, filter.Children, 2)
	assert.Equal(t, "and", filter.Type)
	assert.Equal(t, "and", filter.GetType())
}

func TestOrFilter(t *testing.T) {
	floatFilter := NewFloatMetricFilter("score", 0.5, "<")
	tagFilter := NewTagFilter("source", "mobile", "=")
	filter := NewOrFilter(floatFilter, tagFilter)
	assert.NotNil(t, filter)
	assert.Len(t, filter.Children, 2)
	assert.Equal(t, "or", filter.Type)
	assert.Equal(t, "or", filter.GetType())
}

func TestNotFilter(t *testing.T) {
	floatFilter := NewFloatMetricFilter("score", 0.7, "=")
	filter := NewNotFilter(floatFilter)
	assert.NotNil(t, filter)
	assert.Equal(t, floatFilter, filter.Child)
	assert.Equal(t, "not", filter.Type)
	assert.Equal(t, "not", filter.GetType())
}
