package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

func TestChatDatapointInsert(t *testing.T) {
	input := InferenceInput{
		Messages: []Message{
			{Role: "user", Content: []ContentBlock{NewText("Test message")}},
		},
	}
	datapoint := ChatDatapointInsert{
		FunctionName: "test_chat_function",
		Input:        input,
		Output:       "response",
		AllowedTools: []string{"tool1"},
		Tags:         map[string]string{"env": "dev"},
	}

	assert.Equal(t, "test_chat_function", datapoint.GetFunctionName())
	assert.Equal(t, input, datapoint.Input)
	assert.Equal(t, "response", datapoint.Output)
	assert.Contains(t, datapoint.AllowedTools, "tool1")
	assert.Equal(t, "dev", datapoint.Tags["env"])

	jsonBytes, err := json.Marshal(datapoint)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"function_name":"test_chat_function"`)
}

func TestJsonDatapointInsert(t *testing.T) {
	input := InferenceInput{
		Messages: []Message{
			{Role: "user", Content: []ContentBlock{NewText("Test JSON message")}},
		},
	}
	outputSchema := map[string]interface{}{"type": "object"}
	datapoint := JsonDatapointInsert{
		FunctionName: "test_json_function",
		Input:        input,
		Output:       map[string]interface{}{"key": "value"},
		OutputSchema: outputSchema,
		Tags:         map[string]string{"type": "json"},
	}

	assert.Equal(t, "test_json_function", datapoint.GetFunctionName())
	assert.Equal(t, input, datapoint.Input)
	assert.Equal(t, map[string]interface{}{"key": "value"}, datapoint.Output)
	assert.Equal(t, outputSchema, datapoint.OutputSchema)
	assert.Equal(t, "json", datapoint.Tags["type"])

	jsonBytes, err := json.Marshal(datapoint)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"function_name":"test_json_function"`)
}

func TestListInferencesRequest(t *testing.T) {
	funcName := "my_func"
	episodeID := uuid.New()
	variantName := "my_variant"
	limit := 10
	offset := 0
	filter := NewFloatMetricFilter("score", 0.8, ">")
	orderBy := NewOrderByTimestamp("DESC")

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

func TestInferenceRequest(t *testing.T) {
	input := InferenceInput{Messages: []Message{{Role: "user", Content: []ContentBlock{NewText("Hello")}}}}
	funcName := "my_func"
	modelName := "gpt-4"
	episodeID := uuid.New()
	stream := true
	params := map[string]interface{}{"temp": 0.5}
	variantName := "default"
	dryrun := false
	outputSchema := map[string]interface{}{"type": "string"}
	allowedTools := []string{"tool1"}
	additionalTools := []map[string]interface{}{{"name": "tool2"}}
	toolChoice := ToolChoice("auto")
	parallelToolCalls := true
	internal := false
	tags := map[string]string{"source": "test"}
	credentials := map[string]string{"key": "value"}
	cacheOptions := map[string]interface{}{"ttl": 300}
	extraBody := []ExtraBody{}
	extraHeaders := []map[string]interface{}{{"X-Custom": "header"}}
	includeOriginalResponse := true

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

func TestStoredInference(t *testing.T) {
	infID := uuid.New()
	episodeID := uuid.New()
	input := InferenceInput{Messages: []Message{{Role: "user", Content: []ContentBlock{NewText("Question")}}}}
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

func TestDatapoint(t *testing.T) {
	id := uuid.New()
	input := InferenceInput{Messages: []Message{{Role: "user", Content: []ContentBlock{NewText("Data")}}}}
	output := "Processed Data"
	outputSchema := map[string]interface{}{"format": "json"}

	dp := Datapoint{
		ID:           id,
		Input:        input,
		Output:       output,
		DatasetName:  "my_dataset",
		FunctionName: "data_processor",
		OutputSchema: outputSchema,
		IsCustom:     true,
	}

	assert.Equal(t, id, dp.ID)
	assert.Equal(t, input, dp.Input)
	assert.Equal(t, output, dp.Output)
	assert.Equal(t, "my_dataset", dp.DatasetName)
	assert.Equal(t, "data_processor", dp.FunctionName)
	assert.Equal(t, outputSchema, dp.OutputSchema)
	assert.True(t, dp.IsCustom)
}

func TestFeedbackRequest(t *testing.T) {
	inferenceID := uuid.New()
	episodeID := uuid.New()
	dryrun := true
	internal := false
	tags := map[string]string{"feedback_type": "quality"}

	req := FeedbackRequest{
		MetricName:  "relevance",
		Value:       0.9,
		InferenceID: &inferenceID,
		EpisodeID:   &episodeID,
		Dryrun:      &dryrun,
		Internal:    &internal,
		Tags:        tags,
	}

	assert.Equal(t, "relevance", req.MetricName)
	assert.Equal(t, 0.9, req.Value)
	assert.Equal(t, inferenceID, *req.InferenceID)
	assert.Equal(t, episodeID, *req.EpisodeID)
	assert.True(t, *req.Dryrun)
	assert.False(t, *req.Internal)
	assert.Equal(t, tags, req.Tags)
}

func TestDynamicEvaluationRunRequest(t *testing.T) {
	variants := map[string]string{"A": "variant1", "B": "variant2"}
	tags := map[string]string{"eval_group": "groupA"}
	projectName := "project_X"
	displayName := "Eval Run 1"

	req := DynamicEvaluationRunRequest{
		Variants:    variants,
		Tags:        tags,
		ProjectName: &projectName,
		DisplayName: &displayName,
	}

	assert.Equal(t, variants, req.Variants)
	assert.Equal(t, tags, req.Tags)
	assert.Equal(t, projectName, *req.ProjectName)
	assert.Equal(t, displayName, *req.DisplayName)
}

func TestDynamicEvaluationRunEpisodeRequest(t *testing.T) {
	runID := uuid.New()
	taskName := "task_Y"
	datapointName := "dp_Z"
	tags := map[string]string{"type": "episode"}

	req := DynamicEvaluationRunEpisodeRequest{
		RunID:         runID,
		TaskName:      &taskName,
		DatapointName: &datapointName,
		Tags:          tags,
	}

	assert.Equal(t, runID, req.RunID)
	assert.Equal(t, taskName, *req.TaskName)
	assert.Equal(t, datapointName, *req.DatapointName)
	assert.Equal(t, tags, req.Tags)
}

func TestListDatapointsRequest(t *testing.T) {
	funcName := "filter_func"
	limit := 5
	offset := 0

	req := ListDatapointsRequest{
		DatasetName:  "my_dataset",
		FunctionName: &funcName,
		Limit:        &limit,
		Offset:       &offset,
	}

	assert.Equal(t, "my_dataset", req.DatasetName)
	assert.Equal(t, funcName, *req.FunctionName)
	assert.Equal(t, limit, *req.Limit)
	assert.Equal(t, offset, *req.Offset)
}
