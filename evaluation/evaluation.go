// Package evaluation provides types and functionality for TensorZero dynamic evaluation operations.
// This includes evaluation run requests, responses, and episode management.
package evaluation

import (
	"github.com/google/uuid"
)

// RunRequest represents a dynamic evaluation run request
type RunRequest struct {
	Variants    map[string]string `json:"variants"`
	Tags        map[string]string `json:"tags,omitempty"`
	ProjectName *string           `json:"project_name,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
}

// RunResponse represents dynamic evaluation run response
type RunResponse struct {
	RunID uuid.UUID `json:"run_id"`
}

// EpisodeRequest represents a dynamic evaluation run episode request
type EpisodeRequest struct {
	RunID         uuid.UUID         `json:"run_id"`
	TaskName      *string           `json:"task_name,omitempty"`
	DatapointName *string           `json:"datapoint_name,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// EpisodeResponse represents dynamic evaluation run episode response
type EpisodeResponse struct {
	EpisodeID uuid.UUID `json:"episode_id"`
}