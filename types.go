package tensorzero

import (
	"fmt"

	"github.com/google/uuid"
)

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// FinishReason represents the reason why inference finished
type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCall      FinishReason = "tool_call"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonUnknown       FinishReason = "unknown"
)

// ContentBlock represents a piece of content in a message
type ContentBlock interface {
	GetType() string
	ToMap() map[string]interface{}
}

// Text represents text content
type Text struct {
	Text      *string     `json:"text,omitempty"`
	Arguments interface{} `json:"arguments,omitempty"`
	Type      string      `json:"type"`
}

func NewText(text string) *Text {
	return &Text{
		Text: &text,
		Type: "text",
	}
}

func NewTextWithArguments(arguments interface{}) *Text {
	return &Text{
		Arguments: arguments,
		Type:      "text",
	}
}

func (t *Text) GetType() string {
	return t.Type
}

func (t *Text) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"type": t.Type,
	}
	if t.Text != nil {
		result["text"] = *t.Text
	}
	if t.Arguments != nil {
		result["arguments"] = t.Arguments
	}
	return result
}

// RawText represents raw text content
type RawText struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func NewRawText(value string) *RawText {
	return &RawText{
		Value: value,
		Type:  "raw_text",
	}
}

func (rt *RawText) GetType() string {
	return rt.Type
}

func (rt *RawText) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":  rt.Type,
		"value": rt.Value,
	}
}

// ImageBase64 represents base64-encoded image content
type ImageBase64 struct {
	Data     string `json:"data"`
	MimeType string `json:"mime_type"`
	Type     string `json:"type"`
}

func NewImageBase64(data, mimeType string) *ImageBase64 {
	return &ImageBase64{
		Data:     data,
		MimeType: mimeType,
		Type:     "image",
	}
}

func (img *ImageBase64) GetType() string {
	return img.Type
}

func (img *ImageBase64) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":      img.Type,
		"data":      img.Data,
		"mime_type": img.MimeType,
	}
}

// ImageURL represents image content from URL
type ImageURL struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type,omitempty"`
	Type     string  `json:"type"`
}

func NewImageURL(url string) *ImageURL {
	return &ImageURL{
		URL:  url,
		Type: "image",
	}
}

func NewImageURLWithMimeType(url, mimeType string) *ImageURL {
	return &ImageURL{
		URL:      url,
		MimeType: &mimeType,
		Type:     "image",
	}
}

func (img *ImageURL) GetType() string {
	return img.Type
}

func (img *ImageURL) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"type": img.Type,
		"url":  img.URL,
	}
	if img.MimeType != nil {
		result["mime_type"] = *img.MimeType
	}
	return result
}

// FileBase64 represents base64-encoded file content
type FileBase64 struct {
	Data     string `json:"data"`
	MimeType string `json:"mime_type"`
	Type     string `json:"type"`
}

func NewFileBase64(data, mimeType string) *FileBase64 {
	return &FileBase64{
		Data:     data,
		MimeType: mimeType,
		Type:     "file",
	}
}

func (f *FileBase64) GetType() string {
	return f.Type
}

func (f *FileBase64) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":      f.Type,
		"data":      f.Data,
		"mime_type": f.MimeType,
	}
}

// FileURL represents file content from URL
type FileURL struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func NewFileURL(url string) *FileURL {
	return &FileURL{
		URL:  url,
		Type: "file",
	}
}

func (f *FileURL) GetType() string {
	return f.Type
}

func (f *FileURL) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type": f.Type,
		"url":  f.URL,
	}
}

// ToolCall represents a tool call
type ToolCall struct {
	ID           string                 `json:"id"`
	RawArguments string                 `json:"raw_arguments"`
	RawName      string                 `json:"raw_name"`
	Arguments    map[string]interface{} `json:"arguments,omitempty"`
	Name         *string                `json:"name,omitempty"`
	Type         string                 `json:"type"`
}

func NewToolCall(id, rawArguments, rawName string) *ToolCall {
	return &ToolCall{
		ID:           id,
		RawArguments: rawArguments,
		RawName:      rawName,
		Type:         "tool_call",
	}
}

func (tc *ToolCall) GetType() string {
	return tc.Type
}

