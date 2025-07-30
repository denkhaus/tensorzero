// Package feedback provides types and functionality for TensorZero feedback operations.
// This includes feedback requests, responses, and metric handling.
package feedback

import "github.com/google/uuid"

// Request represents a feedback request used to provide feedback on inferences or episodes.
// Feedback is essential for measuring and improving the performance of your AI functions.
// Each feedback is associated with a metric defined in your TensorZero configuration.
type Request struct {
	// MetricName is the name of the metric to provide feedback for.
	// This must match a metric defined in your TensorZero configuration file.
	// Special reserved metric names include "comment" for free-form text feedback
	// and "demonstration" for providing example outputs.
	MetricName string `json:"metric_name"`

	// Value is the feedback value whose type depends on the metric type:
	// - Boolean metrics: true/false values
	// - Float metrics: numeric values (e.g., ratings, scores)
	// - Comment metrics: string values with free-form text
	// - Demonstration metrics: objects that represent valid outputs
	Value interface{} `json:"value"`

	// InferenceID is the unique identifier of the inference to provide feedback for.
	// Use this field when the metric level is "inference". Only use inference IDs
	// that were returned by the TensorZero gateway. Either InferenceID or EpisodeID
	// must be provided, but not both.
	InferenceID *uuid.UUID `json:"inference_id,omitempty"`

	// EpisodeID is the unique identifier of the episode to provide feedback for.
	// Use this field when the metric level is "episode". Only use episode IDs
	// that were returned by the TensorZero gateway. Either InferenceID or EpisodeID
	// must be provided, but not both.
	EpisodeID *uuid.UUID `json:"episode_id,omitempty"`

	// Dryrun, when set to true, executes the feedback request without storing it
	// to the database. This is primarily for debugging and testing purposes
	// and should not be used in production environments.
	Dryrun *bool `json:"dryrun,omitempty"`

	// Internal indicates whether this feedback is generated internally by the system
	// rather than from external sources like user interactions. This helps distinguish
	// between automated and human-generated feedback.
	Internal *bool `json:"internal,omitempty"`

	// Tags are user-provided key-value pairs to associate with the feedback.
	// These can be used for categorization, filtering, and analysis purposes.
	// Example: {"user_id": "123", "source": "user_rating", "version": "v2.1"}
	Tags map[string]string `json:"tags,omitempty"`
}

// Response represents the response returned after successfully submitting feedback.
// This response contains the unique identifier assigned to the feedback entry.
type Response struct {
	// FeedbackID is the unique identifier (UUIDv7) assigned to the feedback.
	// This ID can be used to reference, update, or track this specific feedback
	// entry in subsequent operations or analytics.
	FeedbackID uuid.UUID `json:"feedback_id"`
}
