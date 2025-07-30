package datapoint

import (
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/tool"
	"github.com/google/uuid"
)

// DatapointInsert represents a datapoint for insertion
type DatapointInsert interface {
	GetFunctionName() string
}

// Datapoint represents a datapoint
type Datapoint struct {
	ID           uuid.UUID                `json:"id"`
	Input        inference.InferenceInput `json:"input"`
	Output       interface{}              `json:"output"`
	DatasetName  string                   `json:"dataset_name"`
	FunctionName string                   `json:"function_name"`
	ToolParams   *tool.ToolParams         `json:"tool_params,omitempty"`
	OutputSchema interface{}              `json:"output_schema,omitempty"`
	IsCustom     bool                     `json:"is_custom"`
}

// ListDatapointsRequest represents a list datapoints request
type ListDatapointsRequest struct {
	DatasetName  string  `json:"dataset_name"`
	FunctionName *string `json:"function_name,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Offset       *int    `json:"offset,omitempty"`
}
