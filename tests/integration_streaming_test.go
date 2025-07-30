//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_StreamingInference(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("BasicChatStreaming", func(t *testing.T) {
		req := &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Tell me a short story about a robot")},
						},
					},
				},
			},
			Stream: util.BoolPtr(true),
			Tags: map[string]string{
				"test_type": "streaming",
				"story_type": "robot",
			},
		}

		chunks, errs := client.InferenceStream(ctx, req)
		
		var receivedChunks []inference.InferenceChunk
		var finalChunk inference.InferenceChunk
		timeout := time.After(30 * time.Second)

		// Collect all chunks
		for {
			select {
			case chunk, ok := <-chunks:
				if !ok {
					// Channel closed, streaming finished
					goto StreamComplete
				}
				receivedChunks = append(receivedChunks, chunk)
				finalChunk = chunk
				
				// Verify chunk structure
				assert.NotEqual(t, uuid.Nil, chunk.GetInferenceID())
				assert.NotEqual(t, uuid.Nil, chunk.GetEpisodeID())
				assert.NotEmpty(t, chunk.GetVariantName())

			case err := <-errs:
				if err != nil {
					t.Fatalf("Streaming error: %v", err)
				}

			case <-timeout:
				t.Fatal("Streaming test timed out")
			}
		}

	StreamComplete:
		// Verify we received chunks
		assert.Greater(t, len(receivedChunks), 0, "Should receive at least one chunk")

		// Verify final chunk has usage information
		if chatChunk, ok := finalChunk.(*inference.ChatChunk); ok {
			if chatChunk.Usage != nil {
				assert.Greater(t, chatChunk.Usage.InputTokens, 0)
				assert.Greater(t, chatChunk.Usage.OutputTokens, 0)
			}
			if chatChunk.FinishReason != nil {
				assert.Contains(t, []inference.FinishReason{
					inference.FinishReasonStop,
					inference.FinishReasonLength,
				}, *chatChunk.FinishReason)
			}
		}
	})

	t.Run("JsonStreaming", func(t *testing.T) {
		req := &inference.InferenceRequest{
			FunctionName: util.StringPtr("json_success"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Extract entities from: Alice works at Microsoft in Seattle")},
						},
					},
				},
			},
			Stream: util.BoolPtr(true),
			Tags: map[string]string{
				"test_type": "json_streaming",
				"task": "entity_extraction",
			},
		}

		chunks, errs := client.InferenceStream(ctx, req)
		
		var receivedChunks []inference.InferenceChunk
		var finalChunk inference.InferenceChunk
		timeout := time.After(30 * time.Second)

		// Collect all chunks
		for {
			select {
			case chunk, ok := <-chunks:
				if !ok {
					goto JsonStreamComplete
				}
				receivedChunks = append(receivedChunks, chunk)
				finalChunk = chunk

			case err := <-errs:
				if err != nil {
					t.Fatalf("JSON streaming error: %v", err)
				}

			case <-timeout:
				t.Fatal("JSON streaming test timed out")
			}
		}

	JsonStreamComplete:
		// Verify we received chunks
		assert.Greater(t, len(receivedChunks), 0, "Should receive at least one JSON chunk")

		// Verify final chunk structure for JSON
		if jsonChunk, ok := finalChunk.(*inference.JsonChunk); ok {
			assert.NotEmpty(t, jsonChunk.Raw, "Should have raw JSON content")
			if jsonChunk.Usage != nil {
				assert.Greater(t, jsonChunk.Usage.InputTokens, 0)
				assert.Greater(t, jsonChunk.Usage.OutputTokens, 0)
			}
		}
	})

	t.Run("StreamingWithToolCalls", func(t *testing.T) {
		req := &inference.InferenceRequest{
			FunctionName: util.StringPtr("weather_helper"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("What's the weather like in San Francisco?")},
						},
					},
				},
			},
			Stream: util.BoolPtr(true),
			Tags: map[string]string{
				"test_type": "tool_streaming",
				"location": "san_francisco",
			},
		}

		chunks, errs := client.InferenceStream(ctx, req)
		
		var receivedChunks []inference.InferenceChunk
		var hasToolCall bool
		timeout := time.After(45 * time.Second) // Longer timeout for tool calls

		// Collect all chunks
		for {
			select {
			case chunk, ok := <-chunks:
				if !ok {
					goto ToolStreamComplete
				}
				receivedChunks = append(receivedChunks, chunk)

				// Check if this chunk contains tool calls
				if chatChunk, ok := chunk.(*inference.ChatChunk); ok {
					for _, content := range chatChunk.Content {
						if content.GetType() == "tool_call" {
							hasToolCall = true
						}
					}
				}

			case err := <-errs:
				if err != nil {
					t.Fatalf("Tool streaming error: %v", err)
				}

			case <-timeout:
				t.Fatal("Tool streaming test timed out")
			}
		}

	ToolStreamComplete:
		// Verify we received chunks
		assert.Greater(t, len(receivedChunks), 0, "Should receive at least one chunk")
		
		// For weather helper, we expect tool calls
		assert.True(t, hasToolCall, "Should have received tool call chunks for weather query")
	})

	t.Run("StreamingErrorHandling", func(t *testing.T) {
		// Test streaming with invalid function
		req := &inference.InferenceRequest{
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
			Stream: util.BoolPtr(true),
		}

		chunks, errs := client.InferenceStream(ctx, req)
		
		timeout := time.After(10 * time.Second)
		var receivedError error

		// Should receive an error
		select {
		case chunk := <-chunks:
			if chunk != nil {
				t.Fatal("Should not receive chunks for invalid function")
			}

		case err := <-errs:
			receivedError = err

		case <-timeout:
			t.Fatal("Should receive error quickly for invalid function")
		}

		// Verify we got an error
		require.Error(t, receivedError)
		var tzErr *shared.TensorZeroError
		assert.ErrorAs(t, receivedError, &tzErr)
	})

	t.Run("StreamingContextCancellation", func(t *testing.T) {
		// Create a context that we'll cancel
		streamCtx, cancel := context.WithCancel(ctx)

		req := &inference.InferenceRequest{
			FunctionName: util.StringPtr("basic_test"),
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							&shared.Text{Text: util.StringPtr("Tell me a very long story that takes time to generate")},
						},
					},
				},
			},
			Stream: util.BoolPtr(true),
		}

		chunks, errs := client.InferenceStream(streamCtx, req)
		
		// Cancel the context after a short delay
		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		timeout := time.After(10 * time.Second)
		var receivedError error

		// Should receive cancellation error
		for {
			select {
			case chunk := <-chunks:
				if chunk == nil {
					// Channel closed due to cancellation
					goto CancelComplete
				}
				// Continue receiving chunks until cancellation

			case err := <-errs:
				receivedError = err
				goto CancelComplete

			case <-timeout:
				t.Fatal("Context cancellation test timed out")
			}
		}

	CancelComplete:
		// Should have received a context cancellation error
		if receivedError != nil {
			assert.Contains(t, receivedError.Error(), "context canceled")
		}
	})
}