package tensorzero

// Config represents TensorZero configuration
type Config struct {
	Functions FunctionsConfig `json:"functions"`
}

// FunctionsConfig represents function configurations
type FunctionsConfig map[string]FunctionConfig

// FunctionConfig represents a function configuration
type FunctionConfig interface {
	GetType() string
	GetVariants() VariantsConfig
}

// ChatFunctionConfig represents a chat function configuration
type ChatFunctionConfig struct {
	Type             string                 `json:"type"`
	Variants         VariantsConfig         `json:"variants"`
	SystemSchema     map[string]interface{} `json:"system_schema,omitempty"`
	UserSchema       map[string]interface{} `json:"user_schema,omitempty"`
	AssistantSchema  map[string]interface{} `json:"assistant_schema,omitempty"`
}

func (c *ChatFunctionConfig) GetType() string {
	return c.Type
}

func (c *ChatFunctionConfig) GetVariants() VariantsConfig {
	return c.Variants
}

// JsonFunctionConfig represents a JSON function configuration
type JsonFunctionConfig struct {
	Type             string                 `json:"type"`
	Variants         VariantsConfig         `json:"variants"`
	SystemSchema     map[string]interface{} `json:"system_schema,omitempty"`
	UserSchema       map[string]interface{} `json:"user_schema,omitempty"`
	AssistantSchema  map[string]interface{} `json:"assistant_schema,omitempty"`
	OutputSchema     map[string]interface{} `json:"output_schema,omitempty"`
}

func (j *JsonFunctionConfig) GetType() string {
	return j.Type
}

func (j *JsonFunctionConfig) GetVariants() VariantsConfig {
	return j.Variants
}

// VariantsConfig represents variant configurations
type VariantsConfig map[string]VariantConfig

// VariantConfig represents a variant configuration
type VariantConfig interface {
	GetType() string
}

// ChatCompletionConfig represents a chat completion variant
type ChatCompletionConfig struct {
	Type               string  `json:"type"`
	SystemTemplate     *string `json:"system_template,omitempty"`
	UserTemplate       *string `json:"user_template,omitempty"`
	AssistantTemplate  *string `json:"assistant_template,omitempty"`
	Model              string  `json:"model"`
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

// OptimizationConfig represents optimization configurations
type OptimizationConfig interface {
	GetType() string
}

// OpenAISFTConfig represents OpenAI SFT optimization configuration
type OpenAISFTConfig struct {
	Model                   string   `json:"model"`
	BatchSize               *int     `json:"batch_size,omitempty"`
	LearningRateMultiplier  *float64 `json:"learning_rate_multiplier,omitempty"`
	NEpochs                 *int     `json:"n_epochs,omitempty"`
	Credentials             *string  `json:"credentials,omitempty"`
	APIBase                 *string  `json:"api_base,omitempty"`
	Seed                    *int     `json:"seed,omitempty"`
	Suffix                  *string  `json:"suffix,omitempty"`
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
	Model                      string   `json:"model"`
	BucketName                 string   `json:"bucket_name"`
	ProjectID                  string   `json:"project_id"`
	Region                     string   `json:"region"`
	LearningRateMultiplier     *float64 `json:"learning_rate_multiplier,omitempty"`
	AdapterSize                *int     `json:"adapter_size,omitempty"`
	NEpochs                    *int     `json:"n_epochs,omitempty"`
	ExportLastCheckpointOnly   *bool    `json:"export_last_checkpoint_only,omitempty"`
	Credentials                *string  `json:"credentials,omitempty"`
	APIBase                    *string  `json:"api_base,omitempty"`
	Seed                       *int     `json:"seed,omitempty"`
	ServiceAccount             *string  `json:"service_account,omitempty"`
	KMSKeyName                 *string  `json:"kms_key_name,omitempty"`
	TunedModelDisplayName      *string  `json:"tuned_model_display_name,omitempty"`
	BucketPathPrefix           *string  `json:"bucket_path_prefix,omitempty"`
}

func (g *GCPVertexGeminiSFTConfig) GetType() string {
	return "gcp_vertex_gemini_sft"
}

// OptimizationJobHandle represents an optimization job handle
type OptimizationJobHandle interface {
	GetType() string
	GetJobID() string
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
	Message          string                 `json:"message"`
	Status           OptimizationJobStatus  `json:"status"`
	Output           interface{}            `json:"output,omitempty"`
	EstimatedFinish  *int64                 `json:"estimated_finish,omitempty"`
}