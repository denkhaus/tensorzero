package tensorzero

import (
	"context"

	"github.com/denkhaus/tensorzero/datapoint"
	"github.com/denkhaus/tensorzero/evaluation"
	"github.com/denkhaus/tensorzero/feedback"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/google/uuid"
)

// Gateway represents the base interface for TensorZero gateways
type Gateway interface {
	Inference(ctx context.Context, req *inference.InferenceRequest) (inference.InferenceResponse, error)
	InferenceStream(ctx context.Context, req *inference.InferenceRequest) (<-chan inference.InferenceChunk, <-chan error)
	Feedback(ctx context.Context, req *feedback.Request) (*feedback.Response, error)
	DynamicEvaluationRun(ctx context.Context, req *evaluation.RunRequest) (*evaluation.RunResponse, error)
	DynamicEvaluationRunEpisode(ctx context.Context, req *evaluation.EpisodeRequest) (*evaluation.EpisodeResponse, error)
	BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []datapoint.DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapoints(ctx context.Context, req *datapoint.ListDatapointsRequest) ([]datapoint.Datapoint, error)
	ListInferences(ctx context.Context, req *inference.ListInferencesRequest) ([]inference.StoredInference, error)
	Close() error
}

// OptimizationJobHandle represents an optimization job handle
type OptimizationJobHandle interface {
	GetType() string
	GetJobID() string
}
