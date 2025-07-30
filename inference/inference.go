// Package inference provides types and functionality for TensorZero inference operations.
// This includes request/response types, streaming support, and inference configuration.
package inference

import (
	"github.com/denkhaus/tensorzero/filter"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/tool"
	"github.com/google/uuid"
)

// FinishReason represents the reason why inference finished
type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCall      FinishReason = "tool_call"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonUnknown       FinishReason = "unknown"
)

// StoredInference represents a stored inference from the list API
type StoredInference struct {
	ID             uuid.UUID              `json:"id"`
	EpisodeID      uuid.UUID              `json:"episode_id"`
	FunctionName   string                 `json:"function_name"`
	VariantName    string                 `json:"variant_name"`
	Input          InferenceInput         `json:"input"`
	Output         interface{}            `json:"output"`
	ToolParams     *tool.ToolParams       `json:"tool_params,omitempty"`
	ProcessingTime *float64               `json:"processing_time,omitempty"`
	Timestamp      string                 `json:"timestamp"` // RFC 3339
	Tags           map[string]string      `json:"tags,omitempty"`
	MetricValues   map[string]interface{} `json:"metric_values,omitempty"`
}

// InferenceResponse represents either chat or JSON inference response
type InferenceResponse interface {
	GetInferenceID() uuid.UUID
	GetEpisodeID() uuid.UUID
	GetVariantName() string
	GetUsage() shared.Usage
	GetFinishReason() *FinishReason
	GetOriginalResponse() *string
}

// InferenceChunk represents either chat or JSON chunk
type InferenceChunk interface {
	GetInferenceID() uuid.UUID
	GetEpisodeID() uuid.UUID
	GetVariantName() string
}

// ContentBlockChunk represents streaming content chunks
type ContentBlockChunk interface {
	GetType() string
	GetID() string
}

// ChatChunk represents streaming chat chunk
type ChatChunk struct {
	InferenceID  uuid.UUID           `json:"inference_id"`
	EpisodeID    uuid.UUID           `json:"episode_id"`
	VariantName  string              `json:"variant_name"`
	Content      []ContentBlockChunk `json:"content"`
	Usage        *shared.Usage       `json:"usage,omitempty"`
	FinishReason *FinishReason       `json:"finish_reason,omitempty"`
}

func (c *ChatChunk) GetInferenceID() uuid.UUID { return c.InferenceID }
func (c *ChatChunk) GetEpisodeID() uuid.UUID   { return c.EpisodeID }
func (c *ChatChunk) GetVariantName() string    { return c.VariantName }

// JsonChunk represents streaming JSON chunk
type JsonChunk struct {
	InferenceID  uuid.UUID     `json:"inference_id"`
	EpisodeID    uuid.UUID     `json:"episode_id"`
	VariantName  string        `json:"variant_name"`
	Raw          string        `json:"raw"`
	Usage        *shared.Usage `json:"usage,omitempty"`
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
}

func (j *JsonChunk) GetInferenceID() uuid.UUID { return j.InferenceID }
func (j *JsonChunk) GetEpisodeID() uuid.UUID   { return j.EpisodeID }
func (j *JsonChunk) GetVariantName() string    { return j.VariantName }

// InferenceRequest represents an inference request
type InferenceRequest struct {
	Input                   InferenceInput           `json:"input"`
	FunctionName            *string                  `json:"function_name,omitempty"`
	ModelName               *string                  `json:"model_name,omitempty"`
	EpisodeID               *uuid.UUID               `json:"episode_id,omitempty"`
	Stream                  *bool                    `json:"stream,omitempty"`
	Params                  map[string]interface{}   `json:"params,omitempty"`
	VariantName             *string                  `json:"variant_name,omitempty"`
	Dryrun                  *bool                    `json:"dryrun,omitempty"`
	OutputSchema            map[string]interface{}   `json:"output_schema,omitempty"`
	AllowedTools            []string                 `json:"allowed_tools,omitempty"`
	AdditionalTools         []map[string]interface{} `json:"additional_tools,omitempty"`
	ToolChoice              tool.ToolChoice          `json:"tool_choice,omitempty"`
	ParallelToolCalls       *bool                    `json:"parallel_tool_calls,omitempty"`
	Internal                *bool                    `json:"internal,omitempty"`
	Tags                    map[string]string        `json:"tags,omitempty"`
	Credentials             map[string]string        `json:"credentials,omitempty"`
	CacheOptions            map[string]interface{}   `json:"cache_options,omitempty"`
	ExtraBody               []ExtraBody              `json:"extra_body,omitempty"`
	ExtraHeaders            []map[string]interface{} `json:"extra_headers,omitempty"`
	IncludeOriginalResponse *bool                    `json:"include_original_response,omitempty"`
}

// InferenceInput represents input to an inference request
type InferenceInput struct {
	Messages []shared.Message `json:"messages,omitempty"`
	System   shared.System    `json:"system,omitempty"`
}

