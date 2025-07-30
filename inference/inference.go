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

// InferenceRequest represents an inference request to the TensorZero Gateway.
// This is the core request structure for making AI inferences through TensorZero,
// supporting both function-based and direct model calls with extensive configuration options.
type InferenceRequest struct {
	// Input contains the input data for the inference, including messages, system prompts,
	// and other parameters. The structure varies based on the function type (chat vs JSON).
	Input InferenceInput `json:"input"`

	// FunctionName is the name of the function to call as defined in your TensorZero configuration.
	// Either FunctionName or ModelName must be provided, but not both.
	// Functions provide structured inference with predefined prompts, variants, and configurations.
	FunctionName *string `json:"function_name,omitempty"`

	// ModelName allows direct model calls using a built-in passthrough function.
	// Either FunctionName or ModelName must be provided, but not both.
	// This bypasses function-level configurations and calls the model directly.
	ModelName *string `json:"model_name,omitempty"`

	// EpisodeID optionally associates this inference with an existing episode.
	// For the first inference of a new episode, leave this empty and TensorZero
	// will generate a new episode ID. Only use episode IDs returned by TensorZero.
	EpisodeID *uuid.UUID `json:"episode_id,omitempty"`

	// Stream, when set to true, enables streaming responses from the model provider.
	// This allows real-time processing of partial responses as they're generated.
	Stream *bool `json:"stream,omitempty"`

	// Params allows dynamic override of inference parameters at runtime.
	// Format: {"variant_type": {"param": value, ...}, ...}
	// Prefer setting these in configuration when possible.
	Params map[string]interface{} `json:"params,omitempty"`

	// VariantName optionally pins the request to a specific variant.
	// This is not recommended for production use and is primarily for
	// testing or debugging purposes. Let TensorZero assign variants automatically.
	VariantName *string `json:"variant_name,omitempty"`

	// Dryrun, when set to true, executes the inference without storing it to the database.
	// The gateway will still call downstream model providers. This is primarily for
	// debugging and testing and should not be used in production.
	Dryrun *bool `json:"dryrun,omitempty"`

	// OutputSchema optionally overrides the output schema for JSON functions.
	// This dynamic schema is used for output validation and sent to providers
	// that support structured outputs.
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`

	// AllowedTools specifies which tools the model is allowed to call.
	// Tools must be defined in the configuration file. Any tools provided
	// in AdditionalTools are always allowed regardless of this field.
	AllowedTools []string `json:"allowed_tools,omitempty"`

	// AdditionalTools defines tools at inference time for dynamic tool use.
	// Each tool object contains: description, name, parameters, and strict fields.
	// Prefer defining tools in configuration when possible.
	AdditionalTools []map[string]interface{} `json:"additional_tools,omitempty"`

	// ToolChoice overrides the tool choice strategy for this request.
	// Supported strategies: "none", "auto", "required", or specific tool selection.
	ToolChoice tool.ToolChoice `json:"tool_choice,omitempty"`

	// ParallelToolCalls, when true, allows multiple tool calls in a single turn.
	// Only supported by certain providers (OpenAI, Fireworks AI). Defaults to
	// the function configuration value if not specified.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// Internal indicates whether this inference is generated internally by the system
	// rather than from external requests. This helps distinguish between automated
	// and user-initiated inferences.
	Internal *bool `json:"internal,omitempty"`

	// Tags are user-provided key-value pairs to associate with the inference.
	// These can be used for tracking, categorization, and analysis purposes.
	// Example: {"user_id": "123", "session": "abc", "version": "v2.1"}
	Tags map[string]string `json:"tags,omitempty"`

	// Credentials provides dynamic API keys for model providers configured
	// with dynamic credential locations. Required when providers expect
	// credentials at inference time.
	Credentials map[string]string `json:"credentials,omitempty"`

	// CacheOptions controls inference caching behavior with options like
	// enabled mode ("write_only", "read_only", "on", "off") and max_age_s.
	CacheOptions map[string]interface{} `json:"cache_options,omitempty"`

	// ExtraBody allows modification of the request body sent to model providers.
	// This is an advanced "escape hatch" for provider-specific functionality
	// not yet implemented in TensorZero.
	ExtraBody []ExtraBody `json:"extra_body,omitempty"`

	// ExtraHeaders allows modification of request headers sent to model providers.
	// This is an advanced "escape hatch" for provider-specific functionality
	// not yet implemented in TensorZero.
	ExtraHeaders []map[string]interface{} `json:"extra_headers,omitempty"`

	// IncludeOriginalResponse, when true, includes the original model provider
	// response in the response as a string. Useful for debugging and analysis.
	IncludeOriginalResponse *bool `json:"include_original_response,omitempty"`
}

// InferenceInput represents the input data for an inference request.
// This contains the messages and system prompts that will be sent to the AI model.
type InferenceInput struct {
	// Messages is a list of conversation messages to provide to the model.
	// Each message has a role (user/assistant) and content (text or content blocks).
	// This represents the conversation history and current user input.
	Messages []shared.Message `json:"messages,omitempty"`

	// System contains the system message or prompt that provides context and instructions
	// to the model. For functions without a system schema, this should be a string.
	// For functions with a system schema, this should match the defined schema structure.
	System shared.System `json:"system,omitempty"`
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

// ListInferencesRequest represents a request to list stored inferences with filtering and pagination.
// This allows you to retrieve and analyze historical inference data for monitoring and evaluation.
type ListInferencesRequest struct {
	// FunctionName optionally filters inferences to only those from a specific function.
	// If not provided, inferences from all functions will be included.
	FunctionName *string `json:"function_name,omitempty"`

	// EpisodeID optionally filters inferences to only those from a specific episode.
	// Useful for analyzing conversation flows or multi-turn interactions.
	EpisodeID *uuid.UUID `json:"episode_id,omitempty"`

	// VariantName optionally filters inferences to only those using a specific variant.
	// Helpful for comparing performance across different model configurations.
	VariantName *string `json:"variant_name,omitempty"`

	// Filter provides advanced filtering capabilities using a tree-based filter structure.
	// Supports filtering by metrics, tags, timestamps, and other inference properties.
	Filter filter.InferenceFilterTreeNode `json:"filter,omitempty"`

	// OrderBy specifies how to sort the results. Can order by timestamp, processing time,
	// or other inference properties in ascending or descending order.
	OrderBy *shared.OrderBy `json:"order_by,omitempty"`

	// Limit optionally specifies the maximum number of inferences to return.
	// Useful for pagination and controlling response size.
	Limit *int `json:"limit,omitempty"`

	// Offset optionally specifies the number of inferences to skip before
	// starting to return results. Used for pagination in combination with Limit.
	Offset *int `json:"offset,omitempty"`
}

// ChatInferenceResponse represents the response from a chat function inference.
// This contains the generated content blocks, usage metrics, and metadata about the inference.
type ChatInferenceResponse struct {
	// InferenceID is the unique identifier (UUIDv7) assigned to this inference.
	InferenceID uuid.UUID `json:"inference_id"`

	// EpisodeID is the unique identifier of the episode this inference belongs to.
	EpisodeID uuid.UUID `json:"episode_id"`

	// VariantName is the name of the variant that was used for this inference.
	VariantName string `json:"variant_name"`

	// Content contains the generated content blocks from the model, which can include
	// text, tool calls, thoughts (for reasoning models), and other content types.
	Content []shared.ContentBlock `json:"content"`

	// Usage contains token usage metrics for this inference, including input and output tokens.
	Usage shared.Usage `json:"usage"`

	// FinishReason indicates why the inference completed (stop, length, tool_call, etc.).
	FinishReason *FinishReason `json:"finish_reason,omitempty"`

	// OriginalResponse contains the raw response from the model provider when
	// IncludeOriginalResponse was set to true in the request.
	OriginalResponse *string `json:"original_response,omitempty"`
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

// JsonInferenceResponse represents the response from a JSON function inference.
// This contains the structured output, usage metrics, and metadata about the inference.
type JsonInferenceResponse struct {
	// InferenceID is the unique identifier (UUIDv7) assigned to this inference.
	InferenceID uuid.UUID `json:"inference_id"`

	// EpisodeID is the unique identifier of the episode this inference belongs to.
	EpisodeID uuid.UUID `json:"episode_id"`

	// VariantName is the name of the variant that was used for this inference.
	VariantName string `json:"variant_name"`

	// Output contains both the raw response and parsed JSON output from the model.
	Output JsonInferenceOutput `json:"output"`

	// Usage contains token usage metrics for this inference, including input and output tokens.
	Usage shared.Usage `json:"usage"`

	// FinishReason indicates why the inference completed (stop, length, tool_call, etc.).
	FinishReason *FinishReason `json:"finish_reason,omitempty"`

	// OriginalResponse contains the raw response from the model provider when
	// IncludeOriginalResponse was set to true in the request.
	OriginalResponse *string `json:"original_response,omitempty"`
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
