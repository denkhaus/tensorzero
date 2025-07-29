package types

import (
	"fmt"
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
	Data              interface{} `json:"data"`
	ModelProviderName *string     `json:"model_provider_name,omitempty"`
	Type              string      `json:"type"`
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
	Role    string         `json:"role"` // "user" or "assistant"
	Content []ContentBlock `json:"content"`
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
	ToolsAvailable    []Tool `json:"tools_available"`
	ToolChoice        string `json:"tool_choice"`
	ParallelToolCalls *bool  `json:"parallel_tool_calls,omitempty"`
}

// ExtraBody represents extra body content for requests

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

// OrderBy specifies ordering for list inferences
type OrderBy struct {
	By        string  `json:"by"`             // "timestamp" or "metric"
	Name      *string `json:"name,omitempty"` // metric name if by="metric"
	Direction string  `json:"direction"`      // "ASC" or "DESC"
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
