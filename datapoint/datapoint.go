package datapoint

import (
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/tool"
	"github.com/google/uuid"
)

// DatapointInsert represents a datapoint that can be inserted into a dataset.
// Datapoints are used to store training data, evaluation data, or examples
// that can be used for testing and improving your AI functions.
type DatapointInsert interface {
	// GetFunctionName returns the name of the function this datapoint is associated with.
	// This must match a function defined in your TensorZero configuration.
	GetFunctionName() string
}

// Datapoint represents a stored datapoint in a TensorZero dataset.
// Datapoints contain input-output pairs that can be used for training,
// evaluation, testing, or as examples for your AI functions.
type Datapoint struct {
	// ID is the unique identifier (UUIDv7) for this datapoint.
	ID uuid.UUID `json:"id"`

	// Input contains the input data for the function, including messages,
	// system prompts, and other parameters that would be passed to an inference.
	Input inference.InferenceInput `json:"input"`

	// Output contains the expected or actual output for this datapoint.
	// For chat functions, this would be content blocks; for JSON functions,
	// this would be the structured output object.
	Output interface{} `json:"output"`

	// DatasetName is the name of the dataset this datapoint belongs to.
	// Datasets help organize related datapoints for specific use cases.
	DatasetName string `json:"dataset_name"`

	// FunctionName is the name of the function this datapoint is associated with.
	// This must match a function defined in your TensorZero configuration.
	FunctionName string `json:"function_name"`

	// ToolParams contains any tool-related parameters used for this datapoint,
	// such as tool choice strategies or available tools.
	ToolParams *tool.ToolParams `json:"tool_params,omitempty"`

	// OutputSchema defines the expected schema for the output when dealing
	// with JSON functions. This ensures output validation and consistency.
	OutputSchema interface{} `json:"output_schema,omitempty"`

	// IsCustom indicates whether this datapoint was created through custom
	// processes rather than standard TensorZero operations.
	IsCustom bool `json:"is_custom"`
}

// ListDatapointsRequest represents a request to list datapoints from a dataset.
// This allows you to retrieve and paginate through datapoints for analysis,
// evaluation, or other processing needs.
type ListDatapointsRequest struct {
	// DatasetName is the name of the dataset to list datapoints from.
	// This is required and must match an existing dataset.
	DatasetName string `json:"dataset_name"`

	// FunctionName optionally filters datapoints to only those associated
	// with a specific function. If not provided, datapoints for all functions
	// in the dataset will be returned.
	FunctionName *string `json:"function_name,omitempty"`

	// Limit optionally specifies the maximum number of datapoints to return.
	// This is useful for pagination and controlling response size.
	Limit *int `json:"limit,omitempty"`

	// Offset optionally specifies the number of datapoints to skip before
	// starting to return results. This is used for pagination in combination with Limit.
	Offset *int `json:"offset,omitempty"`
}
