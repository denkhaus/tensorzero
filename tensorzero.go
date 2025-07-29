// Package tensorzero provides a Go client for the TensorZero AI inference gateway.
//
// TensorZero is an AI inference gateway that provides a unified interface for
// multiple AI model providers with features like A/B testing, optimization,
// and observability.
//
// Basic usage:
//
//	client := tensorzero.NewHTTPGateway("http://localhost:3000")
//	defer client.Close()
//
//	response, err := client.Inference(context.Background(), &tensorzero.InferenceRequest{
//		Input: tensorzero.InferenceInput{
//			Messages: []tensorzero.Message{
//				{
//					Role: "user",
//					Content: []tensorzero.ContentBlock{
//						tensorzero.NewText("Hello, world!"),
//					},
//				},
//			},
//		},
//		FunctionName: tensorzero.StringPtr("my_function"),
//	})
//
// For streaming responses:
//
//	chunks, errs := client.InferenceStream(context.Background(), &tensorzero.InferenceRequest{
//		Input: tensorzero.InferenceInput{
//			Messages: []tensorzero.Message{
//				{
//					Role: "user",
//					Content: []tensorzero.ContentBlock{
//						tensorzero.NewText("Tell me a story"),
//					},
//				},
//			},
//		},
//		FunctionName: tensorzero.StringPtr("story_function"),
//		Stream:       tensorzero.BoolPtr(true),
//	})
//
//	for {
//		select {
//		case chunk, ok := <-chunks:
//			if !ok {
//				return // Stream finished
//			}
//			// Process chunk
//		case err := <-errs:
//			if err != nil {
//				// Handle error
//			}
//		}
//	}
package tensorzero

const (
	// Version is the current version of the tensorzero-go package
	Version = "0.1.0"
)

// Default configuration values
const (
	DefaultTimeout = 30 // seconds
	DefaultBaseURL = "http://localhost:3000"
)

// Role constants for messages
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

// Tool choice constants
const (
	ToolChoiceAuto     = "auto"
	ToolChoiceRequired = "required"
	ToolChoiceOff      = "off"
)