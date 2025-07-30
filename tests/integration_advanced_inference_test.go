//go:build integration

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/tool"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_AdvancedInferenceScenarios(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("InferenceWithToolCalls", func(t *testing.T) {
		// Test inference that should trigger tool calls
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("weather_helper"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What's the weather like in New York City today?")},
						},
					},
				},
			},
			Tags: map[string]string{
				"test_type": "tool_calling",
				"location": "nyc",
			},
		})
		require.NoError(t, err)

		// Verify response structure
		assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
		assert.NotEqual(t, uuid.Nil, resp.GetEpisodeID())
		assert.NotEmpty(t, resp.GetVariantName())

		// For weather helper, we expect tool calls in the response
		if chatResp, ok := resp.(*inference.ChatInferenceResponse); ok {
			hasToolCall := false
			for _, content := range chatResp.Content {
				if content.GetType() == "tool_call" {
					hasToolCall = true
					break
				}
			}
			assert.True(t, hasToolCall, "Weather query should trigger tool calls")
		}
	})

	t.Run("InferenceWithParallelToolCalls", func(t *testing.T) {
		// Test parallel tool calls if supported
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("weather_helper_parallel"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What's the weather in both San Francisco and Seattle?")},
						},
					},
				},
			},
			ParallelToolCalls: util.BoolPtr(true),
			Tags: map[string]string{
				"test_type": "parallel_tools",
				"locations": "sf_seattle",
			},
		})
		require.NoError(t, err)

		// Verify response
		assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
		
		// Check for multiple tool calls if it's a chat response
		if chatResp, ok := resp.(*inference.ChatInferenceResponse); ok {
			toolCallCount := 0
			for _, content := range chatResp.Content {
				if content.GetType() == "tool_call" {
					toolCallCount++
				}
			}
			// Should have multiple tool calls for multiple locations
			assert.Greater(t, toolCallCount, 0, "Should have tool calls for weather queries")
		}
	})

	t.Run("InferenceWithCustomOutputSchema", func(t *testing.T) {
		// Test JSON inference with custom output schema
		customSchema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"entities": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{"type": "string"},
							"type": map[string]interface{}{"type": "string"},
							"confidence": map[string]interface{}{"type": "number"},
						},
						"required": []string{"name", "type"},
					},
				},
				"summary": map[string]interface{}{"type": "string"},
			},
			"required": []string{"entities"},
		}

		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("json_success"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Extract entities from: Apple Inc. was founded by Steve Jobs in Cupertino, California.")},
						},
					},
				},
			},
			OutputSchema: customSchema,
			Tags: map[string]string{
				"test_type": "custom_schema",
				"task": "entity_extraction",
			},
		})
		require.NoError(t, err)

		// Verify JSON response structure
		if jsonResp, ok := resp.(*inference.JsonInferenceResponse); ok {
			assert.NotNil(t, jsonResp.Output.Parsed, "Should have parsed JSON output")
			
			// Verify the output follows our custom schema
			if parsed := jsonResp.Output.Parsed; parsed != nil {
				if entities, exists := parsed["entities"]; exists {
					assert.IsType(t, []interface{}{}, entities, "Entities should be an array")
				}
			}
		}
	})

	t.Run("InferenceWithAdditionalTools", func(t *testing.T) {
		// Test inference with dynamically defined tools
		additionalTools := []map[string]interface{}{
			{
				"name": "calculate_tip",
				"description": "Calculate tip amount for a bill",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"bill_amount": map[string]interface{}{
							"type": "number",
							"description": "The total bill amount",
						},
						"tip_percentage": map[string]interface{}{
							"type": "number",
							"description": "The tip percentage (e.g., 15, 18, 20)",
						},
					},
					"required": []string{"bill_amount", "tip_percentage"},
				},
				"strict": true,
			},
		}

		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("I have a $50 bill and want to leave an 18% tip. How much should I tip?")},
						},
					},
				},
			},
			AdditionalTools: additionalTools,
			ToolChoice: tool.ToolChoice("auto"),
			Tags: map[string]string{
				"test_type": "dynamic_tools",
				"scenario": "tip_calculation",
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
	})

	t.Run("InferenceWithCacheOptions", func(t *testing.T) {
		// Test inference with caching enabled
		cacheOptions := map[string]interface{}{
			"enabled": "on",
			"max_age_s": 3600, // 1 hour
		}

		// First request - should miss cache
		resp1, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What is the capital of Japan?")},
						},
					},
				},
			},
			CacheOptions: cacheOptions,
			Tags: map[string]string{
				"test_type": "caching",
				"request": "first",
			},
		})
		require.NoError(t, err)
		firstInferenceID := resp1.GetInferenceID()

		// Second identical request - might hit cache
		resp2, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What is the capital of Japan?")},
						},
					},
				},
			},
			CacheOptions: cacheOptions,
			Tags: map[string]string{
				"test_type": "caching",
				"request": "second",
			},
		})
		require.NoError(t, err)
		secondInferenceID := resp2.GetInferenceID()

		// Both should be valid responses
		assert.NotEqual(t, uuid.Nil, firstInferenceID)
		assert.NotEqual(t, uuid.Nil, secondInferenceID)
		
		// Note: Whether they're the same depends on TensorZero's caching implementation
		// We just verify both requests succeeded
	})

	t.Run("InferenceWithCredentials", func(t *testing.T) {
		// Test inference with dynamic credentials (if configured)
		credentials := map[string]string{
			"openrouter_api_key": "test_key_for_dynamic_auth",
		}

		// This test might fail if dynamic credentials aren't configured
		// That's expected behavior
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Test message with dynamic credentials")},
						},
					},
				},
			},
			Credentials: credentials,
			Tags: map[string]string{
				"test_type": "dynamic_credentials",
			},
		})

		// This might succeed or fail depending on configuration
		if err != nil {
			// If it fails, it should be a TensorZero error
			var tzErr *shared.TensorZeroError
			assert.ErrorAs(t, err, &tzErr)
		} else {
			assert.NotEqual(t, uuid.Nil, resp.GetInferenceID())
		}
	})

	t.Run("InferenceWithIncludeOriginalResponse", func(t *testing.T) {
		// Test inference with original response included
		resp, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Explain quantum computing briefly")},
						},
					},
				},
			},
			IncludeOriginalResponse: util.BoolPtr(true),
			Tags: map[string]string{
				"test_type": "original_response",
				"debug": "enabled",
			},
		})
		require.NoError(t, err)

		// Verify original response is included
		originalResponse := resp.GetOriginalResponse()
		assert.NotNil(t, originalResponse, "Should include original response when requested")
		if originalResponse != nil {
			assert.NotEmpty(t, *originalResponse, "Original response should not be empty")
		}
	})

	t.Run("InferenceWithComplexConversation", func(t *testing.T) {
		// Test multi-turn conversation
		episodeID, _ := uuid.NewV7()

		// First turn
		resp1, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("I'm planning a trip to Paris. What should I know?")},
						},
					},
				},
			},
			EpisodeID: util.UUIDPtr(episodeID),
			Tags: map[string]string{
				"test_type": "conversation",
				"turn": "1",
				"topic": "travel_planning",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, episodeID, resp1.GetEpisodeID())

		// Second turn - follow up question
		resp2, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("I'm planning a trip to Paris. What should I know?")},
						},
					},
					{
						Role: "assistant",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Paris is a wonderful destination! Here are some key things to know...")},
						},
					},
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What about the best time to visit?")},
						},
					},
				},
			},
			EpisodeID: util.UUIDPtr(episodeID),
			Tags: map[string]string{
				"test_type": "conversation",
				"turn": "2",
				"topic": "travel_timing",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, episodeID, resp2.GetEpisodeID())

		// Third turn - specific question
		resp3, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("I'm planning a trip to Paris. What should I know?")},
						},
					},
					{
						Role: "assistant",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Paris is a wonderful destination! Here are some key things to know...")},
						},
					},
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What about the best time to visit?")},
						},
					},
					{
						Role: "assistant",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("The best time to visit Paris is typically late spring (April-June) or early fall (September-October)...")},
						},
					},
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Can you recommend some must-see museums?")},
						},
					},
				},
			},
			EpisodeID: util.UUIDPtr(episodeID),
			Tags: map[string]string{
				"test_type": "conversation",
				"turn": "3",
				"topic": "museums",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, episodeID, resp3.GetEpisodeID())

		// Verify all responses are part of the same episode
		assert.Equal(t, resp1.GetEpisodeID(), resp2.GetEpisodeID())
		assert.Equal(t, resp2.GetEpisodeID(), resp3.GetEpisodeID())

		// Verify all have different inference IDs
		assert.NotEqual(t, resp1.GetInferenceID(), resp2.GetInferenceID())
		assert.NotEqual(t, resp2.GetInferenceID(), resp3.GetInferenceID())
		assert.NotEqual(t, resp1.GetInferenceID(), resp3.GetInferenceID())
	})

	t.Run("InferenceErrorScenarios", func(t *testing.T) {
		// Test various error conditions
		
		// Invalid function name
		_, err := client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("nonexistent_function"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("This should fail")},
						},
					},
				},
			},
		})
		assert.Error(t, err)
		var tzErr *shared.TensorZeroError
		assert.ErrorAs(t, err, &tzErr)

		// Invalid model name
		_, err = client.Inference(ctx, &inference.InferenceRequest{
			ModelName: util.StringPtr("nonexistent_model"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("This should also fail")},
						},
					},
				},
			},
		})
		assert.Error(t, err)
		assert.ErrorAs(t, err, &tzErr)

		// Both function and model name (should be invalid)
		_, err = client.Inference(ctx, &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			ModelName:    util.StringPtr("gpt-4"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("This should fail - can't specify both")},
						},
					},
				},
			},
		})
		assert.Error(t, err)

		// Neither function nor model name
		_, err = client.Inference(ctx, &inference.InferenceRequest{
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("This should fail - need function or model")},
						},
					},
				},
			},
		})
		assert.Error(t, err)
	})

	t.Run("InferencePerformanceTest", func(t *testing.T) {
		// Test multiple concurrent inferences
		concurrency := 5
		results := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(index int) {
				_, err := client.Inference(ctx, &inference.InferenceRequest{
					FunctionName: util.StringPtr("basic_test"),
					Input: inference.InferenceInput{
						Messages: []shared.Message{
							{
								Role: "user",
								Content: []shared.ContentBlock{
									&shared.Text{Text: util.StringPtr(fmt.Sprintf("Concurrent test message %d", index))},
								},
							},
						},
					},
					Tags: map[string]string{
						"test_type": "concurrency",
						"index": fmt.Sprintf("%d", index),
					},
				})
				results <- err
			}(i)
		}

		// Wait for all to complete
		for i := 0; i < concurrency; i++ {
			select {
			case err := <-results:
				assert.NoError(t, err, "Concurrent inference %d should succeed", i)
			case <-time.After(30 * time.Second):
				t.Fatal("Concurrent inference test timed out")
			}
		}
	})
}