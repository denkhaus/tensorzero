package filter

// InferenceFilterTreeNode represents the base interface for filter nodes
type InferenceFilterTreeNode interface {
	GetType() string
}

// FloatMetricFilter filters inferences by float metric values
type FloatMetricFilter struct {
	MetricName         string  `json:"metric_name"`
	Value              float64 `json:"value"`
	ComparisonOperator string  `json:"comparison_operator"` // "<", "<=", "=", ">", ">=", "!="
	Type               string  `json:"type"`
}

func (f *FloatMetricFilter) GetType() string { return f.Type }

// NewFloatMetricFilter creates a new float metric filter
func NewFloatMetricFilter(metricName string, value float64, operator string) *FloatMetricFilter {
	return &FloatMetricFilter{
		MetricName:         metricName,
		Value:              value,
		ComparisonOperator: operator,
		Type:               "float_metric",
	}
}

// BooleanMetricFilter filters inferences by boolean metric values
type BooleanMetricFilter struct {
	MetricName string `json:"metric_name"`
	Value      bool   `json:"value"`
	Type       string `json:"type"`
}

func (f *BooleanMetricFilter) GetType() string { return f.Type }

// NewBooleanMetricFilter creates a new boolean metric filter
func NewBooleanMetricFilter(metricName string, value bool) *BooleanMetricFilter {
	return &BooleanMetricFilter{
		MetricName: metricName,
		Value:      value,
		Type:       "boolean_metric",
	}
}

// TagFilter filters inferences by tag values
type TagFilter struct {
	Key                string `json:"key"`
	Value              string `json:"value"`
	ComparisonOperator string `json:"comparison_operator"` // "=", "!="
	Type               string `json:"type"`
}

func (f *TagFilter) GetType() string { return f.Type }

// NewTagFilter creates a new tag filter
func NewTagFilter(key, value, operator string) *TagFilter {
	return &TagFilter{
		Key:                key,
		Value:              value,
		ComparisonOperator: operator,
		Type:               "tag",
	}
}

// TimeFilter filters inferences by timestamp
type TimeFilter struct {
	Time               string `json:"time"`                // RFC 3339 timestamp
	ComparisonOperator string `json:"comparison_operator"` // "<", "<=", "=", ">", ">=", "!="
	Type               string `json:"type"`
}

func (f *TimeFilter) GetType() string { return f.Type }

// NewTimeFilter creates a new time filter
func NewTimeFilter(time, operator string) *TimeFilter {
	return &TimeFilter{
		Time:               time,
		ComparisonOperator: operator,
		Type:               "time",
	}
}

// AndFilter combines multiple filters with AND logic
type AndFilter struct {
	Children []InferenceFilterTreeNode `json:"children"`
	Type     string                    `json:"type"`
}

func (f *AndFilter) GetType() string { return f.Type }

// NewAndFilter creates a new AND filter
func NewAndFilter(children ...InferenceFilterTreeNode) *AndFilter {
	return &AndFilter{
		Children: children,
		Type:     "and",
	}
}

// OrFilter combines multiple filters with OR logic
type OrFilter struct {
	Children []InferenceFilterTreeNode `json:"children"`
	Type     string                    `json:"type"`
}

func (f *OrFilter) GetType() string { return f.Type }

// NewOrFilter creates a new OR filter
func NewOrFilter(children ...InferenceFilterTreeNode) *OrFilter {
	return &OrFilter{
		Children: children,
		Type:     "or",
	}
}

// NotFilter negates a filter
type NotFilter struct {
	Child InferenceFilterTreeNode `json:"child"`
	Type  string                  `json:"type"`
}

func (f *NotFilter) GetType() string { return f.Type }

// NewNotFilter creates a new NOT filter
func NewNotFilter(child InferenceFilterTreeNode) *NotFilter {
	return &NotFilter{
		Child: child,
		Type:  "not",
	}
}
