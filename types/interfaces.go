package types

import (
	"context"

	"github.com/google/uuid"
)

// Gateway represents the base interface for TensorZero gateways
type Gateway interface {
	Inference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error)
	InferenceStream(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error)
	Feedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error)
	DynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error)
	DynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error)
	BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error)
	ListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error)
	Close() error
}

// ContentBlock represents a piece of content in a message
type ContentBlock interface {
	GetType() string
	ToMap() map[string]interface{}
}

// InferenceFilterTreeNode represents the base interface for filter nodes
type InferenceFilterTreeNode interface {
	GetType() string
}

// System represents system content
type System interface{}

// ContentBlockChunk represents streaming content chunks
type ContentBlockChunk interface {
	GetType() string
	GetID() string
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

// InferenceChunk represents either chat or JSON chunk
type InferenceChunk interface {
	GetInferenceID() uuid.UUID
	GetEpisodeID() uuid.UUID
	GetVariantName() string
}

// DatapointInsert represents a datapoint for insertion
type DatapointInsert interface {
	GetFunctionName() string
}

// FunctionConfig represents a function configuration
type FunctionConfig interface {
	GetType() string
	GetVariants() VariantsConfig
}

type VariantConfig interface {
	GetType() string
}

// OptimizationConfig represents optimization configurations
type OptimizationConfig interface {
	GetType() string
}

// OptimizationJobHandle represents an optimization job handle
type OptimizationJobHandle interface {
	GetType() string
	GetJobID() string
}
