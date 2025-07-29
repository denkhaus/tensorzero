package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configJSON := `{
		"functions": {
			"chat_function": {
				"type": "chat",
				"variants": {
					"default": {
						"type": "chat_completion",
						"model": "gpt-4"
					}
				}
			}
		}
	}`

	var cfg Config
	err := json.Unmarshal([]byte(configJSON), &cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Functions)
	assert.Contains(t, cfg.Functions, "chat_function")

	chatFunc := cfg.Functions["chat_function"]
	assert.NotNil(t, chatFunc)
	assert.Equal(t, "chat", chatFunc.(*ChatFunctionConfig).Type) // Type assertion needed here
	assert.Contains(t, chatFunc.(*ChatFunctionConfig).Variants, "default")
}

func TestChatFunctionConfig(t *testing.T) {
	chatConfigJSON := `{
		"type": "chat",
		"variants": {
			"default": {
				"type": "chat_completion",
				"model": "gpt-4"
			}
		},
		"system_schema": {"type": "object"},
		"user_schema": null
	}`

	var cfg ChatFunctionConfig
	err := json.Unmarshal([]byte(chatConfigJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "chat", cfg.Type)
	assert.NotNil(t, cfg.Variants)
	assert.Contains(t, cfg.Variants, "default")
	assert.NotNil(t, cfg.SystemSchema)
	assert.Nil(t, cfg.UserSchema)
	assert.Nil(t, cfg.AssistantSchema)

	assert.Equal(t, "chat", cfg.GetType())
	assert.Equal(t, cfg.Variants, cfg.GetVariants())
}

func TestJsonFunctionConfig(t *testing.T) {
	jsonConfigJSON := `{
		"type": "json",
		"variants": {
			"default": {
				"type": "chat_completion",
				"model": "gpt-3.5-turbo"
			}
		},
		"output_schema": {"type": "array"}
	}`

	var cfg JsonFunctionConfig
	err := json.Unmarshal([]byte(jsonConfigJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "json", cfg.Type)
	assert.NotNil(t, cfg.Variants)
	assert.Contains(t, cfg.Variants, "default")
	assert.NotNil(t, cfg.OutputSchema)

	assert.Equal(t, "json", cfg.GetType())
	assert.Equal(t, cfg.Variants, cfg.GetVariants())
}

func TestChatCompletionConfig(t *testing.T) {
	chatCompletionJSON := `{
		"type": "chat_completion",
		"system_template": "You are a helpful assistant.",
		"user_template": "Hello!",
		"model": "gpt-4"
	}`

	var cfg ChatCompletionConfig
	err := json.Unmarshal([]byte(chatCompletionJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "chat_completion", cfg.Type)
	assert.NotNil(t, cfg.SystemTemplate)
	assert.Equal(t, "You are a helpful assistant.", *cfg.SystemTemplate)
	assert.NotNil(t, cfg.UserTemplate)
	assert.Equal(t, "Hello!", *cfg.UserTemplate)
	assert.Nil(t, cfg.AssistantTemplate)
	assert.Equal(t, "gpt-4", cfg.Model)

	assert.Equal(t, "chat_completion", cfg.GetType())

	// Test with omitempty fields
	chatCompletionNoTemplatesJSON := `{
		"type": "chat_completion",
		"model": "gpt-4"
	}`
	var cfg2 ChatCompletionConfig
	err = json.Unmarshal([]byte(chatCompletionNoTemplatesJSON), &cfg2)
	assert.NoError(t, err)
	assert.Nil(t, cfg2.SystemTemplate)
	assert.Nil(t, cfg2.UserTemplate)
	assert.Nil(t, cfg2.AssistantTemplate)
}

func TestBestOfNSamplingConfig(t *testing.T) {
	cfg := BestOfNSamplingConfig{Type: "best_of_n"}
	assert.Equal(t, "best_of_n", cfg.GetType())
}

func TestDiclConfig(t *testing.T) {
	cfg := DiclConfig{Type: "dicl"}
	assert.Equal(t, "dicl", cfg.GetType())
}

func TestMixtureOfNConfig(t *testing.T) {
	cfg := MixtureOfNConfig{Type: "mixture_of_n"}
	assert.Equal(t, "mixture_of_n", cfg.GetType())
}

func TestChainOfThoughtConfig(t *testing.T) {
	cfg := ChainOfThoughtConfig{Type: "chain_of_thought"}
	assert.Equal(t, "chain_of_thought", cfg.GetType())
}

func TestOpenAISFTConfig(t *testing.T) {
	sftConfigJSON := `{
		"model": "gpt-3.5-turbo",
		"batch_size": 4,
		"learning_rate_multiplier": 2.0
	}`
	var cfg OpenAISFTConfig
	err := json.Unmarshal([]byte(sftConfigJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "gpt-3.5-turbo", cfg.Model)
	assert.NotNil(t, cfg.BatchSize)
	assert.Equal(t, 4, *cfg.BatchSize)
	assert.NotNil(t, cfg.LearningRateMultiplier)
	assert.Equal(t, 2.0, *cfg.LearningRateMultiplier)

	assert.Equal(t, "openai_sft", cfg.GetType())
}

func TestFireworksSFTConfig(t *testing.T) {
	sftConfigJSON := `{
		"model": "mixtral-8x7b-instruct",
		"account_id": "12345"
	}`
	var cfg FireworksSFTConfig
	err := json.Unmarshal([]byte(sftConfigJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "mixtral-8x7b-instruct", cfg.Model)
	assert.Equal(t, "12345", cfg.AccountID)

	assert.Equal(t, "fireworks_sft", cfg.GetType())
}

func TestGCPVertexGeminiSFTConfig(t *testing.T) {
	sftConfigJSON := `{
		"model": "gemini-pro",
		"bucket_name": "my-bucket",
		"project_id": "my-project",
		"region": "us-central1"
	}`
	var cfg GCPVertexGeminiSFTConfig
	err := json.Unmarshal([]byte(sftConfigJSON), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "gemini-pro", cfg.Model)
	assert.Equal(t, "my-bucket", cfg.BucketName)
	assert.Equal(t, "my-project", cfg.ProjectID)
	assert.Equal(t, "us-central1", cfg.Region)

	assert.Equal(t, "gcp_vertex_gemini_sft", cfg.GetType())
}

func TestOptimizationJobStatus(t *testing.T) {
	assert.Equal(t, "pending", string(OptimizationJobStatusPending))
	assert.Equal(t, "completed", string(OptimizationJobStatusCompleted))
	assert.Equal(t, "failed", string(OptimizationJobStatusFailed))
}

func TestOptimizationJobInfo(t *testing.T) {
	infoJSON := `{
		"message": "Job started",
		"status": "pending",
		"estimated_finish": 1678886400
	}`
	var info OptimizationJobInfo
	err := json.Unmarshal([]byte(infoJSON), &info)
	assert.NoError(t, err)
	assert.Equal(t, "Job started", info.Message)
	assert.Equal(t, OptimizationJobStatusPending, info.Status)
	assert.NotNil(t, info.EstimatedFinish)
	assert.Equal(t, int64(1678886400), *info.EstimatedFinish)
	assert.Nil(t, info.Output)
}
