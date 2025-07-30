//go:build unit

package inference

import (
	"encoding/json"
	"testing"

	"github.com/denkhaus/tensorzero/filter"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/tool"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
func TestContentBlockChunkInterface(t *testing.T) {
	mockChunk := &MockContentBlockChunk{Type: "chunk_type", ID: "chunk_id"}
	var _ ContentBlockChunk = mockChunk // Assert that MockContentBlockChunk implements ContentBlockChunk
	assert.Equal(t, "chunk_type", mockChunk.GetType())
	assert.Equal(t, "chunk_id", mockChunk.GetID())
}

func TestInferenceRequest(t *testing.T) {
	input := InferenceInput{
		Messages: []shared.Message{
			{
				Role: "user",
				Content: []shared.ContentBlock{
					shared.NewText("Hello"),
				},
			},
		},
	}

	funcName := "my_func"
	modelName := "gpt-4"

	stream := true
	params := map[string]interface{}{"temp": 0.5}
	variantName := "default"
	dryrun := false
	outputSchema := map[string]interface{}{"type": "string"}
	allowedTools := []string{"tool1"}
	additionalTools := []map[string]interface{}{{"name": "tool2"}}
	toolChoice := tool.ToolChoice("auto")
	parallelToolCalls := true
	internal := false
	tags := map[string]string{"source": "test"}
	credentials := map[string]string{"key": "value"}
	cacheOptions := map[string]interface{}{"ttl": 300}
	extraBody := []ExtraBody{}
	extraHeaders := []map[string]interface{}{{"X-Custom": "header"}}
	includeOriginalResponse := true

	episodeID, err := uuid.NewV7()
	assert.NoError(t, err)

	req := InferenceRequest{
		Input:                   input,
		FunctionName:            &funcName,
		ModelName:               &modelName,
		EpisodeID:               &episodeID,
		Stream:                  &stream,
		Params:                  params,
		VariantName:             &variantName,
		Dryrun:                  &dryrun,
		OutputSchema:            outputSchema,
		AllowedTools:            allowedTools,
		AdditionalTools:         additionalTools,
		ToolChoice:              toolChoice,
		ParallelToolCalls:       &parallelToolCalls,
		Internal:                &internal,
		Tags:                    tags,
		Credentials:             credentials,
		CacheOptions:            cacheOptions,
		ExtraBody:               extraBody,
		ExtraHeaders:            extraHeaders,
		IncludeOriginalResponse: &includeOriginalResponse,
	}

	assert.Equal(t, input, req.Input)
	assert.Equal(t, funcName, *req.FunctionName)
	assert.Equal(t, modelName, *req.ModelName)
	assert.Equal(t, episodeID, *req.EpisodeID)
	assert.Equal(t, stream, *req.Stream)
	assert.Equal(t, params, req.Params)
	assert.Equal(t, variantName, *req.VariantName)
	assert.Equal(t, dryrun, *req.Dryrun)
	assert.Equal(t, outputSchema, req.OutputSchema)
	assert.Equal(t, allowedTools, req.AllowedTools)
	assert.Equal(t, additionalTools, req.AdditionalTools)
	assert.Equal(t, toolChoice, req.ToolChoice)
	assert.Equal(t, parallelToolCalls, *req.ParallelToolCalls)
	assert.Equal(t, internal, *req.Internal)
	assert.Equal(t, tags, req.Tags)
	assert.Equal(t, credentials, req.Credentials)
	assert.Equal(t, cacheOptions, req.CacheOptions)
	assert.Equal(t, extraBody, req.ExtraBody)
	assert.Equal(t, extraHeaders, req.ExtraHeaders)
	assert.Equal(t, includeOriginalResponse, *req.IncludeOriginalResponse)
}

func TestInferenceResponseInterface(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonStop
	originalResponse := "original"
	mockResponse := &MockInferenceResponse{
		InferenceID:         infID,
		EpisodeID:           epID,
		VariantName:         "default",
		UsageVal:            shared.Usage{InputTokens: 10, OutputTokens: 20},
		FinishReasonVal:     &finishReason,
		OriginalResponseVal: &originalResponse,
	}
	var _ InferenceResponse = mockResponse // Assert that MockInferenceResponse implements InferenceResponse
	assert.Equal(t, infID, mockResponse.GetInferenceID())
	assert.Equal(t, epID, mockResponse.GetEpisodeID())
	assert.Equal(t, "default", mockResponse.GetVariantName())
	assert.Equal(t, shared.Usage{InputTokens: 10, OutputTokens: 20}, mockResponse.GetUsage())
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

func TestChatInferenceResponse(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonStop
	originalResponse := "original response"
	response := ChatInferenceResponse{
		InferenceID:      infID,
		EpisodeID:        epID,
		VariantName:      "default",
		Content:          []shared.ContentBlock{shared.NewText("Hello")},
		Usage:            shared.Usage{InputTokens: 10, OutputTokens: 20},
		FinishReason:     &finishReason,
		OriginalResponse: &originalResponse,
	}

	assert.Equal(t, infID, response.GetInferenceID())
	assert.Equal(t, epID, response.GetEpisodeID())
	assert.Equal(t, "default", response.GetVariantName())
	assert.Len(t, response.Content, 1)
	assert.Equal(t, "Hello", *response.Content[0].(*shared.Text).Text)
	assert.Equal(t, shared.Usage{InputTokens: 10, OutputTokens: 20}, response.GetUsage())
	assert.Equal(t, &finishReason, response.GetFinishReason())
	assert.Equal(t, &originalResponse, response.GetOriginalResponse())
}

func TestJsonInferenceResponse(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonLength
	response := JsonInferenceResponse{
		InferenceID:  infID,
		EpisodeID:    epID,
		VariantName:  "json_variant",
		Output:       JsonInferenceOutput{Raw: func() *string { s := `{"data": "value"}`; return &s }(), Parsed: map[string]interface{}{"data": "value"}},
		Usage:        shared.Usage{InputTokens: 5, OutputTokens: 15},
		FinishReason: &finishReason,
	}

	assert.Equal(t, infID, response.GetInferenceID())
	assert.Equal(t, epID, response.GetEpisodeID())
	assert.Equal(t, "json_variant", response.GetVariantName())
	assert.Equal(t, `{"data": "value"}`, *response.Output.Raw)
	assert.Equal(t, map[string]interface{}{"data": "value"}, response.Output.Parsed)
	assert.Equal(t, shared.Usage{InputTokens: 5, OutputTokens: 15}, response.GetUsage())
	assert.Equal(t, &finishReason, response.GetFinishReason())
	assert.Nil(t, response.GetOriginalResponse())
}

func TestJsonInferenceOutput(t *testing.T) {
	outputJSON := `{
		"raw": "{\"name\": \"test\"}",
		"parsed": {"name": "test"}
	}`
	var output JsonInferenceOutput
	err := json.Unmarshal([]byte(outputJSON), &output)
	assert.NoError(t, err)
	assert.NotNil(t, output.Raw)
	assert.Equal(t, "{\"name\": \"test\"}", *output.Raw)
	assert.NotNil(t, output.Parsed)
	assert.Equal(t, map[string]interface{}{"name": "test"}, output.Parsed)
}

func TestUsage(t *testing.T) {
	usage := shared.Usage{InputTokens: 10, OutputTokens: 20}
	assert.Equal(t, 10, usage.InputTokens)
	assert.Equal(t, 20, usage.OutputTokens)

	jsonBytes, err := json.Marshal(usage)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"input_tokens":10`)
	assert.Contains(t, string(jsonBytes), `"output_tokens":20`)
}

func TestFinishReason(t *testing.T) {
	assert.Equal(t, "stop", string(FinishReasonStop))
	assert.Equal(t, "length", string(FinishReasonLength))
	assert.Equal(t, "tool_call", string(FinishReasonToolCall))
	assert.Equal(t, "content_filter", string(FinishReasonContentFilter))
	assert.Equal(t, "unknown", string(FinishReasonUnknown))
}

func TestListInferencesRequest(t *testing.T) {
	funcName := "my_func"
	episodeID := uuid.New()
	variantName := "my_variant"
	limit := 10
	offset := 0
	filter := filter.NewFloatMetricFilter("score", 0.8, ">")
	orderBy := shared.NewOrderByTimestamp("DESC")

	req := ListInferencesRequest{
		FunctionName: &funcName,
		EpisodeID:    &episodeID,
		VariantName:  &variantName,
		Filter:       filter,
		OrderBy:      orderBy,
		Limit:        &limit,
		Offset:       &offset,
	}

	assert.Equal(t, funcName, *req.FunctionName)
	assert.Equal(t, episodeID, *req.EpisodeID)
	assert.Equal(t, variantName, *req.VariantName)
	assert.Equal(t, filter, req.Filter)
	assert.Equal(t, orderBy, req.OrderBy)
	assert.Equal(t, limit, *req.Limit)
	assert.Equal(t, offset, *req.Offset)
}

func TestStoredInference(t *testing.T) {
	infID := uuid.New()
	episodeID := uuid.New()
	input := InferenceInput{
		Messages: []shared.Message{
			{
				Role: "user",
				Content: []shared.ContentBlock{
					shared.NewText("Question"),
				},
			},
		},
	}
	output := "Answer"
	timestamp := "2023-01-01T12:00:00Z"
	tags := map[string]string{"stage": "prod"}
	metricValues := map[string]interface{}{"accuracy": 0.9}

	inf := StoredInference{
		ID:           infID,
		EpisodeID:    episodeID,
		FunctionName: "test_func",
		VariantName:  "test_variant",
		Input:        input,
		Output:       output,
		Timestamp:    timestamp,
		Tags:         tags,
		MetricValues: metricValues,
	}

	assert.Equal(t, infID, inf.ID)
	assert.Equal(t, episodeID, inf.EpisodeID)
	assert.Equal(t, "test_func", inf.FunctionName)
	assert.Equal(t, "test_variant", inf.VariantName)
	assert.Equal(t, input, inf.Input)
	assert.Equal(t, output, inf.Output)
	assert.Equal(t, timestamp, inf.Timestamp)
	assert.Equal(t, tags, inf.Tags)
	assert.Equal(t, metricValues, inf.MetricValues)
}

func TestInferenceInput(t *testing.T) {
	inputJSON := `{
		"messages": [
			{
				"role": "user",
				"content": [
					{"type": "text", "text": "Hello"}
				]
			}
		],
		"system": {}
	}`

	var input InferenceInput
	err := json.Unmarshal([]byte(inputJSON), &input)
	assert.NoError(t, err)
	assert.Len(t, input.Messages, 1)
	assert.Equal(t, "user", input.Messages[0].Role)
	assert.NotNil(t, input.System)
}

func TestChatChunk(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	chunk := ChatChunk{
		InferenceID: infID,
		EpisodeID:   epID,
		VariantName: "chat_stream",
		Content:     []ContentBlockChunk{&MockContentBlockChunk{Type: "text", ID: "123"}},
		Usage:       &shared.Usage{InputTokens: 1, OutputTokens: 1},
	}

	assert.Equal(t, infID, chunk.GetInferenceID())
	assert.Equal(t, epID, chunk.GetEpisodeID())
	assert.Equal(t, "chat_stream", chunk.GetVariantName())
	assert.Len(t, chunk.Content, 1)
	assert.Equal(t, "text", chunk.Content[0].GetType())
	assert.Equal(t, "123", chunk.Content[0].GetID())
	assert.Equal(t, &shared.Usage{InputTokens: 1, OutputTokens: 1}, chunk.Usage)
	assert.Nil(t, chunk.FinishReason)
}

func TestJsonChunk(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonStop
	chunk := JsonChunk{
		InferenceID:  infID,
		EpisodeID:    epID,
		VariantName:  "json_stream",
		Raw:          `{"partial": "json"}`,
		Usage:        &shared.Usage{InputTokens: 2, OutputTokens: 2},
		FinishReason: &finishReason,
	}

	assert.Equal(t, infID, chunk.GetInferenceID())
	assert.Equal(t, epID, chunk.GetEpisodeID())
	assert.Equal(t, "json_stream", chunk.GetVariantName())
	assert.Equal(t, `{"partial": "json"}`, chunk.Raw)
	assert.Equal(t, &shared.Usage{InputTokens: 2, OutputTokens: 2}, chunk.Usage)
	assert.Equal(t, &finishReason, chunk.FinishReason)
}
