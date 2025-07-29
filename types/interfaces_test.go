package types

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for interfaces

// MockGateway implements the Gateway interface
type MockGateway struct {
	InferenceFn                   func(ctx context.Context, req *InferenceRequest) (InferenceResponse, error)
	InferenceStreamFn             func(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error)
	FeedbackFn                    func(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error)
	DynamicEvaluationRunFn        func(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error)
	DynamicEvaluationRunEpisodeFn func(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error)
	BulkInsertDatapointsFn        func(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapointFn             func(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapointsFn              func(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error)
	ListInferencesFn              func(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error)
	CloseFn                       func() error
}

func (m *MockGateway) Inference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error) {
	return m.InferenceFn(ctx, req)
}
func (m *MockGateway) InferenceStream(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error) {
	return m.InferenceStreamFn(ctx, req)
}
func (m *MockGateway) Feedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error) {
	return m.FeedbackFn(ctx, req)
}
func (m *MockGateway) DynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error) {
	return m.DynamicEvaluationRunFn(ctx, req)
}
func (m *MockGateway) DynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error) {
	return m.DynamicEvaluationRunEpisodeFn(ctx, req)
}
func (m *MockGateway) BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error) {
	return m.BulkInsertDatapointsFn(ctx, datasetName, datapoints)
}
func (m *MockGateway) DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error {
	return m.DeleteDatapointFn(ctx, datasetName, datapointID)
}
func (m *MockGateway) ListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error) {
	return m.ListDatapointsFn(ctx, req)
}
func (m *MockGateway) ListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error) {
	return m.ListInferencesFn(ctx, req)
}
func (m *MockGateway) Close() error {
	return m.CloseFn()
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

// MockContentBlockChunk implements the ContentBlockChunk interface
type MockContentBlockChunk struct {
	Type string
	ID   string
}

func (m *MockContentBlockChunk) GetType() string {
	return m.Type
}
func (m *MockContentBlockChunk) GetID() string {
	return m.ID
}

// MockInferenceResponse implements the InferenceResponse interface
type MockInferenceResponse struct {
	InferenceID         uuid.UUID
	EpisodeID           uuid.UUID
	Variant             string
	UsageVal            Usage
	FinishReasonVal     *FinishReason
	OriginalResponseVal *string
}

func (m *MockInferenceResponse) GetInferenceID() uuid.UUID {
	return m.InferenceID
}
func (m *MockInferenceResponse) GetEpisodeID() uuid.UUID {
	return m.EpisodeID
}
func (m *MockInferenceResponse) GetVariantName() string {
	return m.Variant
}
func (m *MockInferenceResponse) GetUsage() Usage {
	return m.UsageVal
}
func (m *MockInferenceResponse) GetFinishReason() *FinishReason {
	return m.FinishReasonVal
}
func (m *MockInferenceResponse) GetOriginalResponse() *string {
	return m.OriginalResponseVal
}

// MockInferenceChunk implements the InferenceChunk interface
type MockInferenceChunk struct {
	InferenceID uuid.UUID
	EpisodeID   uuid.UUID
	Variant     string
}

func (m *MockInferenceChunk) GetInferenceID() uuid.UUID {
	return m.InferenceID
}
func (m *MockInferenceChunk) GetEpisodeID() uuid.UUID {
	return m.EpisodeID
}
func (m *MockInferenceChunk) GetVariantName() string {
	return m.Variant
}

// MockDatapointInsert implements the DatapointInsert interface
type MockDatapointInsert struct {
	FunctionName string
}

func (m *MockDatapointInsert) GetFunctionName() string {
	return m.FunctionName
}

// MockFunctionConfig implements the FunctionConfig interface
type MockFunctionConfig struct {
	Type     string
	Variants VariantsConfig
}

func (m *MockFunctionConfig) GetType() string {
	return m.Type
}
func (m *MockFunctionConfig) GetVariants() VariantsConfig {
	return m.Variants
}

// MockVariantConfig implements the VariantConfig interface
type MockVariantConfig struct {
	Type string
}

func (m *MockVariantConfig) GetType() string {
	return m.Type
}

// MockOptimizationConfig implements the OptimizationConfig interface
type MockOptimizationConfig struct {
	Type string
}

func (m *MockOptimizationConfig) GetType() string {
	return m.Type
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
		InferenceFn: func(ctx context.Context, req *InferenceRequest) (InferenceResponse, error) {
			return &MockInferenceResponse{}, nil
		},
		CloseFn: func() error {
			return nil
		},
	}
	var _ Gateway = mockGateway // Assert that MockGateway implements Gateway
	_, err := mockGateway.Inference(context.Background(), &InferenceRequest{})
	assert.NoError(t, err)
	assert.NoError(t, mockGateway.Close())
}

