// Package feedback provides types and functionality for TensorZero feedback operations.
// This includes feedback requests, responses, and metric handling.
package feedback

import "github.com/google/uuid"

// Request represents a feedback request
type Request struct {
	MetricName  string            `json:"metric_name"`
	Value       interface{}       `json:"value"`
	InferenceID *uuid.UUID        `json:"inference_id,omitempty"`
	EpisodeID   *uuid.UUID        `json:"episode_id,omitempty"`
	Dryrun      *bool             `json:"dryrun,omitempty"`
	Internal    *bool             `json:"internal,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// Response represents feedback response
type Response struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
}
