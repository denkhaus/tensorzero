//go:build unit

package tensorzero

import (
	"context"
	"testing"

	"github.com/denkhaus/tensorzero/datapoint"
	"github.com/denkhaus/tensorzero/evaluation"
	"github.com/denkhaus/tensorzero/feedback"
	"github.com/denkhaus/tensorzero/filter"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for interfaces

// MockGateway implements the Gateway interface
type MockGateway struct {
	InferenceFn                   func(ctx context.Context, req *inference.InferenceRequest) (inference.InferenceResponse, error)
	InferenceStreamFn             func(ctx context.Context, req *inference.InferenceRequest) (<-chan inference.InferenceChunk, <-chan error)
	FeedbackFn                    func(ctx context.Context, req *feedback.Request) (*feedback.Response, error)
	DynamicEvaluationRunFn        func(ctx context.Context, req *evaluation.RunRequest) (*evaluation.RunResponse, error)
	DynamicEvaluationRunEpisodeFn func(ctx context.Context, req *evaluation.EpisodeRequest) (*evaluation.EpisodeResponse, error)
	BulkInsertDatapointsFn        func(ctx context.Context, datasetName string, datapoints []datapoint.DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapointFn             func(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapointsFn              func(ctx context.Context, req *datapoint.ListDatapointsRequest) ([]datapoint.Datapoint, error)
	ListInferencesFn              func(ctx context.Context, req *inference.ListInferencesRequest) ([]inference.StoredInference, error)
	CloseFn                       func() error
}

func (m *MockGateway) Inference(ctx context.Context, req *inference.InferenceRequest) (inference.InferenceResponse, error) {
	return m.InferenceFn(ctx, req)
}
func (m *MockGateway) InferenceStream(ctx context.Context, req *inference.InferenceRequest) (<-chan inference.InferenceChunk, <-chan error) {
	return m.InferenceStreamFn(ctx, req)
}
func (m *MockGateway) DynamicEvaluationRun(ctx context.Context, req *evaluation.RunRequest) (*evaluation.RunResponse, error) {
	return m.DynamicEvaluationRunFn(ctx, req)
}
func (m *MockGateway) DynamicEvaluationRunEpisode(ctx context.Context, req *evaluation.EpisodeRequest) (*evaluation.EpisodeResponse, error) {
	return m.DynamicEvaluationRunEpisodeFn(ctx, req)
}
func (m *MockGateway) BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []datapoint.DatapointInsert) ([]uuid.UUID, error) {
	return m.BulkInsertDatapointsFn(ctx, datasetName, datapoints)
}
func (m *MockGateway) DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error {
	return m.DeleteDatapointFn(ctx, datasetName, datapointID)
}
func (m *MockGateway) ListDatapoints(ctx context.Context, req *datapoint.ListDatapointsRequest) ([]datapoint.Datapoint, error) {
	return m.ListDatapointsFn(ctx, req)
}
func (m *MockGateway) ListInferences(ctx context.Context, req *inference.ListInferencesRequest) ([]inference.StoredInference, error) {
	return m.ListInferencesFn(ctx, req)
}
func (m *MockGateway) Close() error {
	return m.CloseFn()
}

func (m *MockGateway) Feedback(ctx context.Context, req *feedback.Request) (*feedback.Response, error) {
	return m.FeedbackFn(ctx, req)
}

// MockContentBlock implements the ContentBlock interface
type MockContentBlock struct {
	Type string
	Map  map[string]interface{}
}

func (m *MockContentBlock) GetType() string {
	return m.Type
}
func (m *MockContentBlock) ToMap() map[string]interface{} {
	return m.Map
}

// MockInferenceFilterTreeNode implements the InferenceFilterTreeNode interface
type MockInferenceFilterTreeNode struct {
	Type string
}

func (m *MockInferenceFilterTreeNode) GetType() string {
	return m.Type
}

// MockSystem implements the System interface (empty interface, so any type can implement it)
type MockSystem struct{}

// MockDatapointInsert implements the DatapointInsert interface
type MockDatapointInsert struct {
	FunctionName string
}

func (m *MockDatapointInsert) GetFunctionName() string {
	return m.FunctionName
}

// MockOptimizationJobHandle implements the OptimizationJobHandle interface
type MockOptimizationJobHandle struct {
	Type  string
	JobID string
}

func (m *MockOptimizationJobHandle) GetType() string {
	return m.Type
}
func (m *MockOptimizationJobHandle) GetJobID() string {
	return m.JobID
}

// Test cases for each interface

func TestGatewayInterface(t *testing.T) {
	mockGateway := &MockGateway{
		InferenceFn: func(ctx context.Context, req *inference.InferenceRequest) (inference.InferenceResponse, error) {
			return &inference.MockInferenceResponse{}, nil
		},
		FeedbackFn: func(ctx context.Context, req *feedback.Request) (*feedback.Response, error) {
			return &feedback.Response{}, nil
		},
		CloseFn: func() error {
			return nil
		},
	}
	var _ Gateway = mockGateway // Assert that MockGateway implements Gateway
	_, err := mockGateway.Inference(context.Background(), &inference.InferenceRequest{})
	assert.NoError(t, err)
	_, err = mockGateway.Feedback(context.Background(), &feedback.Request{})
	assert.NoError(t, err)
	assert.NoError(t, mockGateway.Close())
}

func TestContentBlockInterface(t *testing.T) {
	mockBlock := &MockContentBlock{Type: "test_type", Map: map[string]interface{}{"key": "value"}}
	var _ shared.ContentBlock = mockBlock // Assert that MockContentBlock implements ContentBlock
	assert.Equal(t, "test_type", mockBlock.GetType())
	assert.Equal(t, map[string]interface{}{"key": "value"}, mockBlock.ToMap())
}

func TestInferenceFilterTreeNodeInterface(t *testing.T) {
	mockNode := &MockInferenceFilterTreeNode{Type: "filter_type"}
	var _ filter.InferenceFilterTreeNode = mockNode // Assert that MockInferenceFilterTreeNode implements InferenceFilterTreeNode
	assert.Equal(t, "filter_type", mockNode.GetType())
}

func TestSystemInterface(t *testing.T) {
	mockSystem := &MockSystem{}
	var _ shared.System = mockSystem // Assert that MockSystem implements System (empty interface)
	assert.NotNil(t, mockSystem)
}

func TestDatapointInsertInterface(t *testing.T) {
	mockDatapoint := &MockDatapointInsert{FunctionName: "test_func"}
	var _ datapoint.DatapointInsert = mockDatapoint // Assert that MockDatapointInsert implements DatapointInsert
	assert.Equal(t, "test_func", mockDatapoint.GetFunctionName())
}

func TestOptimizationJobHandleInterface(t *testing.T) {
	mockHandle := &MockOptimizationJobHandle{Type: "job_handle", JobID: "job_123"}
	var _ OptimizationJobHandle = mockHandle // Assert that MockOptimizationJobHandle implements OptimizationJobHandle
	assert.Equal(t, "job_handle", mockHandle.GetType())
	assert.Equal(t, "job_123", mockHandle.GetJobID())
}