func TestContentBlockInterface(t *testing.T) {
	mockBlock := &MockContentBlock{Type: "test_type", Map: map[string]interface{}{"key": "value"}}
	var _ ContentBlock = mockBlock // Assert that MockContentBlock implements ContentBlock
	assert.Equal(t, "test_type", mockBlock.GetType())
	assert.Equal(t, map[string]interface{}{"key": "value"}, mockBlock.ToMap())
}

func TestInferenceFilterTreeNodeInterface(t *testing.T) {
	mockNode := &MockInferenceFilterTreeNode{Type: "filter_type"}
	var _ InferenceFilterTreeNode = mockNode // Assert that MockInferenceFilterTreeNode implements InferenceFilterTreeNode
	assert.Equal(t, "filter_type", mockNode.GetType())
}

func TestSystemInterface(t *testing.T) {
	mockSystem := &MockSystem{}
	var _ System = mockSystem // Assert that MockSystem implements System (empty interface)
	assert.NotNil(t, mockSystem)
}

func TestContentBlockChunkInterface(t *testing.T) {
	mockChunk := &MockContentBlockChunk{Type: "chunk_type", ID: "chunk_id"}
	var _ ContentBlockChunk = mockChunk // Assert that MockContentBlockChunk implements ContentBlockChunk
	assert.Equal(t, "chunk_type", mockChunk.GetType())
	assert.Equal(t, "chunk_id", mockChunk.GetID())
}

func TestInferenceResponseInterface(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonStop
	originalResponse := "original"
	mockResponse := &MockInferenceResponse{
		InferenceID:         infID,
		EpisodeID:           epID,
		Variant:             "default",
		UsageVal:            Usage{InputTokens: 10, OutputTokens: 20},
		FinishReasonVal:     &finishReason,
		OriginalResponseVal: &originalResponse,
	}
	var _ InferenceResponse = mockResponse // Assert that MockInferenceResponse implements InferenceResponse
	assert.Equal(t, infID, mockResponse.GetInferenceID())
	assert.Equal(t, epID, mockResponse.GetEpisodeID())
	assert.Equal(t, "default", mockResponse.GetVariantName())
	assert.Equal(t, Usage{InputTokens: 10, OutputTokens: 20}, mockResponse.GetUsage())
	assert.Equal(t, &finishReason, mockResponse.GetFinishReason())
	assert.Equal(t, &originalResponse, mockResponse.GetOriginalResponse())
}

func TestInferenceChunkInterface(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	mockChunk := &MockInferenceChunk{
		InferenceID: infID,
		EpisodeID:   epID,
		Variant:     "stream_variant",
	}
	var _ InferenceChunk = mockChunk // Assert that MockInferenceChunk implements InferenceChunk
	assert.Equal(t, infID, mockChunk.GetInferenceID())
	assert.Equal(t, epID, mockChunk.GetEpisodeID())
	assert.Equal(t, "stream_variant", mockChunk.GetVariantName())
}

func TestDatapointInsertInterface(t *testing.T) {
	mockDatapoint := &MockDatapointInsert{FunctionName: "test_func"}
	var _ DatapointInsert = mockDatapoint // Assert that MockDatapointInsert implements DatapointInsert
	assert.Equal(t, "test_func", mockDatapoint.GetFunctionName())
}

func TestFunctionConfigInterface(t *testing.T) {
	mockConfig := &MockFunctionConfig{Type: "chat_func", Variants: VariantsConfig{}}
	var _ FunctionConfig = mockConfig // Assert that MockFunctionConfig implements FunctionConfig
	assert.Equal(t, "chat_func", mockConfig.GetType())
	assert.NotNil(t, mockConfig.GetVariants())
}

func TestVariantConfigInterface(t *testing.T) {
	mockConfig := &MockVariantConfig{Type: "chat_completion_variant"}
	var _ VariantConfig = mockConfig // Assert that MockVariantConfig implements VariantConfig
	assert.Equal(t, "chat_completion_variant", mockConfig.GetType())
}

func TestOptimizationConfigInterface(t *testing.T) {
	mockConfig := &MockOptimizationConfig{Type: "sft_optimization"}
	var _ OptimizationConfig = mockConfig // Assert that MockOptimizationConfig implements OptimizationConfig
	assert.Equal(t, "sft_optimization", mockConfig.GetType())
}

func TestOptimizationJobHandleInterface(t *testing.T) {
	mockHandle := &MockOptimizationJobHandle{Type: "job_handle", JobID: "job_123"}
	var _ OptimizationJobHandle = mockHandle // Assert that MockOptimizationJobHandle implements OptimizationJobHandle
	assert.Equal(t, "job_handle", mockHandle.GetType())
	assert.Equal(t, "job_123", mockHandle.GetJobID())
}
