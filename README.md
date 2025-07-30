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

    // Create an inference request with comprehensive documentation
    // Each field is well-documented with usage examples and constraints
    resp, err := client.Inference(ctx, &inference.InferenceRequest{
        FunctionName: stringPtr("qa_function"), // Function defined in TensorZero config
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
        // Optional: Add tags for tracking and analysis
        Tags: map[string]string{
            "user_id": "123",
            "session": "abc",
            "source": "api",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %+v\n", resp)
}

func stringPtr(s string) *string {
    return &s
}
    response, err := client.Inference(context.Background(), request)
    if err != nil {
        log.Fatal(err)
    }

    // Handle the response
    switch resp := response.(type) {
    case *tensorzero.ChatInferenceResponse:
        fmt.Printf("Response: %v\n", resp.Content)
    case *tensorzero.JsonInferenceResponse:
        fmt.Printf("JSON Response: %v\n", resp.Output)
    }
}
```

### Advanced Usage

#### Streaming with Error Handling
```go
chunks, errs := client.InferenceStream(context.Background(), request)

for {
    select {
    case chunk, ok := <-chunks:
        if !ok {
            return // Stream finished
        }
        // Process chunk
        fmt.Printf("Received: %v\n", chunk)
        
    case err := <-errs:
        if err != nil {
            log.Printf("Stream error: %v", err)
            return
        }
    }
}
```

#### Feedback and Evaluation
```go
// Submit feedback
feedbackResp, err := client.Feedback(ctx, &feedback.Request{
    MetricName:  "user_rating",
    Value:       5,
    InferenceID: &inferenceID,
})

// Dynamic evaluation - test different variants against datasets
evalResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
    Variants: map[string]string{
        "model_a": "gpt-4",
        "model_b": "claude-3",
    },
    DisplayName: stringPtr("A/B Test: GPT-4 vs Claude-3"),
    Tags: map[string]string{
        "experiment": "model_comparison",
        "version": "v1.0",
    },
})
```

#### Advanced Filtering
```go
// Complex filtering with AND/OR logic
complexFilter := filter.NewAndFilter(
    filter.NewTagFilter("environment", "production", "="),
    filter.NewOrFilter(
        filter.NewFloatMetricFilter("accuracy", 0.8, ">="),
        filter.NewBooleanMetricFilter("is_helpful", true),
    ),
)

inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
    Filters: []filter.InferenceFilterTreeNode{complexFilter},
    OrderBy: shared.NewOrderByTimestamp("desc"),
    Limit:   util.IntPtr(50),
})
```

## Development & Testing

This project includes a comprehensive testing framework with automated setup and execution.

### Prerequisites

*   **Go 1.21+** - Latest Go version
*   **Docker & Docker Compose** - Container runtime
*   **OpenRouter API Key** - For model access (optional for dry-run tests)

### Quick Start

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/denkhaus/tensorzero
    cd tensorzero-go
    ```

2.  **Automated setup and testing:**
    ```bash
    # One-command setup and test execution
    ./scripts/setup-integration-tests.sh
    ./scripts/run-integration-tests.sh --verbose
    ```

### Testing Options

#### Unit Tests (Fast, no external dependencies)
```bash
make test-unit
```

#### Integration Tests (Full TensorZero API testing)
```bash
# Quick integration tests
make test-integration

# Comprehensive test suite with benchmarks
./scripts/run-integration-tests.sh --verbose --benchmark

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
./scripts/setup-integration-tests.sh

# 2. Run tests during development
make test-unit                    # Fast feedback loop
make test-integration            # API validation
./scripts/run-integration-tests.sh --verbose  # Full test suite

# 3. Performance testing
go test -tags=integration ./tests -bench=. -benchmem

# 4. Cleanup
make docker-down
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