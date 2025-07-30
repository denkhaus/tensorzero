# TensorZero Go Client

[![GitHub issues](https://img.shields.io/github/issues/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/pulls)
[![GitHub contributors](https://img.shields.io/github/contributors/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/graphs/contributors)

**Disclaimer:** This is an unofficial, community-maintained Go client for the TensorZero API.

A comprehensive Go client library for [TensorZero](https://github.com/tensorzero/tensorzero), an AI inference gateway that provides a unified interface for multiple AI model providers with features like A/B testing, optimization, and observability.

This Go SDK is a complete and accurate port of the Python SDK, providing full feature parity for the Go ecosystem. For a detailed breakdown of the feature parity, see [PYTHON_SDK_PARITY.md](./docs/PYTHON_SDK_PARITY.md).

## Features

*   **Complete TensorZero API Client:** Implements all TensorZero API endpoints, including inference, streaming, feedback, datapoint management, and dynamic evaluation.
*   **OpenAI SDK Compatibility:** Designed to be compatible with the OpenAI SDK for easy migration.
*   **Go-Idiomatic Design:** Leverages Go's best practices, including context support, channel-based streaming, structured error handling, and comprehensive documentation.
*   **Well-Organized Package Structure:** Logically grouped functionality with dedicated packages for `inference`, `feedback`, `evaluation`, `datapoint`, `tool`, and more.
*   **Comprehensive Documentation:** Extensively documented structs and types with detailed field descriptions, usage examples, and best practices.
*   **Production-Ready Testing:** Comprehensive test suite with unit tests, integration tests, performance benchmarks, and reliability testing. See [INTEGRATION_TESTS.md](./docs/INTEGRATION_TESTS.md) for details.
*   **Automated Test Environment:** Complete Docker-based development and testing environment with automated setup scripts.
*   **Enterprise-Grade Reliability:** Concurrent request handling, proper resource management, graceful error handling, and robust streaming support.

## Package Structure

The TensorZero Go client is organized into logical packages for better maintainability and ease of use:

- **`inference`** - Core inference requests, responses, and streaming functionality
- **`feedback`** - Feedback submission for metrics and model improvement
- **`evaluation`** - Dynamic evaluation runs and episode management
- **`datapoint`** - Dataset management and datapoint operations
- **`tool`** - Tool definitions and parameters for model interactions
- **`config`** - Configuration types and validation
- **`filter`** - Advanced filtering capabilities for queries
- **`shared`** - Common types and utilities used across packages
- **`errors`** - TensorZero-specific error types and handling

Each package contains comprehensive documentation with detailed field descriptions, usage examples, and best practices.

## Getting Started

### Installation

```bash
go get github.com/denkhaus/tensorzero
```

### Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/denkhaus/tensorzero"
    "github.com/denkhaus/tensorzero/inference"
    "github.com/denkhaus/tensorzero/shared"
)

func main() {
    // Create a new HTTP gateway client
    client := tensorzero.NewHTTPGateway("http://localhost:3000")
    defer client.Close()

    ctx := context.Background()

    // Method 1: Direct struct initialization with utility functions
    resp, err := client.Inference(ctx, &inference.InferenceRequest{
        FunctionName: util.StringPtr("qa_function"), // Use util functions for pointers
        Input: inference.InferenceInput{
            Messages: []shared.Message{
                {
                    Role: "user",
                    Content: []shared.ContentBlock{
                        &shared.Text{Text: "What is the capital of France?"},
                    },
                },
            },
        },
        Tags: map[string]string{
            "user_id": "123",
            "session": "abc",
            "source":  "api",
        },
        Dryrun: util.BoolPtr(false), // Use util.BoolPtr for optional booleans
    })

    // Method 2: Clean options pattern with helper function
    resp2, err := makeInference(ctx, client, 
        "What is the capital of France?",
        inference.WithFunctionName("qa_function"),
        inference.WithTags(map[string]string{
            "user_id": "123",
            "session": "abc",
            "source":  "api",
        }),
        inference.WithDryRun(false),
        inference.WithIncludeOriginalResponse(true),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %+v\n", resp)
}

// Helper function demonstrating clean options pattern usage
func makeInference(ctx context.Context, client tensorzero.Gateway, message string, opts ...inference.InferenceRequestOption) (inference.InferenceResponse, error) {
    req := &inference.InferenceRequest{
        Input: inference.InferenceInput{
            Messages: []shared.Message{
                {
                    Role: "user",
                    Content: []shared.ContentBlock{
                        &shared.Text{Text: message},
                    },
                },
            },
        },
    }
    
    // Apply all options
    for _, opt := range opts {
        opt(req)
    }
    
    return client.Inference(ctx, req)
}

// The util package provides helpful pointer functions:
// util.StringPtr(), util.BoolPtr(), util.IntPtr(), util.UUIDPtr(), util.Float64Ptr()
```

### Advanced Usage

#### Streaming with Error Handling
```go
// Create a streaming request using options pattern
streamReq := &inference.InferenceRequest{
    Input: inference.InferenceInput{
        Messages: []shared.Message{
            {
                Role: "user",
                Content: []shared.ContentBlock{
                    &shared.Text{Text: "Tell me a story"},
                },
            },
        },
    },
}

// Apply streaming options cleanly
applyOptions(streamReq,
    inference.WithFunctionName("story_generator"),
    inference.WithStream(true),
    inference.WithTags(map[string]string{
        "type": "streaming",
        "user": "demo",
    }),
)

chunks, errs := client.InferenceStream(ctx, streamReq)

for {
    select {
    case chunk, ok := <-chunks:
        if !ok {
            return // Stream finished
        }
        // Process chunk using the inference package types
        fmt.Printf("Received chunk: %v\n", chunk)
        
    case err := <-errs:
        if err != nil {
            log.Printf("Stream error: %v", err)
            return
        }
    }
}

// Helper function for applying multiple options
func applyOptions(req *inference.InferenceRequest, opts ...inference.InferenceRequestOption) {
    for _, opt := range opts {
        opt(req)
    }
}

// Example: Create inference with clean options pattern
func createChatInference(ctx context.Context, client tensorzero.Gateway, message string) (inference.InferenceResponse, error) {
    return makeInference(ctx, client, message,
        inference.WithFunctionName("chat_assistant"),
        inference.WithTags(map[string]string{"type": "chat"}),
        inference.WithDryRun(false),
    )
}
```

#### Feedback and Evaluation
```go
import (
    "github.com/denkhaus/tensorzero/feedback"
    "github.com/denkhaus/tensorzero/evaluation"
    "github.com/denkhaus/tensorzero/util"
    "github.com/google/uuid"
)

// Submit feedback using utility functions
inferenceID := uuid.NewV7() // TensorZero requires UUIDv7 for proper timestamp ordering
feedbackResp, err := client.Feedback(ctx, &feedback.Request{
    MetricName:  "user_rating",
    Value:       5,
    InferenceID: util.UUIDPtr(inferenceID), // Use util.UUIDPtr
    Tags: map[string]string{
        "source": "user_feedback",
        "rating_type": "helpfulness",
    },
    Internal: util.BoolPtr(false), // External user feedback
})

// Submit boolean metric feedback
boolFeedback, err := client.Feedback(ctx, &feedback.Request{
    MetricName:  "is_helpful",
    Value:       true,
    InferenceID: util.UUIDPtr(inferenceID),
    Tags: map[string]string{
        "evaluator": "human",
        "confidence": "high",
    },
})

// Dynamic evaluation - test different variants against datasets
evalResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
    Variants: map[string]string{
        "model_a": "gpt-4",
        "model_b": "claude-3",
        "model_c": "gemini-pro",
    },
    DisplayName: util.StringPtr("A/B/C Test: Model Comparison"),
    ProjectName: util.StringPtr("q4-model-evaluation"),
    Tags: map[string]string{
        "experiment": "model_comparison",
        "version": "v1.0",
        "dataset": "qa_benchmark",
    },
})

// Create evaluation episodes
episodeResp, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
    RunID: evalResp.RunID,
    TaskName: util.StringPtr("question_answering"),
    DatapointName: util.StringPtr("sample_qa_001"),
    Tags: map[string]string{
        "difficulty": "medium",
        "category": "factual",
    },
})
```

#### Advanced Filtering
```go
import (
    "github.com/denkhaus/tensorzero/filter"
    "github.com/denkhaus/tensorzero/shared"
    "github.com/denkhaus/tensorzero/util"
    "github.com/google/uuid"
)

// Complex filtering with AND/OR logic using the filter package
complexFilter := filter.NewAndFilter(
    filter.NewTagFilter("environment", "production", "="),
    filter.NewOrFilter(
        filter.NewFloatMetricFilter("accuracy", 0.8, ">="),
        filter.NewBooleanMetricFilter("is_helpful", true),
    ),
)

inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
    FunctionName: util.StringPtr("qa_function"), // Filter by function
    Filter:       complexFilter,
    OrderBy:      &shared.OrderBy{Field: "timestamp", Direction: "desc"},
    Limit:        util.IntPtr(50),
    Offset:       util.IntPtr(0), // For pagination
})

