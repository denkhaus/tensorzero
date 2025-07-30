//go:build unit

package tensorzero

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPGateway(t *testing.T) {
	baseURL := "http://localhost:3000"
	client := NewHTTPGateway(baseURL)
	
	assert.NotNil(t, client)
	// Test that we can call methods on the client
	err := client.Close()
	assert.NoError(t, err)
}

func TestWithTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	client := NewHTTPGateway("http://localhost:3000", WithTimeout(customTimeout))
	
	assert.NotNil(t, client)
	// Test that the client was created with custom timeout
	err := client.Close()
	assert.NoError(t, err)
}

func TestWithHTTPClient(t *testing.T) {
	customHTTPClient := &http.Client{Timeout: 45 * time.Second}
	client := NewHTTPGateway("http://localhost:3000", WithHTTPClient(customHTTPClient))
	
	assert.NotNil(t, client)
	// Test that the client was created with custom HTTP client
	err := client.Close()
	assert.NoError(t, err)
}

func TestHTTPGatewayClose(t *testing.T) {
	client := NewHTTPGateway("http://localhost:3000")
	err := client.Close()
	assert.NoError(t, err)
}

func TestInferenceWithMockServer(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/inference", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		
		// Return a mock chat response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"inference_id": "550e8400-e29b-41d4-a716-446655440000",
			"episode_id": "550e8400-e29b-41d4-a716-446655440001",
			"variant_name": "test_variant",
			"content": [{"type": "text", "text": "Hello, world!"}],
			"usage": {"input_tokens": 10, "output_tokens": 5},
			"finish_reason": "stop"
		}`))
	}))
	defer server.Close()
	
	client := NewHTTPGateway(server.URL)
	
	request := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Hello"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("test_function"),
	}
	
	response, err := client.Inference(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, response)
	
	chatResp, ok := response.(*inference.ChatInferenceResponse)
	require.True(t, ok)
	assert.Equal(t, "test_variant", chatResp.VariantName)
	if chatResp.FinishReason != nil {
		assert.Equal(t, "stop", string(*chatResp.FinishReason))
	}
}

func TestInferenceWithError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request"}`))
	}))
	defer server.Close()
	
	client := NewHTTPGateway(server.URL)
	
	request := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Hello"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("test_function"),
	}
	
	response, err := client.Inference(context.Background(), request)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestInferenceStreamChannels(t *testing.T) {
	// Test that InferenceStream returns proper channels
	client := NewHTTPGateway("http://localhost:3000")
	
	request := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Hello"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("test_function"),
		Stream:       util.BoolPtr(true),
	}
	
	chunks, errs := client.InferenceStream(context.Background(), request)
	
	// Verify channels are created
	assert.NotNil(t, chunks)
	assert.NotNil(t, errs)
}

func TestInferenceWithContextCancellation(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay to allow cancellation
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"content": [{"type": "text", "text": "Response"}]}`))
	}))
	defer server.Close()
	
	client := NewHTTPGateway(server.URL)
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	request := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Hello"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("test_function"),
	}
	
	response, err := client.Inference(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestInferenceRequestValidation(t *testing.T) {
	client := NewHTTPGateway("http://localhost:3000")
	
	// Test with nil request
	response, err := client.Inference(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, response)
	
	// Test with empty messages
	request := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{},
		},
		FunctionName: util.StringPtr("test_function"),
	}
	
	response, err = client.Inference(context.Background(), request)
	assert.Error(t, err)
	assert.Nil(t, response)
}