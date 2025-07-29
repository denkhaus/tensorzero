package types

import (
	"github.com/google/uuid"
)

// ExtraBody represents a custom field to be added to the inference request body
type ExtraBody struct {
}

// InferenceInput represents input to an inference request
type InferenceInput struct {
	Messages []Message `json:"messages,omitempty"`
	System   System    `json:"system,omitempty"`
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

// ListInferencesRequest represents a request to list inferences
type ListInferencesRequest struct {
	FunctionName *string                 `json:"function_name,omitempty"`
	EpisodeID    *uuid.UUID              `json:"episode_id,omitempty"`
	VariantName  *string                 `json:"variant_name,omitempty"`
	Filter       InferenceFilterTreeNode `json:"filter,omitempty"`
	OrderBy      *OrderBy                `json:"order_by,omitempty"`
	Limit        *int                    `json:"limit,omitempty"`
	Offset       *int                    `json:"offset,omitempty"`
}

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
	ToolChoice              ToolChoice               `json:"tool_choice,omitempty"`
	ParallelToolCalls       *bool                    `json:"parallel_tool_calls,omitempty"`
	Internal                *bool                    `json:"internal,omitempty"`
	Tags                    map[string]string        `json:"tags,omitempty"`
	Credentials             map[string]string        `json:"credentials,omitempty"`
	CacheOptions            map[string]interface{}   `json:"cache_options,omitempty"`
	ExtraBody               []ExtraBody              `json:"extra_body,omitempty"`
	ExtraHeaders            []map[string]interface{} `json:"extra_headers,omitempty"`
	IncludeOriginalResponse *bool                    `json:"include_original_response,omitempty"`
}

// StoredInference represents a stored inference from the list API
type StoredInference struct {
	ID             uuid.UUID              `json:"id"`
	EpisodeID      uuid.UUID              `json:"episode_id"`
	FunctionName   string                 `json:"function_name"`
	VariantName    string                 `json:"variant_name"`
	Input          InferenceInput         `json:"input"`
	Output         interface{}            `json:"output"`
	ToolParams     *ToolParams            `json:"tool_params,omitempty"`
	ProcessingTime *float64               `json:"processing_time,omitempty"`
	Timestamp      string                 `json:"timestamp"` // RFC 3339
	Tags           map[string]string      `json:"tags,omitempty"`
	MetricValues   map[string]interface{} `json:"metric_values,omitempty"`
}

// Datapoint represents a datapoint
type Datapoint struct {
	ID           uuid.UUID      `json:"id"`
	Input        InferenceInput `json:"input"`
	Output       interface{}    `json:"output"`
	DatasetName  string         `json:"dataset_name"`
	FunctionName string         `json:"function_name"`
	ToolParams   *ToolParams    `json:"tool_params,omitempty"`
	OutputSchema interface{}    `json:"output_schema,omitempty"`
	IsCustom     bool           `json:"is_custom"`
}

// FeedbackRequest represents a feedback request
type FeedbackRequest struct {
	MetricName  string            `json:"metric_name"`
	Value       interface{}       `json:"value"`
	InferenceID *uuid.UUID        `json:"inference_id,omitempty"`
	EpisodeID   *uuid.UUID        `json:"episode_id,omitempty"`
	Dryrun      *bool             `json:"dryrun,omitempty"`
	Internal    *bool             `json:"internal,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// DynamicEvaluationRunRequest represents a dynamic evaluation run request
type DynamicEvaluationRunRequest struct {
	Variants    map[string]string `json:"variants"`
	Tags        map[string]string `json:"tags,omitempty"`
	ProjectName *string           `json:"project_name,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
}

// DynamicEvaluationRunEpisodeRequest represents a dynamic evaluation run episode request
type DynamicEvaluationRunEpisodeRequest struct {
	RunID         uuid.UUID         `json:"run_id"`
	TaskName      *string           `json:"task_name,omitempty"`
	DatapointName *string           `json:"datapoint_name,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// ListDatapointsRequest represents a list datapoints request
type ListDatapointsRequest struct {
	DatasetName  string  `json:"dataset_name"`
	FunctionName *string `json:"function_name,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Offset       *int    `json:"offset,omitempty"`
}