// Example: List inferences for a specific episode
episodeID := uuid.NewV7() // Use UUIDv7 for episode IDs
episodeInferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
    EpisodeID: util.UUIDPtr(episodeID),
    OrderBy:   &shared.OrderBy{Field: "timestamp", Direction: "asc"},
    Limit:     util.IntPtr(100),
})
```

## Development & Testing

This project includes a comprehensive testing framework with automated setup and execution.

### Prerequisites

*   **Go 1.21+** - Latest Go version
*   **Docker & Docker Compose** - Container runtime
*   **OpenRouter API Key** - For model access (create `docker/.env` with `OPENROUTER_API_KEY=your_key`)

### Quick Start

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/denkhaus/tensorzero
    cd tensorzero-go
    ```

2.  **Setup and testing:**
    ```bash
    # Setup development environment
    make setup
    
    # Run all tests
    make test-all
    ```

### Testing Options

#### Unit Tests (Fast, no external dependencies)
```bash
make test-unit
```

#### Integration Tests (Full TensorZero API testing)
```bash
# Run integration tests
make test-integration

# Custom test scenarios
go test -tags=integration ./tests -run TestIntegration_BasicInference -v
```

#### Performance Benchmarks
```bash
go test -tags=integration ./tests -bench=. -benchmem
```

