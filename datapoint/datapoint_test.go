//go:build unit

package datapoint

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/google/uuid"
)

func TestDatapoint(t *testing.T) {
	id := uuid.New()
	input := inference.InferenceInput{
		Messages: []shared.Message{
			{
				Role: "user",
				Content: []shared.ContentBlock{
					shared.NewText("Data"),
				},
			},
		},
	}

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
	assert.Equal(t, interface{}(output), dp.Output)
	assert.Equal(t, "my_dataset", dp.DatasetName)
	assert.Equal(t, "data_processor", dp.FunctionName)
	assert.Equal(t, interface{}(outputSchema), dp.OutputSchema)
	assert.True(t, dp.IsCustom)
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

func TestJsonDatapointInsert(t *testing.T) {
	input := inference.InferenceInput{
		Messages: []shared.Message{
			{Role: "user", Content: []shared.ContentBlock{shared.NewText("Test JSON message")}},
		},
	}
	outputSchema := map[string]interface{}{"type": "object"}
	dpoint := inference.JsonDatapointInsert{
		FunctionName: "test_json_function",
		Input:        input,
		Output:       map[string]interface{}{"key": "value"},
		OutputSchema: outputSchema,
		Tags:         map[string]string{"type": "json"},
	}

	assert.Equal(t, "test_json_function", dpoint.GetFunctionName())
	assert.Equal(t, input, dpoint.Input)
	assert.Equal(t, interface{}(map[string]interface{}{"key": "value"}), dpoint.Output)
	assert.Equal(t, interface{}(outputSchema), dpoint.OutputSchema)
	assert.Equal(t, "json", dpoint.Tags["type"])

	jsonBytes, err := json.Marshal(dpoint)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"function_name":"test_json_function"`)
}

func TestChatDatapointInsert(t *testing.T) {
	input := inference.InferenceInput{
		Messages: []shared.Message{
			{Role: "user", Content: []shared.ContentBlock{shared.NewText("Test message")}},
		},
	}
	dp := inference.ChatDatapointInsert{
		FunctionName: "test_chat_function",
		Input:        input,
		Output:       "response",
		AllowedTools: []string{"tool1"},
		Tags:         map[string]string{"env": "dev"},
	}

	assert.Equal(t, "test_chat_function", dp.GetFunctionName())
	assert.Equal(t, input, dp.Input)
	assert.Equal(t, "response", dp.Output)
	assert.Equal(t, "tool1", dp.AllowedTools[0])
	assert.Equal(t, "dev", dp.Tags["env"])

	jsonBytes, err := json.Marshal(dp)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), `"function_name":"test_chat_function"`)
}
