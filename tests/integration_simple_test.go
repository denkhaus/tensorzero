//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_SimpleInference(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("BasicChatInference", func(t *testing.T) {
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				System: map[string]interface{}{
					"assistant_name": "TestBot",
				},
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What is 2+2?"), Type: "text"},
						},
					},
				},
			},
			Tags: map[string]string{
				"test_type": "simple_integration",
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
		assert.NotEqual(t, uuid.Nil, resp.GetEpisodeID())
		assert.NotEmpty(t, resp.GetVariantName())
	})

	t.Run("JSONInference", func(t *testing.T) {
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("json_success"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Extract JSON data: John is 25 years old"), Type: "text"},
						},
					},
				},
			},
			Tags: map[string]string{
				"test_type": "json_integration",
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
		
		// Verify it's a JSON response
		if jsonResp, ok := resp.(*inference.JsonInferenceResponse); ok {
			assert.NotNil(t, jsonResp.Output)
		}
	})
}