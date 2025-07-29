# TensorZero Go Client

[![GitHub issues](https://img.shields.io/github/issues/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/pulls)
[![GitHub contributors](https://img.shields.io/github/contributors/denkhaus/tensorzero)](https://github.com/denkhaus/tensorzero/graphs/contributors)

**Disclaimer:** This is an unofficial, community-maintained Go client for the TensorZero API.

A comprehensive Go client library for [TensorZero](https://github.com/tensorzero/tensorzero), an AI inference gateway that provides a unified interface for multiple AI model providers with features like A/B testing, optimization, and observability.

This Go SDK is a complete and accurate port of the Python SDK, providing full feature parity for the Go ecosystem. For a detailed breakdown of the feature parity, see [PYTHON_SDK_PARITY.md](./docs/PYTHON_SDK_PARITY.md).

## Features

*   **Complete TensorZero API Client:** Implements all TensorZero API endpoints, including inference, streaming, feedback, and datapoint management.
*   **OpenAI SDK Compatibility:** Designed to be compatible with the OpenAI SDK.
*   **Go-Idiomatic Design:** Leverages Go's best practices, including context support, channel-based streaming, and structured error handling.
*   **Comprehensive Testing:** Includes a robust test suite with unit and integration tests. See [TESTS.md](./docs/TESTS.md) for more details.
*   **Dockerized Development Environment:** Comes with a Docker-based development environment for easy setup and testing.

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
)

func main() {
    // Create a new HTTP gateway client
    client := tensorzero.NewHTTPGateway("http://localhost:3000")
    defer client.Close()

    // Create an inference request
    request := &tensorzero.InferenceRequest{
        Input: tensorzero.InferenceInput{
            Messages: []tensorzero.Message{
                {
                    Role: tensorzero.RoleUser,
                    Content: []tensorzero.ContentBlock{
                        tensorzero.NewText("What is the capital of France?"),
                    },
                },
            },
        },
        FunctionName: tensorzero.StringPtr("qa_function"),
    }

    // Make the inference request
    response, err := client.Inference(context.Background(), request)
    if err != nil {
        log.Fatal(err)
    }

    // Handle the response
    switch resp := response.(type) {
    case *tensorzero.ChatInferenceResponse:
        fmt.Printf("Response: %v
", resp.Content)
    case *tensorzero.JsonInferenceResponse:
        fmt.Printf("JSON Response: %v
", resp.Output)
    }
}
```

## Development & Testing

This project includes a Docker-based development environment that simplifies setup and testing.

### Prerequisites

*   Go 1.21 or later
*   Docker and Docker Compose
*   OpenRouter API key (for model access)

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/denkhaus/tensorzero
    cd tensorzero-go
    ```

2.  **Set up the environment:**
    ```bash
    export OPENROUTER_API_KEY=your_api_key
    make up
    ```

### Running Tests

The project includes a comprehensive test suite. You can run unit tests quickly without external dependencies:

```bash
make test-unit
```

To run all tests, including integration tests that require a running TensorZero instance (managed by Docker Compose), use:

```bash
make test-all
```

For more detailed information on the tests, see [TESTS.md](./docs/TESTS.md).

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for more details.