#### All Tests (Unit + Integration)
```bash
make test-all
```

### Test Suite Features

Our comprehensive test suite includes:

#### Unit Tests
- **15+ test files** covering all packages
- **90%+ code coverage** in utility functions
- **Interface compliance** validation
- **Error handling** verification
- **Type safety** testing

#### Integration Tests
- **Real TensorZero API** testing against live services
- **Complete API coverage** - inference, streaming, feedback, datapoints
- **Advanced filtering** - complex queries with AND/OR/NOT logic
- **Performance benchmarks** - latency and throughput measurements
- **Reliability testing** - concurrent requests and stress scenarios
- **Error scenarios** - network failures and API errors

#### Performance & Reliability
- **Concurrent request handling** - Multi-threaded stress testing
- **Streaming performance** - Real-time data flow validation
- **Resource management** - Memory usage and connection cleanup
- **Context cancellation** - Proper timeout handling
- **Long-running operations** - Extended session testing

#### Automated Testing
- **One-command setup** - Automated environment configuration
- **Docker orchestration** - TensorZero services management
- **Health checks** - Service readiness validation
- **Comprehensive reporting** - JSON, JUnit XML, and console output
- **CI/CD ready** - GitHub Actions integration

### Test Documentation

*   **[INTEGRATION_TESTS.md](./docs/INTEGRATION_TESTS.md)** - Comprehensive integration testing guide
*   **[TESTS.md](./docs/TESTS.md)** - General testing information
*   **[PYTHON_SDK_PARITY.md](./docs/PYTHON_SDK_PARITY.md)** - Feature parity documentation

### Development Workflow

```bash
# 1. Setup development environment
make setup

# 2. Run tests during development
make test-unit                    # Fast feedback loop
make test-integration            # API validation
make test-all                    # Complete test suite

# 3. Performance testing
go test -tags=integration ./tests -bench=. -benchmem

# 4. Check service health
make health-check

# 5. View logs if needed
make logs

# 6. Cleanup
make clean
```

## Architecture

The TensorZero Go client is organized into focused packages:

```
github.com/denkhaus/tensorzero/
├── inference/     # Inference operations and types
├── feedback/      # Feedback and metrics
├── datapoint/     # Training data management  
├── filter/        # Query filtering logic
├── shared/        # Common types and utilities
├── config/        # Configuration management
├── tool/          # Tool calling functionality
├── types/         # Request/response types
├── util/          # Helper functions
└── errors/        # Structured error handling
```

### Key Design Principles

- **Go-Idiomatic**: Follows Go best practices and conventions
- **Context-Aware**: All operations support context cancellation
- **Type-Safe**: Strong typing with compile-time validation
- **Interface-Driven**: Extensible design with clear abstractions
- **Error-Transparent**: Detailed error information with proper wrapping
- **Resource-Conscious**: Proper cleanup and connection management

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

1. **Fork and clone** the repository
2. **Run setup script**: `./scripts/setup-integration-tests.sh`
3. **Make changes** and add tests
4. **Run test suite**: `./scripts/run-integration-tests.sh --verbose`
5. **Submit pull request** with comprehensive test coverage

### Code Quality Standards

- **100% test coverage** for new features
- **Integration tests** for API changes
- **Performance benchmarks** for critical paths
- **Documentation updates** for public APIs
- **Go formatting** with `gofmt` and `go vet`

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for more details.