func (tc *ToolCall) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"type":          tc.Type,
		"id":            tc.ID,
		"raw_arguments": tc.RawArguments,
		"raw_name":      tc.RawName,
	}
	if tc.Arguments != nil {
		result["arguments"] = tc.Arguments
	}
	if tc.Name != nil {
		result["name"] = *tc.Name
	}
	return result
}

// Thought represents a thought content block
type Thought struct {
	Text      *string `json:"text,omitempty"`
	Type      string  `json:"type"`
	Signature *string `json:"signature,omitempty"`
}

func NewThought(text string) *Thought {
	return &Thought{
		Text: &text,
		Type: "thought",
	}
}

func (t *Thought) GetType() string {
	return t.Type
}

func (t *Thought) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"type": t.Type,
	}
	if t.Text != nil {
		result["text"] = *t.Text
	}
	if t.Signature != nil {
		result["signature"] = *t.Signature
	}
	return result
}

// ToolResult represents a tool result
type ToolResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
	ID     string `json:"id"`
	Type   string `json:"type"`
}

func NewToolResult(name, result, id string) *ToolResult {
	return &ToolResult{
		Name:   name,
		Result: result,
		ID:     id,
		Type:   "tool_result",
	}
}

func (tr *ToolResult) GetType() string {
	return tr.Type
}

func (tr *ToolResult) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   tr.Type,
		"name":   tr.Name,
		"result": tr.Result,
		"id":     tr.ID,
	}
}

// UnknownContentBlock represents unknown content
type UnknownContentBlock struct {
	Data                interface{} `json:"data"`
	ModelProviderName   *string     `json:"model_provider_name,omitempty"`
	Type                string      `json:"type"`
}

func NewUnknownContentBlock(data interface{}) *UnknownContentBlock {
	return &UnknownContentBlock{
		Data: data,
		Type: "unknown",
	}
}

func (ucb *UnknownContentBlock) GetType() string {
	return ucb.Type
}

func (ucb *UnknownContentBlock) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"type": ucb.Type,
		"data": ucb.Data,
	}
	if ucb.ModelProviderName != nil {
		result["model_provider_name"] = *ucb.ModelProviderName
	}
	return result
}

// Message represents a message in a conversation
type Message struct {
	Role    string        `json:"role"` // "user" or "assistant"
	Content []ContentBlock `json:"content"`
}

// System represents system content
type System interface{}

// InferenceInput represents input to an inference request
type InferenceInput struct {
	Messages []Message `json:"messages,omitempty"`
	System   System    `json:"system,omitempty"`
}

// JsonInferenceOutput represents JSON inference output
type JsonInferenceOutput struct {
	Raw    *string                `json:"raw,omitempty"`
	Parsed map[string]interface{} `json:"parsed,omitempty"`
}

// ChatInferenceResponse represents a chat inference response
type ChatInferenceResponse struct {
	InferenceID      uuid.UUID     `json:"inference_id"`
	EpisodeID        uuid.UUID     `json:"episode_id"`
	VariantName      string        `json:"variant_name"`
	Content          []ContentBlock `json:"content"`
	Usage            Usage         `json:"usage"`
	FinishReason     *FinishReason `json:"finish_reason,omitempty"`
	OriginalResponse *string       `json:"original_response,omitempty"`
}

// JsonInferenceResponse represents a JSON inference response
type JsonInferenceResponse struct {
	InferenceID      uuid.UUID            `json:"inference_id"`
	EpisodeID        uuid.UUID            `json:"episode_id"`
	VariantName      string               `json:"variant_name"`
	Output           JsonInferenceOutput  `json:"output"`
	Usage            Usage                `json:"usage"`
	FinishReason     *FinishReason        `json:"finish_reason,omitempty"`
	OriginalResponse *string              `json:"original_response,omitempty"`
}

// InferenceResponse represents either chat or JSON inference response
type InferenceResponse interface {
	GetInferenceID() uuid.UUID
	GetEpisodeID() uuid.UUID
	GetVariantName() string
	GetUsage() Usage
	GetFinishReason() *FinishReason
	GetOriginalResponse() *string
}

func (c *ChatInferenceResponse) GetInferenceID() uuid.UUID    { return c.InferenceID }
func (c *ChatInferenceResponse) GetEpisodeID() uuid.UUID      { return c.EpisodeID }
func (c *ChatInferenceResponse) GetVariantName() string       { return c.VariantName }
func (c *ChatInferenceResponse) GetUsage() Usage              { return c.Usage }
func (c *ChatInferenceResponse) GetFinishReason() *FinishReason { return c.FinishReason }
func (c *ChatInferenceResponse) GetOriginalResponse() *string { return c.OriginalResponse }

