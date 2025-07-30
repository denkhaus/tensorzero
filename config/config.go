package config

import (
	"encoding/json"
	"fmt"
)

// FunctionConfig represents a function configuration
type FunctionConfig interface {
	GetType() string
	GetVariants() VariantsConfig
}

type VariantConfig interface {
	GetType() string
}

// OptimizationConfig represents optimization configurations
type OptimizationConfig interface {
	GetType() string
}

// Config represents TensorZero configuration
type Config struct {
	Functions FunctionsConfig `json:"functions"`
}

// FunctionsConfig represents function configurations
type FunctionsConfig map[string]FunctionConfig

func (fc *FunctionsConfig) UnmarshalJSON(data []byte) error {
	rawMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	*fc = make(FunctionsConfig)
	for name, rawFunc := range rawMap {
		var typeField struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawFunc, &typeField); err != nil {
			return err
		}

		var functionConfig FunctionConfig
		switch typeField.Type {
		case "chat":
			var chatCfg ChatFunctionConfig
			if err := json.Unmarshal(rawFunc, &chatCfg); err != nil {
				return err
			}
			functionConfig = &chatCfg
		case "json":
			var jsonCfg JsonFunctionConfig
			if err := json.Unmarshal(rawFunc, &jsonCfg); err != nil {
				return err
			}
			functionConfig = &jsonCfg
		default:
			// Handle unknown types or return an error
			// For now, we'll just skip or unmarshal into a generic type if necessary
			// or simply return an error for unsupported types.
			return fmt.Errorf("unknown function type: %s", typeField.Type)
		}
		(*fc)[name] = functionConfig
	}
	return nil
}

// ChatFunctionConfig represents a chat function configuration
type ChatFunctionConfig struct {
	Type            string                 `json:"type"`
	Variants        VariantsConfig         `json:"variants"`
	SystemSchema    map[string]interface{} `json:"system_schema,omitempty"`
	UserSchema      map[string]interface{} `json:"user_schema,omitempty"`
	AssistantSchema map[string]interface{} `json:"assistant_schema,omitempty"`
}

func (c *ChatFunctionConfig) GetType() string {
	return c.Type
}

func (c *ChatFunctionConfig) GetVariants() VariantsConfig {
	return c.Variants
}

// JsonFunctionConfig represents a JSON function configuration
type JsonFunctionConfig struct {
	Type            string                 `json:"type"`
	Variants        VariantsConfig         `json:"variants"`
	SystemSchema    map[string]interface{} `json:"system_schema,omitempty"`
	UserSchema      map[string]interface{} `json:"user_schema,omitempty"`
	AssistantSchema map[string]interface{} `json:"assistant_schema,omitempty"`
	OutputSchema    map[string]interface{} `json:"output_schema,omitempty"`
}

func (j *JsonFunctionConfig) GetType() string {
	return j.Type
}

func (j *JsonFunctionConfig) GetVariants() VariantsConfig {
	return j.Variants
}

// VariantsConfig represents variant configurations
type VariantsConfig map[string]VariantConfig

func (vc *VariantsConfig) UnmarshalJSON(data []byte) error {
	rawMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	*vc = make(VariantsConfig)
	for name, rawVariant := range rawMap {
		var typeField struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawVariant, &typeField); err != nil {
			return err
		}

		var variantConfig VariantConfig
		switch typeField.Type {
		case "chat_completion":
			var chatCompletionCfg ChatCompletionConfig
			if err := json.Unmarshal(rawVariant, &chatCompletionCfg); err != nil {
				return err
			}
			variantConfig = &chatCompletionCfg
		case "best_of_n":
			var bestOfNCfg BestOfNSamplingConfig
			if err := json.Unmarshal(rawVariant, &bestOfNCfg); err != nil {
				return err
			}
			variantConfig = &bestOfNCfg
		case "dicl":
			var diclCfg DiclConfig
			if err := json.Unmarshal(rawVariant, &diclCfg); err != nil {
				return err
			}
			variantConfig = &diclCfg
		case "mixture_of_n":
			var mixtureOfNCfg MixtureOfNConfig
			if err := json.Unmarshal(rawVariant, &mixtureOfNCfg); err != nil {
				return err
			}
			variantConfig = &mixtureOfNCfg
		case "chain_of_thought":
			var chainOfThoughtCfg ChainOfThoughtConfig
			if err := json.Unmarshal(rawVariant, &chainOfThoughtCfg); err != nil {
				return err
			}
			variantConfig = &chainOfThoughtCfg
		default:
			return fmt.Errorf("unknown variant type: %s", typeField.Type)
		}
		(*vc)[name] = variantConfig
	}
	return nil
}

