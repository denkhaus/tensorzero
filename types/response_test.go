package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

func TestChatInferenceResponse(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	finishReason := FinishReasonStop
	originalResponse := "original response"
	response := ChatInferenceResponse{
		InferenceID:      infID,
		EpisodeID:        epID,
		VariantName:      "default",
		Content:          []ContentBlock{NewText("Hello")},
		Usage:            Usage{InputTokens: 10, OutputTokens: 20},
		FinishReason:     &finishReason,
		OriginalResponse: &originalResponse,
	}

	assert.Equal(t, infID, response.GetInferenceID())
	assert.Equal(t, epID, response.GetEpisodeID())
	assert.Equal(t, "default", response.GetVariantName())
	assert.Len(t, response.Content, 1)
	assert.Equal(t, "Hello", *response.Content[0].(*Text).Text)
	assert.Equal(t, Usage{InputTokens: 10, OutputTokens: 20}, response.GetUsage())
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
		Usage:        Usage{InputTokens: 5, OutputTokens: 15},
		FinishReason: &finishReason,
	}

	assert.Equal(t, infID, response.GetInferenceID())
	assert.Equal(t, epID, response.GetEpisodeID())
	assert.Equal(t, "json_variant", response.GetVariantName())
	assert.Equal(t, `{"data": "value"}`, *response.Output.Raw)
	assert.Equal(t, map[string]interface{}{"data": "value"}, response.Output.Parsed)
	assert.Equal(t, Usage{InputTokens: 5, OutputTokens: 15}, response.GetUsage())
	assert.Equal(t, &finishReason, response.GetFinishReason())
	assert.Nil(t, response.GetOriginalResponse())
}

func TestChatChunk(t *testing.T) {
	infID := uuid.New()
	epID := uuid.New()
	chunk := ChatChunk{
		InferenceID: infID,
		EpisodeID:   epID,
		VariantName: "chat_stream",
		Content:     []ContentBlockChunk{&MockContentBlockChunk{Type: "text", ID: "123"}},
		Usage:       &Usage{InputTokens: 1, OutputTokens: 1},
	}

	assert.Equal(t, infID, chunk.GetInferenceID())
	assert.Equal(t, epID, chunk.GetEpisodeID())
	assert.Equal(t, "chat_stream", chunk.GetVariantName())
	assert.Len(t, chunk.Content, 1)
	assert.Equal(t, "text", chunk.Content[0].GetType())
	assert.Equal(t, "123", chunk.Content[0].GetID())
	assert.Equal(t, &Usage{InputTokens: 1, OutputTokens: 1}, chunk.Usage)
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
		Usage:        &Usage{InputTokens: 2, OutputTokens: 2},
		FinishReason: &finishReason,
	}

	assert.Equal(t, infID, chunk.GetInferenceID())
	assert.Equal(t, epID, chunk.GetEpisodeID())
	assert.Equal(t, "json_stream", chunk.GetVariantName())
	assert.Equal(t, `{"partial": "json"}`, chunk.Raw)
	assert.Equal(t, &Usage{InputTokens: 2, OutputTokens: 2}, chunk.Usage)
	assert.Equal(t, &finishReason, chunk.FinishReason)
}

func TestFeedbackResponse(t *testing.T) {
	feedbackID := uuid.New()
	response := FeedbackResponse{FeedbackID: feedbackID}
	assert.Equal(t, feedbackID, response.FeedbackID)
}

func TestDynamicEvaluationRunResponse(t *testing.T) {
	runID := uuid.New()
	response := DynamicEvaluationRunResponse{RunID: runID}
	assert.Equal(t, runID, response.RunID)
}

func TestDynamicEvaluationRunEpisodeResponse(t *testing.T) {
	episodeID := uuid.New()
	response := DynamicEvaluationRunEpisodeResponse{EpisodeID: episodeID}
	assert.Equal(t, episodeID, response.EpisodeID)
}
