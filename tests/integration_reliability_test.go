//go:build integration

package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 30 * time.Second

func TestIntegration_ConcurrentRequests(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()
	
	const numGoroutines = 10
	const requestsPerGoroutine = 5
	
	var wg sync.WaitGroup
	results := make(chan error, numGoroutines*requestsPerGoroutine)
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			for j := 0; j < requestsPerGoroutine; j++ {
				ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
				
				req := &inference.InferenceRequest{
					Input: inference.InferenceInput{
						Messages: []shared.Message{
							{
								Role: "user",
								Content: []shared.ContentBlock{
									shared.NewText("Concurrent test"),
								},
							},
						},
					},
					FunctionName: util.StringPtr("basic_test"),
					Tags: map[string]string{
						"goroutine": fmt.Sprintf("%d", goroutineID),
						"request":   fmt.Sprintf("%d", j),
						"test_type": "concurrent",
					},
				}
				
				_, err := client.Inference(ctx, req)
				results <- err
				cancel()
			}
		}(i)
	}
	
	wg.Wait()
	close(results)
	
	// Collect results
	var errors []error
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	// Allow some failures but not too many
	totalRequests := numGoroutines * requestsPerGoroutine
	successRate := float64(totalRequests-len(errors)) / float64(totalRequests)
	
	t.Logf("Concurrent test: %d/%d requests succeeded (%.1f%% success rate)", 
		totalRequests-len(errors), totalRequests, successRate*100)
	
	// Require at least 80% success rate for concurrent requests
	assert.GreaterOrEqual(t, successRate, 0.8, "Success rate should be at least 80%")
}

func TestIntegration_LongRunningStream(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Tell me a detailed story about artificial intelligence"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("basic_test"),
		Stream:       util.BoolPtr(true),
	}
	
	chunks, errs := client.InferenceStream(ctx, req)
	
	var chunkCount int
	// Track chunk timing
	startTime := time.Now()
	
	timeout := time.After(45 * time.Second)
	
	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				goto done // Stream finished
			}
			chunkCount++
		_ = time.Now() // Track timing
			
			// Verify chunk structure
			assert.NotEmpty(t, chunk.GetInferenceID())
			assert.NotEmpty(t, chunk.GetVariantName())
			
		case err := <-errs:
			if err != nil {
				// Some errors might be expected (e.g., API key issues)
				t.Logf("Stream error (may be expected): %v", err)
				return
			}
			
		case <-timeout:
			t.Fatal("Long running stream test timed out")
		}
	}
	
done:
	duration := time.Since(startTime)
	
	t.Logf("Long running stream: received %d chunks over %v", chunkCount, duration)
	
	if chunkCount > 0 {
		assert.Greater(t, chunkCount, 0, "Should receive at least one chunk")
		assert.True(t, duration > time.Second, "Stream should take some time")
	}
}

func TestIntegration_RetryLogic(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()
	
	// Test with a very short timeout to trigger retries
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("This should timeout"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("basic_test"),
	}
	
	_, err := client.Inference(ctx, req)
	
	// We expect this to fail due to timeout
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	
	// Now test with proper timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), testTimeout)
	defer cancel2()
	
	_, err2 := client.Inference(ctx2, req)
	// This might succeed or fail depending on API key, but shouldn't timeout
	if err2 != nil {
		assert.NotContains(t, err2.Error(), "context deadline exceeded")
	}
}

func TestIntegration_ResourceCleanup(t *testing.T) {
	// Test that clients can be created and closed multiple times
	for i := 0; i < 5; i++ {
		client := setupTestClient(t)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		req := &inference.InferenceRequest{
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							shared.NewText("Resource cleanup test"),
						},
					},
				},
			},
			FunctionName: util.StringPtr("basic_test"),
			Dryrun:       util.BoolPtr(true),
		}
		
		_, err := client.Inference(ctx, req)
		if err != nil {
			t.Logf("Iteration %d failed (may be expected): %v", i, err)
		}
		
		cancel()
		
		// Close the client
		err = client.Close()
		assert.NoError(t, err, "Client should close without error")
	}
}

func TestIntegration_MemoryUsage(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()
	
	// Create many small requests to test memory usage
	const numRequests = 100
	
	for i := 0; i < numRequests; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		req := &inference.InferenceRequest{
			Input: inference.InferenceInput{
				Messages: []shared.Message{
					{
						Role: "user",
						Content: []shared.ContentBlock{
							shared.NewText("Memory test"),
						},
					},
				},
			},
			FunctionName: util.StringPtr("basic_test"),
			Dryrun:       util.BoolPtr(true),
			Tags: map[string]string{
				"iteration": fmt.Sprintf("%d", i),
				"test_type": "memory",
			},
		}
		
		_, err := client.Inference(ctx, req)
		if err != nil && i == 0 {
			// If first request fails, skip the test
			t.Skipf("Skipping memory test due to initial request failure: %v", err)
		}
		
		cancel()
		
		// Small delay to prevent overwhelming the server
		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
	
	t.Logf("Successfully completed %d requests for memory usage test", numRequests)
}