func (j *JsonInferenceResponse) GetInferenceID() uuid.UUID    { return j.InferenceID }
func (j *JsonInferenceResponse) GetEpisodeID() uuid.UUID      { return j.EpisodeID }
func (j *JsonInferenceResponse) GetVariantName() string       { return j.VariantName }
func (j *JsonInferenceResponse) GetUsage() Usage              { return j.Usage }
func (j *JsonInferenceResponse) GetFinishReason() *FinishReason { return j.FinishReason }
func (j *JsonInferenceResponse) GetOriginalResponse() *string { return j.OriginalResponse }

// ContentBlockChunk represents streaming content chunks
type ContentBlockChunk interface {
	GetType() string
	GetID() string
}

// TextChunk represents streaming text chunk
type TextChunk struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}

func (tc *TextChunk) GetType() string { return tc.Type }
func (tc *TextChunk) GetID() string   { return tc.ID }

// ToolCallChunk represents streaming tool call chunk
type ToolCallChunk struct {
	ID           string `json:"id"`
	RawArguments string `json:"raw_arguments"`
	RawName      string `json:"raw_name"`
	Type         string `json:"type"`
}

func (tcc *ToolCallChunk) GetType() string { return tcc.Type }
func (tcc *ToolCallChunk) GetID() string   { return tcc.ID }

// ThoughtChunk represents streaming thought chunk
type ThoughtChunk struct {
	ID        string  `json:"id"`
	Text      string  `json:"text"`
	Type      string  `json:"type"`
	Signature *string `json:"signature,omitempty"`
}

func (tc *ThoughtChunk) GetType() string { return tc.Type }
func (tc *ThoughtChunk) GetID() string   { return tc.ID }

// ChatChunk represents streaming chat chunk
type ChatChunk struct {
	InferenceID  uuid.UUID            `json:"inference_id"`
	EpisodeID    uuid.UUID            `json:"episode_id"`
	VariantName  string               `json:"variant_name"`
	Content      []ContentBlockChunk  `json:"content"`
	Usage        *Usage               `json:"usage,omitempty"`
	FinishReason *FinishReason        `json:"finish_reason,omitempty"`
}

// JsonChunk represents streaming JSON chunk
type JsonChunk struct {
	InferenceID  uuid.UUID     `json:"inference_id"`
	EpisodeID    uuid.UUID     `json:"episode_id"`
	VariantName  string        `json:"variant_name"`
	Raw          string        `json:"raw"`
	Usage        *Usage        `json:"usage,omitempty"`
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
}

// InferenceChunk represents either chat or JSON chunk
type InferenceChunk interface {
	GetInferenceID() uuid.UUID
	GetEpisodeID() uuid.UUID
	GetVariantName() string
}

func (c *ChatChunk) GetInferenceID() uuid.UUID { return c.InferenceID }
func (c *ChatChunk) GetEpisodeID() uuid.UUID   { return c.EpisodeID }
func (c *ChatChunk) GetVariantName() string    { return c.VariantName }

func (j *JsonChunk) GetInferenceID() uuid.UUID { return j.InferenceID }
func (j *JsonChunk) GetEpisodeID() uuid.UUID   { return j.EpisodeID }
func (j *JsonChunk) GetVariantName() string    { return j.VariantName }

// FeedbackResponse represents feedback response
type FeedbackResponse struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
}

// DynamicEvaluationRunResponse represents dynamic evaluation run response
type DynamicEvaluationRunResponse struct {
	RunID uuid.UUID `json:"run_id"`
}

// DynamicEvaluationRunEpisodeResponse represents dynamic evaluation run episode response
type DynamicEvaluationRunEpisodeResponse struct {
	EpisodeID uuid.UUID `json:"episode_id"`
}

// Tool represents a tool definition
type Tool struct {
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
	Name        string      `json:"name"`
	Strict      bool        `json:"strict"`
}

// ToolChoice represents tool choice options
type ToolChoice interface{}

// ToolParams represents tool parameters
type ToolParams struct {
	ToolsAvailable    []Tool       `json:"tools_available"`
	ToolChoice        string       `json:"tool_choice"`
	ParallelToolCalls *bool        `json:"parallel_tool_calls,omitempty"`
}

