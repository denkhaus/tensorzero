//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
)

// Benchmark tests for integration testing performance characteristics

func BenchmarkIntegration_BasicInference(b *testing.B) {
	client := setupClientForBenchmark(b)
	defer client.Close()
	
	ctx := context.Background()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("What is 2+2?"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("basic_test"),
		Dryrun:       util.BoolPtr(true), // Use dryrun for benchmarking
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := client.Inference(ctx, req)
		if err != nil {
			b.Fatalf("Inference failed: %v", err)
		}
	}
}

func BenchmarkIntegration_ConcurrentInferences(b *testing.B) {
	client := setupClientForBenchmark(b)
	defer client.Close()
	
	ctx := context.Background()
	
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
		Dryrun:       util.BoolPtr(true),
	}
	
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Inference(ctx, req)
			if err != nil {
				b.Errorf("Concurrent inference failed: %v", err)
			}
		}
	})
}

func BenchmarkIntegration_StreamingInference(b *testing.B) {
	client := setupClientForBenchmark(b)
	defer client.Close()
	
	ctx := context.Background()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Stream test"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("basic_test"),
		Stream:       util.BoolPtr(true),
		Dryrun:       util.BoolPtr(true),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		chunks, errs := client.InferenceStream(ctx, req)
		
		// Consume the stream
		for {
			select {
			case _, ok := <-chunks:
				if !ok {
					goto next
				}
			case err := <-errs:
				if err != nil {
					b.Fatalf("Streaming failed: %v", err)
				}
				goto next
			case <-time.After(5 * time.Second):
				b.Fatal("Streaming timed out")
			}
		}
		
	next:
	}
}

func BenchmarkIntegration_JSONInference(b *testing.B) {
	client := setupClientForBenchmark(b)
	defer client.Close()
	
	ctx := context.Background()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Extract JSON: John is 25 years old"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("json_with_schemas"),
		Dryrun:       util.BoolPtr(true),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := client.Inference(ctx, req)
		if err != nil {
			b.Fatalf("JSON inference failed: %v", err)
		}
	}
}

// setupClientForBenchmark creates a client for benchmarking with appropriate timeouts
func setupClientForBenchmark(b *testing.B) tensorzero.Gateway {
	client := setupTestClient(b)
	
	// Quick health check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	req := &inference.InferenceRequest{
		Input: inference.InferenceInput{
			Messages: []shared.Message{
				{
					Role: "user",
					Content: []shared.ContentBlock{
						shared.NewText("Health check"),
					},
				},
			},
		},
		FunctionName: util.StringPtr("basic_test"),
		Dryrun:       util.BoolPtr(true),
	}
	
	_, err := client.Inference(ctx, req)
	if err != nil {
		b.Skipf("TensorZero gateway not available for benchmarking: %v", err)
	}
	
	return client
}