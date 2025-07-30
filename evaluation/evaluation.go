// Package evaluation provides types and functionality for TensorZero dynamic evaluation operations.
// This includes evaluation run requests, responses, and episode management.
package evaluation

import (
	"github.com/google/uuid"
)

// RunRequest represents a dynamic evaluation run request used to create a new evaluation run.
// Dynamic evaluation runs allow you to test different variants of your functions against
// specific datasets or tasks to measure their performance and effectiveness.
type RunRequest struct {
	// Variants is a map of variant names to their configurations for this evaluation run.
	// Each key represents a variant name, and the value specifies the variant configuration
	// to be used during the evaluation process.
	Variants map[string]string `json:"variants"`

	// Tags are user-provided key-value pairs to associate with the evaluation run.
	// These can be used for categorization, filtering, and tracking purposes.
	// Example: {"experiment": "A/B_test", "version": "v1.2"}
	Tags map[string]string `json:"tags,omitempty"`

	// ProjectName optionally specifies the project this evaluation run belongs to.
	// This helps organize evaluation runs within larger projects or experiments.
	ProjectName *string `json:"project_name,omitempty"`

	// DisplayName provides a human-readable name for the evaluation run.
	// This makes it easier to identify and manage evaluation runs in dashboards and reports.
	DisplayName *string `json:"display_name,omitempty"`
}

// RunResponse represents the response returned after successfully creating a dynamic evaluation run.
// This response contains the unique identifier for the newly created evaluation run.
type RunResponse struct {
	// RunID is the unique identifier (UUIDv7) assigned to the evaluation run.
	// This ID should be used in subsequent requests to create episodes within this run
	// or to reference this specific evaluation run in other operations.
	RunID uuid.UUID `json:"run_id"`
}

// EpisodeRequest represents a request to create an episode within a dynamic evaluation run.
// Episodes are individual test cases or scenarios within an evaluation run that test
// specific functionality or behavior of the variants being evaluated.
type EpisodeRequest struct {
	// RunID is the unique identifier of the evaluation run this episode belongs to.
	// This must be a valid RunID returned from a previous RunRequest.
	RunID uuid.UUID `json:"run_id"`

	// TaskName optionally specifies the name of the task or test case this episode represents.
	// This helps categorize and organize different types of evaluations within a run.
	TaskName *string `json:"task_name,omitempty"`

	// DatapointName optionally specifies the name of the specific datapoint being evaluated.
	// This is useful when testing against known datasets or specific test cases.
	DatapointName *string `json:"datapoint_name,omitempty"`

	// Tags are user-provided key-value pairs to associate with this specific episode.
	// These can be used for filtering, categorization, and analysis of episode results.
	// Example: {"difficulty": "hard", "category": "reasoning"}
	Tags map[string]string `json:"tags,omitempty"`
}

// EpisodeResponse represents the response returned after successfully creating an evaluation episode.
// This response contains the unique identifier for the newly created episode within the evaluation run.
type EpisodeResponse struct {
	// EpisodeID is the unique identifier (UUIDv7) assigned to the evaluation episode.
	// This ID can be used to track the progress and results of this specific episode,
	// associate feedback with the episode, or reference it in subsequent operations.
	EpisodeID uuid.UUID `json:"episode_id"`
}