// ExtraBody represents extra body content for requests
type ExtraBody interface{}

// VariantExtraBody represents variant-specific extra body
type VariantExtraBody struct {
	VariantName string      `json:"variant_name"`
	Pointer     string      `json:"pointer"`
	Value       interface{} `json:"value,omitempty"`
	Delete      *bool       `json:"delete,omitempty"`
}

// ProviderExtraBody represents provider-specific extra body
type ProviderExtraBody struct {
	ModelProviderName string      `json:"model_provider_name"`
	Pointer           string      `json:"pointer"`
	Value             interface{} `json:"value,omitempty"`
	Delete            *bool       `json:"delete,omitempty"`
}

// ChatDatapointInsert represents chat datapoint insertion
type ChatDatapointInsert struct {
	FunctionName      string                 `json:"function_name"`
	Input             InferenceInput         `json:"input"`
	Output            interface{}            `json:"output,omitempty"`
	AllowedTools      []string               `json:"allowed_tools,omitempty"`
	AdditionalTools   []interface{}          `json:"additional_tools,omitempty"`
	ToolChoice        *string                `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool                  `json:"parallel_tool_calls,omitempty"`
	Tags              map[string]string      `json:"tags,omitempty"`
}

// JsonDatapointInsert represents JSON datapoint insertion
type JsonDatapointInsert struct {
	FunctionName string            `json:"function_name"`
	Input        InferenceInput    `json:"input"`
	Output       interface{}       `json:"output,omitempty"`
	OutputSchema interface{}       `json:"output_schema,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// TensorZeroError represents an error from TensorZero
type TensorZeroError struct {
	StatusCode int
	Text       string
}

func (e *TensorZeroError) Error() string {
	return fmt.Sprintf("TensorZeroError (status code %d): %s", e.StatusCode, e.Text)
}

// TensorZeroInternalError represents an internal error
type TensorZeroInternalError struct {
	Message string
}

func (e *TensorZeroInternalError) Error() string {
	return e.Message
}
// ============================================================================
// List Inferences API Types (matching Python SDK)
// ============================================================================

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
	Time               string `json:"time"` // RFC 3339 timestamp
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

// OrderBy specifies ordering for list inferences
type OrderBy struct {
	By        string  `json:"by"`        // "timestamp" or "metric"
	Name      *string `json:"name,omitempty"` // metric name if by="metric"
	Direction string  `json:"direction"` // "ASC" or "DESC"
}

// NewOrderByTimestamp creates ordering by timestamp
func NewOrderByTimestamp(direction string) *OrderBy {
	return &OrderBy{
		By:        "timestamp",
		Direction: direction,
	}
}

// NewOrderByMetric creates ordering by metric
func NewOrderByMetric(metricName, direction string) *OrderBy {
	return &OrderBy{
		By:        "metric",
		Name:      &metricName,
		Direction: direction,
	}
}

// ListInferencesRequest represents a request to list inferences
type ListInferencesRequest struct {
	FunctionName *string                   `json:"function_name,omitempty"`
	EpisodeID    *uuid.UUID                `json:"episode_id,omitempty"`
	VariantName  *string                   `json:"variant_name,omitempty"`
	Filter       InferenceFilterTreeNode   `json:"filter,omitempty"`
	OrderBy      *OrderBy                  `json:"order_by,omitempty"`
	Limit        *int                      `json:"limit,omitempty"`
	Offset       *int                      `json:"offset,omitempty"`
}

// StoredInference represents a stored inference from the list API
type StoredInference struct {
	ID              uuid.UUID              `json:"id"`
	EpisodeID       uuid.UUID              `json:"episode_id"`
	FunctionName    string                 `json:"function_name"`
	VariantName     string                 `json:"variant_name"`
	Input           InferenceInput         `json:"input"`
	Output          interface{}            `json:"output"`
	ToolParams      *ToolParams            `json:"tool_params,omitempty"`
	ProcessingTime  *float64               `json:"processing_time,omitempty"`
	Timestamp       string                 `json:"timestamp"` // RFC 3339
	Tags            map[string]string      `json:"tags,omitempty"`
	MetricValues    map[string]interface{} `json:"metric_values,omitempty"`
}