// ChatCompletionConfig represents a chat completion variant
type ChatCompletionConfig struct {
	Type              string  `json:"type"`
	SystemTemplate    *string `json:"system_template,omitempty"`
	UserTemplate      *string `json:"user_template,omitempty"`
	AssistantTemplate *string `json:"assistant_template,omitempty"`
	Model             string  `json:"model"`
}

func (c *ChatCompletionConfig) GetType() string {
	return c.Type
}

// BestOfNSamplingConfig represents a best-of-n sampling variant
type BestOfNSamplingConfig struct {
	Type string `json:"type"`
	// Add specific fields as needed
}

func (b *BestOfNSamplingConfig) GetType() string {
	return b.Type
}

// DiclConfig represents a DICL variant
type DiclConfig struct {
	Type string `json:"type"`
	// Add specific fields as needed
}

func (d *DiclConfig) GetType() string {
	return d.Type
}

// MixtureOfNConfig represents a mixture-of-n variant
type MixtureOfNConfig struct {
	Type string `json:"type"`
	// Add specific fields as needed
}

func (m *MixtureOfNConfig) GetType() string {
	return m.Type
}

// ChainOfThoughtConfig represents a chain-of-thought variant
type ChainOfThoughtConfig struct {
	Type string `json:"type"`
	// Add specific fields as needed
}

func (c *ChainOfThoughtConfig) GetType() string {
	return c.Type
}

// OpenAISFTConfig represents OpenAI SFT optimization configuration
type OpenAISFTConfig struct {
	Model                  string   `json:"model"`
	BatchSize              *int     `json:"batch_size,omitempty"`
	LearningRateMultiplier *float64 `json:"learning_rate_multiplier,omitempty"`
	NEpochs                *int     `json:"n_epochs,omitempty"`
	Credentials            *string  `json:"credentials,omitempty"`
	APIBase                *string  `json:"api_base,omitempty"`
	Seed                   *int     `json:"seed,omitempty"`
	Suffix                 *string  `json:"suffix,omitempty"`
}

func (o *OpenAISFTConfig) GetType() string {
	return "openai_sft"
}

// FireworksSFTConfig represents Fireworks SFT optimization configuration
type FireworksSFTConfig struct {
	Model       string  `json:"model"`
	Credentials *string `json:"credentials,omitempty"`
	AccountID   string  `json:"account_id"`
	APIBase     *string `json:"api_base,omitempty"`
}

func (f *FireworksSFTConfig) GetType() string {
	return "fireworks_sft"
}

// GCPVertexGeminiSFTConfig represents GCP Vertex Gemini SFT optimization configuration
type GCPVertexGeminiSFTConfig struct {
	Model                    string   `json:"model"`
	BucketName               string   `json:"bucket_name"`
	ProjectID                string   `json:"project_id"`
	Region                   string   `json:"region"`
	LearningRateMultiplier   *float64 `json:"learning_rate_multiplier,omitempty"`
	AdapterSize              *int     `json:"adapter_size,omitempty"`
	NEpochs                  *int     `json:"n_epochs,omitempty"`
	ExportLastCheckpointOnly *bool    `json:"export_last_checkpoint_only,omitempty"`
	Credentials              *string  `json:"credentials,omitempty"`
	APIBase                  *string  `json:"api_base,omitempty"`
	Seed                     *int     `json:"seed,omitempty"`
	ServiceAccount           *string  `json:"service_account,omitempty"`
	KMSKeyName               *string  `json:"kms_key_name,omitempty"`
	TunedModelDisplayName    *string  `json:"tuned_model_display_name,omitempty"`
	BucketPathPrefix         *string  `json:"bucket_path_prefix,omitempty"`
}

func (g *GCPVertexGeminiSFTConfig) GetType() string {
	return "gcp_vertex_gemini_sft"
}

// OptimizationJobStatus represents optimization job status
type OptimizationJobStatus string

const (
	OptimizationJobStatusPending   OptimizationJobStatus = "pending"
	OptimizationJobStatusCompleted OptimizationJobStatus = "completed"
	OptimizationJobStatusFailed    OptimizationJobStatus = "failed"
)

// OptimizationJobInfo represents optimization job information
type OptimizationJobInfo struct {
	Message         string                `json:"message"`
	Status          OptimizationJobStatus `json:"status"`
	Output          interface{}           `json:"output,omitempty"`
	EstimatedFinish *int64                `json:"estimated_finish,omitempty"`
}
