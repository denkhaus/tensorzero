package types

import (
	"github.com/google/uuid"
)

// JsonInferenceOutput represents JSON inference output
type JsonInferenceOutput struct {
	Raw    *string                `json:"raw,omitempty"`
	Parsed map[string]interface{} `json:"parsed,omitempty"`
}

// ChatInferenceResponse represents a chat inference response
type ChatInferenceResponse struct {
	InferenceID      uuid.UUID      `json:"inference_id"`
	EpisodeID        uuid.UUID      `json:"episode_id"`
	VariantName      string         `json:"variant_name"`
	Content          []ContentBlock `json:"content"`
	Usage            Usage          `json:"usage"`
	FinishReason     *FinishReason  `json:"finish_reason,omitempty"`
	OriginalResponse *string        `json:"original_response,omitempty"`
}

// JsonInferenceResponse represents a JSON inference response
type JsonInferenceResponse struct {
	InferenceID      uuid.UUID           `json:"inference_id"`
	EpisodeID        uuid.UUID           `json:"episode_id"`
	VariantName      string              `json:"variant_name"`
	Output           JsonInferenceOutput `json:"output"`
	Usage            Usage               `json:"usage"`
	FinishReason     *FinishReason       `json:"finish_reason,omitempty"`
	OriginalResponse *string             `json:"original_response,omitempty"`
}

func (c *ChatInferenceResponse) GetInferenceID() uuid.UUID      { return c.InferenceID }
func (c *ChatInferenceResponse) GetEpisodeID() uuid.UUID        { return c.EpisodeID }
func (c *ChatInferenceResponse) GetVariantName() string         { return c.VariantName }
func (c *ChatInferenceResponse) GetUsage() Usage                { return c.Usage }
func (c *ChatInferenceResponse) GetFinishReason() *FinishReason { return c.FinishReason }
func (c *ChatInferenceResponse) GetOriginalResponse() *string   { return c.OriginalResponse }

func (j *JsonInferenceResponse) GetInferenceID() uuid.UUID      { return j.InferenceID }
func (j *JsonInferenceResponse) GetEpisodeID() uuid.UUID        { return j.EpisodeID }
func (j *JsonInferenceResponse) GetVariantName() string         { return j.VariantName }
func (j *JsonInferenceResponse) GetUsage() Usage                { return j.Usage }
func (j *JsonInferenceResponse) GetFinishReason() *FinishReason { return j.FinishReason }
func (j *JsonInferenceResponse) GetOriginalResponse() *string   { return j.OriginalResponse }

// ChatChunk represents streaming chat chunk
type ChatChunk struct {
	InferenceID  uuid.UUID           `json:"inference_id"`
	EpisodeID    uuid.UUID           `json:"episode_id"`
	VariantName  string              `json:"variant_name"`
	Content      []ContentBlockChunk `json:"content"`
	Usage        *Usage              `json:"usage,omitempty"`
	FinishReason *FinishReason       `json:"finish_reason,omitempty"`
}

// JsonChunk represents streaming JSON chunk
type JsonChunk struct {
	InferenceID  uuid.UUID     `json:"inference_id"`
	EpisodeID    uuid.UUID     `json:"episode_id"`
	VariantName  string        `json:"variant_name"`
	Raw          string        `json:"raw"`
	Usage        *Usage        `json:"usage,omitempty"`
	FinishReason *FinishReason `json:"finish_reason,omitempty"`
}

func (c *ChatChunk) GetInferenceID() uuid.UUID { return c.InferenceID }
func (c *ChatChunk) GetEpisodeID() uuid.UUID   { return c.EpisodeID }
func (c *ChatChunk) GetVariantName() string    { return c.VariantName }

func (j *JsonChunk) GetInferenceID() uuid.UUID { return j.InferenceID }
func (j *JsonChunk) GetEpisodeID() uuid.UUID   { return j.EpisodeID }
func (j *JsonChunk) GetVariantName() string    { return j.VariantName }

// FeedbackResponse represents feedback response
type FeedbackResponse struct {
	FeedbackID uuid.UUID `json:"feedback_id"`
}

// DynamicEvaluationRunResponse represents dynamic evaluation run response
type DynamicEvaluationRunResponse struct {
	RunID uuid.UUID `json:"run_id"`
}

// DynamicEvaluationRunEpisodeResponse represents dynamic evaluation run episode response
type DynamicEvaluationRunEpisodeResponse struct {
	EpisodeID uuid.UUID `json:"episode_id"`
}