// ChatDatapointInsert represents chat datapoint insertion
type ChatDatapointInsert struct {
	FunctionName      string            `json:"function_name"`
	Input             InferenceInput    `json:"input"`
	Output            interface{}       `json:"output,omitempty"`
	AllowedTools      []string          `json:"allowed_tools,omitempty"`
	AdditionalTools   []interface{}     `json:"additional_tools,omitempty"`
	ToolChoice        *string           `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool             `json:"parallel_tool_calls,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

func (c *ChatDatapointInsert) GetFunctionName() string { return c.FunctionName }

// JsonDatapointInsert represents JSON datapoint insertion
type JsonDatapointInsert struct {
	FunctionName string            `json:"function_name"`
	Input        InferenceInput    `json:"input"`
	Output       interface{}       `json:"output,omitempty"`
	OutputSchema interface{}       `json:"output_schema,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

func (j *JsonDatapointInsert) GetFunctionName() string { return j.FunctionName }

// ExtraBody represents a custom field to be added to the inference request body
type ExtraBody struct {
}

// ListInferencesRequest represents a request to list inferences
type ListInferencesRequest struct {
	FunctionName *string                        `json:"function_name,omitempty"`
	EpisodeID    *uuid.UUID                     `json:"episode_id,omitempty"`
	VariantName  *string                        `json:"variant_name,omitempty"`
	Filter       filter.InferenceFilterTreeNode `json:"filter,omitempty"`
	OrderBy      *shared.OrderBy                `json:"order_by,omitempty"`
	Limit        *int                           `json:"limit,omitempty"`
	Offset       *int                           `json:"offset,omitempty"`
}

// ChatInferenceResponse represents a chat inference response
type ChatInferenceResponse struct {
	InferenceID      uuid.UUID             `json:"inference_id"`
	EpisodeID        uuid.UUID             `json:"episode_id"`
	VariantName      string                `json:"variant_name"`
	Content          []shared.ContentBlock `json:"content"`
	Usage            shared.Usage          `json:"usage"`
	FinishReason     *FinishReason         `json:"finish_reason,omitempty"`
	OriginalResponse *string               `json:"original_response,omitempty"`
}

func (c *ChatInferenceResponse) GetInferenceID() uuid.UUID      { return c.InferenceID }
func (c *ChatInferenceResponse) GetEpisodeID() uuid.UUID        { return c.EpisodeID }
func (c *ChatInferenceResponse) GetVariantName() string         { return c.VariantName }
func (c *ChatInferenceResponse) GetUsage() shared.Usage         { return c.Usage }
func (c *ChatInferenceResponse) GetFinishReason() *FinishReason { return c.FinishReason }
func (c *ChatInferenceResponse) GetOriginalResponse() *string   { return c.OriginalResponse }

// JsonInferenceOutput represents JSON inference output
type JsonInferenceOutput struct {
	Raw    *string                `json:"raw,omitempty"`
	Parsed map[string]interface{} `json:"parsed,omitempty"`
}

// JsonInferenceResponse represents a JSON inference response
type JsonInferenceResponse struct {
	InferenceID      uuid.UUID           `json:"inference_id"`
	EpisodeID        uuid.UUID           `json:"episode_id"`
	VariantName      string              `json:"variant_name"`
	Output           JsonInferenceOutput `json:"output"`
	Usage            shared.Usage        `json:"usage"`
	FinishReason     *FinishReason       `json:"finish_reason,omitempty"`
	OriginalResponse *string             `json:"original_response,omitempty"`
}

func (j *JsonInferenceResponse) GetInferenceID() uuid.UUID      { return j.InferenceID }
func (j *JsonInferenceResponse) GetEpisodeID() uuid.UUID        { return j.EpisodeID }
func (j *JsonInferenceResponse) GetVariantName() string         { return j.VariantName }
func (j *JsonInferenceResponse) GetUsage() shared.Usage         { return j.Usage }
func (j *JsonInferenceResponse) GetFinishReason() *FinishReason { return j.FinishReason }
func (j *JsonInferenceResponse) GetOriginalResponse() *string   { return j.OriginalResponse }

// MockInferenceResponse is a mock implementation of InferenceResponse for testing
type MockInferenceResponse struct {
	InferenceID         uuid.UUID
	EpisodeID           uuid.UUID
	VariantName         string
	UsageVal            shared.Usage
	FinishReasonVal     *FinishReason
	OriginalResponseVal *string
}

func (m *MockInferenceResponse) GetInferenceID() uuid.UUID      { return m.InferenceID }
func (m *MockInferenceResponse) GetEpisodeID() uuid.UUID        { return m.EpisodeID }
func (m *MockInferenceResponse) GetVariantName() string         { return m.VariantName }
func (m *MockInferenceResponse) GetUsage() shared.Usage         { return m.UsageVal }
func (m *MockInferenceResponse) GetFinishReason() *FinishReason { return m.FinishReasonVal }
func (m *MockInferenceResponse) GetOriginalResponse() *string   { return m.OriginalResponseVal }
