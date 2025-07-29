import re
import json

go_files = {
    "types/config.go": """
  1 | package types
  2 |
  3 | // Config represents TensorZero configuration
  4 | type Config struct {
  5 | 	Functions FunctionsConfig `json:"functions"`
  6 | }
  7 |
  8 | // FunctionsConfig represents function configurations
  9 | type FunctionsConfig map[string]FunctionConfig
 10 |
 11 | // ChatFunctionConfig represents a chat function configuration
 12 | type ChatFunctionConfig struct {
 13 | 	Type            string                 `json:"type"`
 14 | 	Variants        VariantsConfig         `json:"variants"`
 15 | 	SystemSchema    map[string]interface{} `json:"system_schema,omitempty"`
 16 | 	UserSchema      map[string]interface{} `json:"user_schema,omitempty"`
 17 | 	AssistantSchema map[string]interface{} `json:"assistant_schema,omitempty"`
 18 | }
 19 |
 20 | func (c *ChatFunctionConfig) GetType() string {
 21 | 	return c.Type
 22 | }
 23 |
 24 | func (c *ChatFunctionConfig) GetVariants() VariantsConfig {
 25 | 	return c.Variants
 26 | }
 27 |
 28 | // JsonFunctionConfig represents a JSON function configuration
 29 | type JsonFunctionConfig struct {
 30 | 	Type            string                 `json:"type"`
 31 | 	Variants        VariantsConfig         `json:"variants"`
 32 | 	SystemSchema    map[string]interface{} `json:"system_schema,omitempty"`
 33 | 	UserSchema      map[string]interface{} `json:"user_schema,omitempty"`
 34 | 	AssistantSchema map[string]interface{} `json:"assistant_schema,omitempty"`
 35 | 	OutputSchema    map[string]interface{} `json:"output_schema,omitempty"`
 36 | }
 37 |
 38 | func (j *JsonFunctionConfig) GetType() string {
 39 | 	return j.Type
 40 | }
 41 |
 42 | func (j *JsonFunctionConfig) GetVariants() VariantsConfig {
 43 | 	return j.Variants
 44 | }
 45 |
 46 | // VariantsConfig represents variant configurations
 47 | type VariantsConfig map[string]VariantConfig
 48 |
 49 | // ChatCompletionConfig represents a chat completion variant
 50 | type ChatCompletionConfig struct {
 51 | 	Type              string  `json:"type"`
 52 | 	SystemTemplate    *string `json:"system_template,omitempty"`
 53 | 	UserTemplate      *string `json:"user_template,omitempty"`
 54 | 	AssistantTemplate *string `json:"assistant_template,omitempty"`
 55 | 	Model             string  `json:"model"`
 56 | }
 57 |
 58 | func (c *ChatCompletionConfig) GetType() string {
 59 | 	return c.Type
 60 | }
 61 |
 62 | // BestOfNSamplingConfig represents a best-of-n sampling variant
 63 | type BestOfNSamplingConfig struct {
 64 | 	Type string `json:"type"`
 65 | 	// Add specific fields as needed
 66 | }
 67 |
 68 | func (b *BestOfNSamplingConfig) GetType() string {
 69 | 	return b.Type
 70 | }
 71 |
 72 | // DiclConfig represents a DICL variant
 73 | type DiclConfig struct {
 74 | 	Type string `json:"type"`
 75 | 	// Add specific fields as needed
 76 | }
 77 |
 78 | func (d *DiclConfig) GetType() string {
 79 | 	return d.Type
 80 | }
 81 |
 82 | // MixtureOfNConfig represents a mixture-of-n variant
 83 | type MixtureOfNConfig struct {
 84 | 	Type string `json:"type"`
 85 | 	// Add specific fields as needed
 86 | }
 87 |
 88 | func (m *MixtureOfNConfig) GetType() string {
 89 | 	return m.Type
 90 | }
 91 |
 92 | // ChainOfThoughtConfig represents a chain-of-thought variant
 93 | type ChainOfThoughtConfig struct {
 94 | 	Type string `json:"type"`
 95 | 	// Add specific fields as needed
 96 | }
 97 |
 98 | func (c *ChainOfThoughtConfig) GetType() string {
 99 | 	return c.Type
100 | }
101 |
102 | // OpenAISFTConfig represents OpenAI SFT optimization configuration
103 | type OpenAISFTConfig struct {
104 | 	Model                  string   `json:"model"`
105 | 	BatchSize              *int     `json:"batch_size,omitempty"`
106 | 	LearningRateMultiplier *float64 `json:"learning_rate_multiplier,omitempty"`
107 | 	NEpochs                *int     `json:"n_epochs,omitempty"`
108 | 	Credentials            *string  `json:"credentials,omitempty"`
109 | 	APIBase                *string  `json:"api_base,omitempty"`
110 | 	Seed                   *int     `json:"seed,omitempty"`
111 | 	Suffix                 *string  `json:"suffix,omitempty"`
112 | }
113 |
114 | func (o *OpenAISFTConfig) GetType() string {
115 | 	return "openai_sft"
116 | }
117 |
118 | // FireworksSFTConfig represents Fireworks SFT optimization configuration
119 | type FireworksSFTConfig struct {
120 | 	Model       string  `json:"model"`
121 | 	Credentials *string `json:"credentials,omitempty"`
122 | 	AccountID   string  `json:"account_id"`
123 | 	APIBase     *string `json:"api_base,omitempty"`
124 | }
125 |
126 | func (f *FireworksSFTConfig) GetType() string {
127 | 	return "fireworks_sft"
128 | }
129 |
130 | // GCPVertexGeminiSFTConfig represents GCP Vertex Gemini SFT optimization configuration
131 | type GCPVertexGeminiSFTConfig struct {
132 | 	Model                    string   `json:"model"`
133 | 	BucketName               string   `json:"bucket_name"`
134 | 	ProjectID                string   `json:"project_id"`
135 | 	Region                   string   `json:"region"`
136 | 	LearningRateMultiplier   *float64 `json:"learning_rate_multiplier,omitempty"`
137 | 	AdapterSize              *int     `json:"adapter_size,omitempty"`
138 | 	NEpochs                  *int     `json:"n_epochs,omitempty"`
139 | 	ExportLastCheckpointOnly *bool    `json:"export_last_checkpoint_only,omitempty"`
140 | 	Credentials              *string  `json:"credentials,omitempty"`
141 | 	APIBase                  *string  `json:"api_base,omitempty"`
142 | 	Seed                     *int     `json:"seed,omitempty"`
143 | 	ServiceAccount           *string  `json:"service_account,omitempty"`
144 | 	KMSKeyName               *string  `json:"kms_key_name,omitempty"`
145 | 	TunedModelDisplayName    *string  `json:"tuned_model_display_name,omitempty"`
146 | 	BucketPathPrefix         *string  `json:"bucket_path_prefix,omitempty"`
147 | }
148 |
149 | func (g *GCPVertexGeminiSFTConfig) GetType() string {
150 | 	return "gcp_vertex_gemini_sft"
151 | }
152 |
153 | // OptimizationJobStatus represents optimization job status
154 | type OptimizationJobStatus string
155 |
156 | const (
157 | 	OptimizationJobStatusPending   OptimizationJobStatus = "pending"
158 | 	OptimizationJobStatusCompleted OptimizationJobStatus = "completed"
159 | 	OptimizationJobStatusFailed    OptimizationJobStatus = "failed"
160 | )
161 |
162 | // OptimizationJobInfo represents optimization job information
163 | type OptimizationJobInfo struct {
164 | 	Message         string                `json:"message"`
165 | 	Status          OptimizationJobStatus `json:"status"`
166 | 	Output          interface{}           `json:"output,omitempty"`
167 | 	EstimatedFinish *int64                `json:"estimated_finish,omitempty"`
168 | }
""",
    "types/filters.go": """
  1 | package types
  2 |
  3 | // FloatMetricFilter filters inferences by float metric values
  4 | type FloatMetricFilter struct {
  5 | 	MetricName         string  `json:"metric_name"`
  6 | 	Value              float64 `json:"value"`
  7 | 	ComparisonOperator string  `json:"comparison_operator"` // "<", "<=", "=", ">", ">=", "!="
  8 | 	Type               string  `json:"type"`
  9 | }
 10 |
 11 | func (f *FloatMetricFilter) GetType() string { return f.Type }
 12 |
 13 | // NewFloatMetricFilter creates a new float metric filter
 14 | func NewFloatMetricFilter(metricName string, value float64, operator string) *FloatMetricFilter {
 15 | 	return &FloatMetricFilter{
 16 | 		MetricName:         metricName,
 17 | 		Value:              value,
 18 | 		ComparisonOperator: operator,
 19 | 		Type:               "float_metric",
 20 | 	}
 21 | }
 22 |
 23 | // BooleanMetricFilter filters inferences by boolean metric values
 24 | type BooleanMetricFilter struct {
 25 | 	MetricName string `json:"metric_name"`
 26 | 	Value      bool   `json:"value"`
 27 | 	Type       string `json:"type"`
 28 | }
 29 |
 30 | func (f *BooleanMetricFilter) GetType() string { return f.Type }
 31 |
 32 | // NewBooleanMetricFilter creates a new boolean metric filter
 33 | func NewBooleanMetricFilter(metricName string, value bool) *BooleanMetricFilter {
 34 | 	return &BooleanMetricFilter{
 35 | 		MetricName: metricName,
 36 | 		Value:      value,
 37 | 		Type:       "boolean_metric",
 38 | 	}
 39 | }
 40 |
 41 | // TagFilter filters inferences by tag values
 42 | type TagFilter struct {
 43 | 	Key                string `json:"key"`
 44 | 	Value              string `json:"value"`
 45 | 	ComparisonOperator string `json:"comparison_operator"` // "=", "!="
 46 | 	Type               string `json:"type"`
 47 | }
 48 |
 49 | func (f *TagFilter) GetType() string { return f.Type }
 50 |
 51 | // NewTagFilter creates a new tag filter
 52 | func NewTagFilter(key, value, operator string) *TagFilter {
 53 | 	return &TagFilter{
 54 | 		Key:                key,
 55 | 		Value:              value,
 56 | 		ComparisonOperator: operator,
 57 | 		Type:               "tag",
 58 | 	}
 59 | }
 60 |
 61 | // TimeFilter filters inferences by timestamp
 62 | type TimeFilter struct {
 63 | 	Time               string `json:"time"`                // RFC 3339 timestamp
 64 | 	ComparisonOperator string `json:"comparison_operator"` // "<", "<=", "=", ">", ">=", "!="
 65 | 	Type               string `json:"type"`
 66 | }
 67 |
 68 | func (f *TimeFilter) GetType() string { return f.Type }
 69 |
 70 | // NewTimeFilter creates a new time filter
 71 | func NewTimeFilter(time, operator string) *TimeFilter {
 72 | 	return &TimeFilter{
 73 | 		Time:               time,
 74 | 		ComparisonOperator: operator,
 75 | 		Type:               "time",
 76 | 	}
 77 | }
 78 |
 79 | // AndFilter combines multiple filters with AND logic
 80 | type AndFilter struct {
 81 | 	Children []InferenceFilterTreeNode `json:"children"`
 82 | 	Type     string                    `json:"type"`
 83 | }
 84 |
 85 | func (f *AndFilter) GetType() string { return f.Type }
 86 |
 87 | // NewAndFilter creates a new AND filter
 88 | func NewAndFilter(children ...InferenceFilterTreeNode) *AndFilter {
 89 | 	return &AndFilter{
 90 | 		Children: children,
 91 | 		Type:     "and",
 92 | 	}
 93 | }
 94 |
 95 | // OrFilter combines multiple filters with OR logic
 96 | type OrFilter struct {
 97 | 	Children []InferenceFilterTreeNode `json:"children"`
 98 | 	Type     string                    `json:"type"`
 99 | }
100 |
101 | func (f *OrFilter) GetType() string { return f.Type }
102 |
103 | // NewOrFilter creates a new OR filter
104 | func NewOrFilter(children ...InferenceFilterTreeNode) *OrFilter {
105 | 	return &OrFilter{
106 | 		Children: children,
107 | 		Type:     "or",
108 | 	}
109 | }
110 |
111 | // NotFilter negates a filter
112 | type NotFilter struct {
113 | 	Child InferenceFilterTreeNode `json:"child"`
114 | 	Type  string                  `json:"type"`
115 | }
116 |
117 | func (f *NotFilter) GetType() string { return f.Type }
118 |
119 | // NewNotFilter creates a new NOT filter
120 | func NewNotFilter(child InferenceFilterTreeNode) *NotFilter {
121 | 	return &NotFilter{
122 | 		Child: child,
123 | 		Type:  "not",
124 | 	}
125 | }
""",
    "types/inference-request-options.go": """
  1 | package types
  2 |
  3 | import (
  4 | 	"github.com/denkhaus/tensorzero/util"
  5 | 	"github.com/google/uuid"
  6 | )
  7 |
  8 | type InferenceRequestOption func(*InferenceRequest)
  9 |
 10 | // WithDryRun sets the dry run option for the inference request
 11 | func WithDryRun(dryRun bool) InferenceRequestOption {
 12 | 	return func(g *InferenceRequest) {
 13 | 		g.Dryrun = util.BoolPtr(dryRun)
 14 | 	}
 15 | }
 16 |
 17 | // WithStream sets the stream option for the inference request
 18 | func WithStream(stream bool) InferenceRequestOption {
 19 | 	return func(g *InferenceRequest) {
 20 | 		g.Stream = util.BoolPtr(stream)
 21 | 	}
 22 | }
 23 |
 24 | // WithFunctionName sets the function name for the inference request
 25 | func WithFunctionName(functionName string) InferenceRequestOption {
 26 | 	return func(g *InferenceRequest) {
 27 | 		g.FunctionName = util.StringPtr(functionName)
 28 | 	}
 29 | }
 30 |
 31 | // WithModelName sets the model name for the inference request
 32 | func WithModelName(modelName string) InferenceRequestOption {
 33 | 	return func(g *InferenceRequest) {
 34 | 		g.ModelName = util.StringPtr(modelName)
 35 | 	}
 36 | }
 37 |
 38 | // WithEpisodeID sets the episode ID for the inference request
 39 | func WithEpisodeID(episodeID uuid.UUID) InferenceRequestOption {
 40 | 	return func(g *InferenceRequest) {
 41 | 		g.EpisodeID = &episodeID
 42 | 	}
 43 |
 44 | }
 45 |
 46 | // WithOutputSchema sets the output schema for the inference request
 47 | func WithOutputSchema(outputSchema map[string]interface{}) InferenceRequestOption {
 48 | 	return func(g *InferenceRequest) {
 49 | 		g.OutputSchema = outputSchema
 50 | 	}
 51 | }
 52 |
 53 | // WithAllowedTools sets the allowed tools for the inference request
 54 | func WithAllowedTools(allowedTools []string) InferenceRequestOption {
 55 | 	return func(g *InferenceRequest) {
 56 | 		g.AllowedTools = allowedTools
 57 | 	}
 58 |
 59 | }
 60 |
 61 | // WithAdditionalTools sets the additional tools for the inference request
 62 | func WithAdditionalTools(additionalTools []map[string]interface{}) InferenceRequestOption {
 63 | 	return func(g *InferenceRequest) {
 64 | 		g.AdditionalTools = additionalTools
 65 | 	}
 66 | }
 67 |
 68 | // WithParams sets the params for the inference request
 69 | func WithParams(params map[string]interface{}) InferenceRequestOption {
 70 | 	return func(g *InferenceRequest) {
 71 | 		g.Params = params
 72 | 	}
 73 | }
 74 |
 75 | // WithVariantName sets the variant name for the inference request
 76 | func WithVariantName(variantName string) InferenceRequestOption {
 77 | 	return func(g *InferenceRequest) {
 78 | 		g.VariantName = util.StringPtr(variantName)
 79 | 	}
 80 | }
 81 |
 82 | // WithToolChoice sets the tool choice for the inference request
 83 | func WithToolChoice(toolChoice ToolChoice) InferenceRequestOption {
 84 | 	return func(g *InferenceRequest) {
 85 | 		g.ToolChoice = toolChoice
 86 | 	}
 87 | }
 88 |
 89 | // WithParallelToolCalls sets the parallel tool calls option for the inference request
 90 | func WithParallelToolCalls(parallelToolCalls bool) InferenceRequestOption {
 91 | 	return func(g *InferenceRequest) {
 92 | 		g.ParallelToolCalls = util.BoolPtr(parallelToolCalls)
 93 | 	}
 94 | }
 95 |
 96 | // WithInternal sets the internal option for the inference request
 97 | func WithInternal(internal bool) InferenceRequestOption {
 98 | 	return func(g *InferenceRequest) {
 99 | 		g.Internal = util.BoolPtr(internal)
100 | 	}
101 | }
102 |
103 | // WithTags sets the tags for the inference request
104 | func WithTags(tags map[string]string) InferenceRequestOption {
105 | 	return func(g *InferenceRequest) {
106 | 		g.Tags = tags
107 | 	}
108 | }
109 |
110 | // WithCredentials sets the credentials for the inference request
111 | func WithCredentials(credentials map[string]string) InferenceRequestOption {
112 | 	return func(g *InferenceRequest) {
113 | 		g.Credentials = credentials
114 | 	}
115 | }
116 |
117 | // WithCacheOptions sets the cache options for the inference request
118 | func WithCacheOptions(cacheOptions map[string]interface{}) InferenceRequestOption {
119 | 	return func(g *InferenceRequest) {
120 | 		g.CacheOptions = cacheOptions
121 | 	}
122 | }
123 |
124 | // WithExtraBody sets the extra body for the inference request
125 | func WithExtraBody(extraBody []ExtraBody) InferenceRequestOption {
126 | 	return func(g *InferenceRequest) {
127 | 		g.ExtraBody = extraBody
128 | 	}
129 | }
130 |
131 | // WithExtraHeaders sets the extra headers for the inference request
132 | func WithExtraHeaders(extraHeaders []map[string]interface{}) InferenceRequestOption {
133 | 	return func(g *InferenceRequest) {
134 | 		g.ExtraHeaders = extraHeaders
135 | 	}
136 | }
137 |
138 | // WithIncludeOriginalResponse sets the include original response option for the inference request
139 | func WithIncludeOriginalResponse(includeOriginalResponse bool) InferenceRequestOption {
140 | 	return func(g *InferenceRequest) {
141 | 		g.IncludeOriginalResponse = util.BoolPtr(includeOriginalResponse)
142 | 	}
143 | }
""",
    "types/interfaces.go": """
 1 | package types
 2 |
 3 | import (
 4 | 	"context"
 5 |
 6 | 	"github.com/google/uuid"
 7 | )
 8 |
 9 | // Gateway represents the base interface for TensorZero gateways
10 | type Gateway interface {
11 | 	Inference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error)
12 | 	InferenceStream(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error)
13 | 	Feedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error)
14 | 	DynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error)
15 | 	DynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error)
16 | 	BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error)
17 | 	DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error
18 | 	ListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error)
19 | 	ListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error)
20 | 	Close() error
21 | }
22 |
23 | // ContentBlock represents a piece of content in a message
24 | type ContentBlock interface {
25 | 	GetType() string
26 | 	ToMap() map[string]interface{}
27 | }
28 |
29 | // InferenceFilterTreeNode represents the base interface for filter nodes
30 | type InferenceFilterTreeNode interface {
31 | 	GetType() string
32 | }
33 |
34 | // System represents system content
35 | type System interface{}
36 |
37 | // ContentBlockChunk represents streaming content chunks
38 | type ContentBlockChunk interface {
39 | 	GetType() string
40 | 	GetID() string
41 | }
42 |
43 | // InferenceResponse represents either chat or JSON inference response
44 | type InferenceResponse interface {
45 | 	GetInferenceID() uuid.UUID
46 | 	GetEpisodeID() uuid.UUID
47 | 	GetVariantName() string
48 | 	GetUsage() Usage
49 | 	GetFinishReason() *FinishReason
50 | 	GetOriginalResponse() *string
51 | }
52 |
53 | // InferenceChunk represents either chat or JSON chunk
54 | type InferenceChunk interface {
55 | 	GetInferenceID() uuid.UUID
56 | 	GetEpisodeID() uuid.UUID
57 | 	GetVariantName() string
58 | }
59 |
60 | // DatapointInsert represents a datapoint for insertion
61 | type DatapointInsert interface {
62 | 	GetFunctionName() string
63 | }
66 | // FunctionConfig represents a function configuration
67 | type FunctionConfig interface {
68 | 	GetType() string
69 | 	GetVariants() VariantsConfig
70 | }
71 |
72 | type VariantConfig interface {
73 | 	GetType() string
74 | }
75 |
76 | // OptimizationConfig represents optimization configurations
77 | type OptimizationConfig interface {
78 | 	GetType() string
79 | }
80 |
81 | // OptimizationJobHandle represents an optimization job handle
82 | type OptimizationJobHandle interface {
83 | 	GetType() string
84 | 	GetJobID() string
85 | }
""",
    "types/request.go": """
  1 | package types
  2 |
  3 | import (
  4 | 	"github.com/google/uuid"
  5 | )
  6 |
  7 | // ExtraBody represents a custom field to be added to the inference request body
  8 | type ExtraBody struct {
  9 | }
 10 |
 11 | // InferenceInput represents input to an inference request
 12 | type InferenceInput struct {
 13 | 	Messages []Message `json:"messages,omitempty"`
 14 | 	System   System    `json:"system,omitempty"`
 15 | }
 16 |
 17 | // ChatDatapointInsert represents chat datapoint insertion
 18 | type ChatDatapointInsert struct {
 19 | 	FunctionName      string            `json:"function_name"`
 20 | 	Input             InferenceInput    `json:"input"`
 21 | 	Output            interface{}       `json:"output,omitempty"`
 22 | 	AllowedTools      []string          `json:"allowed_tools,omitempty"`
 23 | 	AdditionalTools   []interface{}     `json:"additional_tools,omitempty"`
 24 | 	ToolChoice        *string           `json:"tool_choice,omitempty"`
 25 | 	ParallelToolCalls *bool             `json:"parallel_tool_calls,omitempty"`
 26 | 	Tags              map[string]string `json:"tags,omitempty"`
 27 | }
 28 |
 29 | func (c *ChatDatapointInsert) GetFunctionName() string { return c.FunctionName }
 30 |
 31 | // JsonDatapointInsert represents JSON datapoint insertion
 32 | type JsonDatapointInsert struct {
 33 | 	FunctionName string            `json:"function_name"`
 34 | 	Input        InferenceInput    `json:"input"`
 35 | 	Output       interface{}       `json:"output,omitempty"`
 36 | 	OutputSchema interface{}       `json:"output_schema,omitempty"`
 37 | 	Tags         map[string]string `json:"tags,omitempty"`
 38 | }
 39 |
 40 | func (j *JsonDatapointInsert) GetFunctionName() string { return j.FunctionName }
 41 |
 42 | // ListInferencesRequest represents a request to list inferences
 43 | type ListInferencesRequest struct {
 44 | 	FunctionName *string                 `json:"function_name,omitempty"`
 45 | 	EpisodeID    *uuid.UUID              `json:"episode_id,omitempty"`
 46 | 	VariantName  *string                 `json:"variant_name,omitempty"`
 47 | 	Filter       InferenceFilterTreeNode `json:"filter,omitempty"`
 48 | 	OrderBy      *OrderBy                `json:"order_by,omitempty"`
 49 | 	Limit        *int                    `json:"limit,omitempty"`
 50 | 	Offset       *int                    `json:"offset,omitempty"`
 51 | }
 52 |
 53 | // InferenceRequest represents an inference request
 54 | type InferenceRequest struct {
 55 | 	Input                   InferenceInput           `json:"input"`
 56 | 	FunctionName            *string                  `json:"function_name,omitempty"`
 57 | 	ModelName               *string                  `json:"model_name,omitempty"`
 58 | 	EpisodeID               *uuid.UUID               `json:"episode_id,omitempty"`
 59 | 	Stream                  *bool                    `json:"stream,omitempty"`
 60 | 	Params                  map[string]interface{}   `json:"params,omitempty"`
 61 | 	VariantName             *string                  `json:"variant_name,omitempty"`
 62 | 	Dryrun                  *bool                    `json:"dryrun,omitempty"`
 63 | 	OutputSchema            map[string]interface{}   `json:"output_schema,omitempty"`
 64 | 	AllowedTools            []string                 `json:"allowed_tools,omitempty"`
 65 | 	AdditionalTools         []map[string]interface{} `json:"additional_tools,omitempty"`
 66 | 	ToolChoice              ToolChoice               `json:"tool_choice,omitempty"`
 67 | 	ParallelToolCalls       *bool                    `json:"parallel_tool_calls,omitempty"`
 68 | 	Internal                *bool                    `json:"internal,omitempty"`
 69 | 	Tags                    map[string]string        `json:"tags,omitempty"`
 70 | 	Credentials             map[string]string        `json:"credentials,omitempty"`
 71 | 	CacheOptions            map[string]interface{}   `json:"cache_options,omitempty"`
 72 | 	ExtraBody               []ExtraBody              `json:"extra_body,omitempty"`
 73 | 	ExtraHeaders            []map[string]interface{} `json:"extra_headers,omitempty"`
 74 | 	IncludeOriginalResponse *bool                    `json:"include_original_response,omitempty"`
 75 | }
 76 |
 77 | // StoredInference represents a stored inference from the list API
 77 | type StoredInference struct {
 78 | 	ID             uuid.UUID              `json:"id"`
 79 | 	EpisodeID      uuid.UUID              `json:"episode_id"`
 80 | 	FunctionName   string                 `json:"function_name"`
 81 | 	VariantName    string                 `json:"variant_name"`
 82 | 	Input          InferenceInput         `json:"input"`
 83 | 	Output         interface{}            `json:"output"`
 84 | 	ToolParams     *ToolParams            `json:"tool_params,omitempty"`
 85 | 	ProcessingTime *float64               `json:"processing_time,omitempty"`
 86 | 	Timestamp      string                 `json:"timestamp"` // RFC 3339
 87 | 	Tags           map[string]string      `json:"tags,omitempty"`
 88 | 	MetricValues   map[string]interface{} `json:"metric_values,omitempty"`
 89 | }
 90 |
 91 | // Datapoint represents a datapoint
 92 | type Datapoint struct {
 93 | 	ID           uuid.UUID      `json:"id"`
 94 | 	Input        InferenceInput `json:"input"`
 95 | 	Output       interface{}    `json:"output"`
 96 | 	DatasetName  string         `json:"dataset_name"`
 97 | 	FunctionName string         `json:"function_name"`
 98 | 	ToolParams   *ToolParams    `json:"tool_params,omitempty"`
 99 | 	OutputSchema interface{}    `json:"output_schema,omitempty"`
100 | 	IsCustom     bool           `json:"is_custom"`
101 | }
102 |
103 | // FeedbackRequest represents a feedback request
104 | type FeedbackRequest struct {
105 | 	MetricName  string            `json:"metric_name"`
106 | 	Value       interface{}       `json:"value"`
107 | 	InferenceID *uuid.UUID        `json:"inference_id,omitempty"`
108 | 	EpisodeID   *uuid.UUID        `json:"episode_id,omitempty"`
109 | 	Dryrun      *bool             `json:"dryrun,omitempty"`
110 | 	Internal    *bool             `json:"internal,omitempty"`
111 | 	Tags        map[string]string `json:"tags,omitempty"`
112 | }
113 |
114 | // DynamicEvaluationRunRequest represents a dynamic evaluation run request
115 | type DynamicEvaluationRunRequest struct {
116 | 	Variants    map[string]string `json:"variants"`
117 | 	Tags        map[string]string `json:"tags,omitempty"`
118 | 	ProjectName *string           `json:"project_name,omitempty"`
119 | 	DisplayName *string           `json:"display_name,omitempty"`
120 | }
121 |
122 | // DynamicEvaluationRunEpisodeRequest represents a dynamic evaluation run episode request
123 | type DynamicEvaluationRunEpisodeRequest struct {
124 | 	RunID         uuid.UUID         `json:"run_id"`
125 | 	TaskName      *string           `json:"task_name,omitempty"`
126 | 	DatapointName *string           `json:"datapoint_name,omitempty"`
127 | 	Tags          map[string]string `json:"tags,omitempty"`
128 | }
129 |
130 | // ListDatapointsRequest represents a list datapoints request
131 | type ListDatapointsRequest struct {
132 | 	DatasetName  string  `json:"dataset_name"`
133 | 	FunctionName *string `json:"function_name,omitempty"`
134 | 	Limit        *int    `json:"limit,omitempty"`
135 | 	Offset       *int    `json:"offset,omitempty"`
136 | }
""",
    "types/response.go": """
 1 | package types
 2 |
 3 | import (
 4 | 	"github.com/google/uuid"
 5 | )
 6 |
 7 | // JsonInferenceOutput represents JSON inference output
 8 | type JsonInferenceOutput struct {
 9 | 	Raw    *string                `json:"raw,omitempty"`
10 | 	Parsed map[string]interface{} `json:"parsed,omitempty"`
11 | }
12 |
13 | // ChatInferenceResponse represents a chat inference response
14 | type ChatInferenceResponse struct {
15 | 	InferenceID      uuid.UUID      `json:"inference_id"`
16 | 	EpisodeID        uuid.UUID      `json:"episode_id"`
17 | 	VariantName      string         `json:"variant_name"`
18 | 	Content          []ContentBlock `json:"content"`
19 | 	Usage            Usage          `json:"usage"`
20 | 	FinishReason     *FinishReason  `json:"finish_reason,omitempty"`
21 | 	OriginalResponse *string        `json:"original_response,omitempty"`
22 | }
23 |
24 | // JsonInferenceResponse represents a JSON inference response
25 | type JsonInferenceResponse struct {
26 | 	InferenceID      uuid.UUID           `json:"inference_id"`
27 | 	EpisodeID        uuid.UUID           `json:"episode_id"`
28 | 	VariantName      string              `json:"variant_name"`
29 | 	Output           JsonInferenceOutput `json:"output"`
30 | 	Usage            Usage               `json:"usage"`
31 | 	FinishReason     *FinishReason       `json:"finish_reason,omitempty"`
32 | 	OriginalResponse *string             `json:"original_response,omitempty"`
33 | }
34 |
35 | func (c *ChatInferenceResponse) GetInferenceID() uuid.UUID      { return c.InferenceID }
36 | func (c *ChatInferenceResponse) GetEpisodeID() uuid.UUID        { return c.EpisodeID }
37 | func (c *ChatInferenceResponse) GetVariantName() string         { return c.VariantName }
38 | func (c *ChatInferenceResponse) GetUsage() Usage                { return c.Usage }
39 | func (c *ChatInferenceResponse) GetFinishReason() *FinishReason { return c.FinishReason }
40 | func (c *ChatInferenceResponse) GetOriginalResponse() *string   { return c.OriginalResponse }
41 |
42 | func (j *JsonInferenceResponse) GetInferenceID() uuid.UUID      { return j.InferenceID }
43 | func (j *JsonInferenceResponse) GetEpisodeID() uuid.UUID        { return j.EpisodeID }
44 | func (j *JsonInferenceResponse) GetVariantName() string         { return j.VariantName }
45 | func (j *JsonInferenceResponse) GetUsage() Usage                { return j.Usage }
46 | func (j *JsonInferenceResponse) GetFinishReason() *FinishReason { return j.FinishReason }
47 | func (j *JsonInferenceResponse) GetOriginalResponse() *string   { return j.OriginalResponse }
48 |
49 | // ChatChunk represents streaming chat chunk
50 | type ChatChunk struct {
51 | 	InferenceID  uuid.UUID           `json:"inference_id"`
52 | 	EpisodeID    uuid.UUID           `json:"episode_id"`
53 | 	VariantName  string              `json:"variant_name"`
54 | 	Content      []ContentBlockChunk `json:"content"`
55 | 	Usage        *Usage              `json:"usage,omitempty"`
56 | 	FinishReason *FinishReason       `json:"finish_reason,omitempty"`
57 | }
58 |
59 | // JsonChunk represents streaming JSON chunk
60 | type JsonChunk struct {
61 | 	InferenceID  uuid.UUID     `json:"inference_id"`
62 | 	EpisodeID    uuid.UUID     `json:"episode_id"`
63 | 	VariantName  string        `json:"variant_name"`
64 | 	Raw          string        `json:"raw"`
65 | 	Usage        *Usage        `json:"usage,omitempty"`
66 | 	FinishReason *FinishReason `json:"finish_reason,omitempty"`
67 | }
68 |
69 | func (c *ChatChunk) GetInferenceID() uuid.UUID { return c.InferenceID }
70 | func (c *ChatChunk) GetEpisodeID() uuid.UUID   { return c.EpisodeID }
71 | func (c *ChatChunk) GetVariantName() string    { return c.VariantName }
72 |
73 | func (j *JsonChunk) GetInferenceID() uuid.UUID { return j.InferenceID }
74 | func (j *JsonChunk) GetEpisodeID() uuid.UUID   { return j.EpisodeID }
75 | func (j *JsonChunk) GetVariantName() string    { return j.VariantName }
76 |
77 | // FeedbackResponse represents feedback response
78 | type FeedbackResponse struct {
79 | 	FeedbackID uuid.UUID `json:"feedback_id"`
80 | }
81 |
 82 | // DynamicEvaluationRunResponse represents dynamic evaluation run response
 83 | type DynamicEvaluationRunResponse struct {
 84 | 	RunID uuid.UUID `json:"run_id"`
 85 | }
 86 |
 87 | // DynamicEvaluationRunEpisodeResponse represents dynamic evaluation run episode response
 88 | type DynamicEvaluationRunEpisodeResponse struct {
 89 | 	EpisodeID uuid.UUID `json:"episode_id"`
 90 | }
""",
    "types/shared.go": """
  1 | package types
  2 |
  3 | import (
  4 | 	"fmt"
  5 | )
  6 |
  7 | // Usage represents token usage information
  8 | type Usage struct {
  9 | 	InputTokens  int `json:"input_tokens"`
 10 | 	OutputTokens int `json:"output_tokens"`
 11 | }
 12 |
 13 | // FinishReason represents the reason why inference finished
 14 | type FinishReason string
 15 |
 16 | const (
 17 | 	FinishReasonStop          FinishReason = "stop"
 18 | 	FinishReasonLength        FinishReason = "length"
 19 | 	FinishReasonToolCall      FinishReason = "tool_call"
 20 | 	FinishReasonContentFilter FinishReason = "content_filter"
 21 | 	FinishReasonUnknown       FinishReason = "unknown"
 22 | )
 23 |
 24 | // Text represents text content
 25 | type Text struct {
 26 | 	Text      *string     `json:"text,omitempty"`
 27 | 	Arguments interface{} `json:"arguments,omitempty"`
 28 | 	Type      string      `json:"type"`
 29 | }
 30 |
 31 | func NewText(text string) *Text {
 32 | 	return &Text{
 33 | 		Text: &text,
 34 | 		Type: "text",
 35 | 	}
 36 | }
 37 |
 38 | func NewTextWithArguments(arguments interface{}) *Text {
 39 | 	return &Text{
 40 | 		Arguments: arguments,
 41 | 		Type:      "text",
 42 | 	}
 43 | }
 44 |
 45 | func (t *Text) GetType() string {
 46 | 	return t.Type
 47 | }
 48 |
 49 | func (t *Text) ToMap() map[string]interface{} {
 50 | 	result := map[string]interface{}{
 51 | 		"type": t.Type,
 52 | 	}
 53 | 	if t.Text != nil {
 54 | 		result["text"] = *t.Text
 55 | 	}
 56 | 	if t.Arguments != nil {
 57 | 		result["arguments"] = t.Arguments
 58 | 	}
 59 | 	return result
 60 | }
 61 |
 62 | // RawText represents raw text content
 63 | type RawText struct {
 64 | 	Value string `json:"value"`
 65 | 	Type  string `json:"type"`
 66 | }
 67 |
 68 | func NewRawText(value string) *RawText {
 69 | 	return &RawText{
 70 | 		Value: value,
 71 | 		Type:  "raw_text",
 72 | 	}
 73 | }
 74 |
 75 | func (rt *RawText) GetType() string {
 76 | 	return rt.Type
 77 | }
 78 |
 79 | func (rt *RawText) ToMap() map[string]interface{} {
 80 | 	return map[string]interface{}{
 81 | 		"type":  rt.Type,
 82 | 		"value": rt.Value,
 83 | 	}
 84 | }
 85 |
 86 | // ImageBase64 represents base64-encoded image content
 87 | type ImageBase64 struct {
 88 | 	Data     string `json:"data"`
 89 | 	MimeType string `json:"mime_type"`
 90 | 	Type     string `json:"type"`
 91 | }
 92 |
 93 | func NewImageBase64(data, mimeType string) *ImageBase64 {
 94 | 	return &ImageBase64{
 95 | 		Data:     data,
 96 | 		MimeType: mimeType,
 97 | 		Type:     "image",
 98 | 	}
 99 | }
100 |
101 | func (img *ImageBase64) GetType() string {
102 | 	return img.Type
103 | }
104 |
105 | func (img *ImageBase64) ToMap() map[string]interface{} {
106 | 	return map[string]interface{}{
107 | 		"type":      img.Type,
108 | 		"data":      img.Data,
109 | 		"mime_type": img.MimeType,
110 | 	}
111 | }
112 |
113 | // ImageURL represents image content from URL
114 | type ImageURL struct {
115 | 	URL      string  `json:"url"`
116 | 	MimeType *string `json:"mime_type,omitempty"`
117 | 	Type     string  `json:"type"`
118 | }
119 |
120 | func NewImageURL(url string) *ImageURL {
121 | 	return &ImageURL{
122 | 		URL:  url,
123 | 		Type: "image",
124 | 	}
125 | }
126 |
127 | func NewImageURLWithMimeType(url, mimeType string) *ImageURL {
128 | 	return &ImageURL{
129 | 		URL:      url,
130 | 		MimeType: &mimeType,
131 | 		Type:     "image",
132 | 	}
133 | }
134 |
135 | func (img *ImageURL) GetType() string {
136 | 	return img.Type
137 | }
138 |
139 | func (img *ImageURL) ToMap() map[string]interface{} {
140 | 	result := map[string]interface{}{
141 | 		"type": img.Type,
142 | 		"url":  img.URL,
143 | 	}
144 | 	if img.MimeType != nil {
145 | 		result["mime_type"] = *img.MimeType
146 | 	}
147 | 	return result
148 | }
149 |
150 | // FileBase64 represents base64-encoded file content
151 | type FileBase64 struct {
152 | 	Data     string `json:"data"`
153 | 	MimeType string `json:"mime_type"`
154 | 	Type     string `json:"type"`
155 | }
156 |
157 | func NewFileBase64(data, mimeType string) *FileBase64 {
158 | 	return &FileBase64{
159 | 		Data:     data,
160 | 		MimeType: mimeType,
161 | 		Type:     "file",
162 | 	}
163 | }
164 |
165 | func (f *FileBase64) GetType() string {
166 | 	return f.Type
167 | }
168 |
169 | func (f *FileBase64) ToMap() map[string]interface{} {
170 | 	return map[string]interface{}{
171 | 		"type":      f.Type,
172 | 		"data":      f.Data,
173 | 		"mime_type": f.MimeType,
174 | 	}
175 | }
176 |
177 | // FileURL represents file content from URL
178 | type FileURL struct {
179 | 	URL  string `json:"url"`
180 | 	Type string `json:"type"`
181 | }
182 |
183 | func NewFileURL(url string) *FileURL {
184 | 	return &FileURL{
185 | 		URL:  url,
186 | 		Type: "file",
187 | 	}
188 | }
189 |
190 | func (f *FileURL) GetType() string {
191 | 	return f.Type
192 | }
193 |
194 | func (f *FileURL) ToMap() map[string]interface{} {
195 | 	return map[string]interface{}{
196 | 		"type": f.Type,
197 | 		"url":  f.URL,
198 | 	}
199 | }
200 |
201 | // ToolCall represents a tool call
202 | type ToolCall struct {
203 | 	ID           string                 `json:"id"`
204 | 	RawArguments string                 `json:"raw_arguments"`
205 | 	RawName      string                 `json:"raw_name"`
206 | 	Arguments    map[string]interface{} `json:"arguments,omitempty"`
207 | 	Name         *string                `json:"name,omitempty"`
208 | 	Type         string                 `json:"type"`
209 | }
210 |
211 | func NewToolCall(id, rawArguments, rawName string) *ToolCall {
212 | 	return &ToolCall{
213 | 		ID:           id,
214 | 		RawArguments: rawArguments,
215 | 		RawName:      rawName,
216 | 		Type:         "tool_call",
217 | 	}
218 | }
219 |
220 | func (tc *ToolCall) GetType() string {
221 | 	return tc.Type
222 | }
223 |
224 | func (tc *ToolCall) ToMap() map[string]interface{} {
225 | 	result := map[string]interface{}{
226 | 		"type":          tc.Type,
227 | 		"id":            tc.ID,
228 | 		"raw_arguments": tc.RawArguments,
229 | 		"raw_name":      tc.RawName,
230 | 	}
231 | 	if tc.Arguments != nil {
232 | 		result["arguments"] = tc.Arguments
233 | 	}
234 | 	if tc.Name != nil {
235 | 		result["name"] = *tc.Name
236 | 	}
237 | 	return result
238 | }
239 |
240 | // Thought represents a thought content block
241 | type Thought struct {
242 | 	Text      *string `json:"text,omitempty"`
243 | 	Type      string  `json:"type"`
244 | 	Signature *string `json:"signature,omitempty"`
245 | }
246 |
247 | func NewThought(text string) *Thought {
248 | 	return &Thought{
249 | 		Text: &text,
250 | 		Type: "thought",
251 | 	}
252 | }
253 |
254 | func (t *Thought) GetType() string {
255 | 	return t.Type
256 | }
257 |
258 | func (t *Thought) ToMap() map[string]interface{} {
259 | 	result := map[string]interface{}{
260 | 		"type": t.Type,
261 | 	}
262 | 	if t.Text != nil {
263 | 		result["text"] = *t.Text
264 | 	}
265 | 	if t.Signature != nil {
266 | 		result["signature"] = *t.Signature
267 | 	}
268 | 	return result
269 | }
270 |
271 | // ToolResult represents a tool result
272 | type ToolResult struct {
273 | 	Name   string `json:"name"`
274 | 	Result string `json:"result"`
275 | 	ID     string `json:"id"`
276 | 	Type   string `json:"type"`
277 | }
278 |
279 | func NewToolResult(name, result, id string) *ToolResult {
280 | 	return &ToolResult{
281 | 		Name:   name,
282 | 		Result: result,
283 | 		ID:     id,
284 | 		Type:   "tool_result",
285 | 	}
286 | }
287 |
288 | func (tr *ToolResult) GetType() string {
289 | 	return tr.Type
290 | }
291 |
292 | func (tr *ToolResult) ToMap() map[string]interface{} {
293 | 	return map[string]interface{}{
294 | 		"type":   tr.Type,
295 | 		"name":   tr.Name,
296 | 		"result": tr.Result,
297 | 		"id":     tr.ID,
298 | 	}
299 | }
300 |
301 | // UnknownContentBlock represents unknown content
302 | type UnknownContentBlock struct {
303 | 	Data              interface{} `json:"data"`
304 | 	ModelProviderName *string     `json:"model_provider_name,omitempty"`
305 | 	Type              string      `json:"type"`
306 | }
307 |
308 | func NewUnknownContentBlock(data interface{}) *UnknownContentBlock {
309 | 	return &UnknownContentBlock{
310 | 		Data: data,
311 | 		Type: "unknown",
312 | 	}
313 | }
314 |
315 | func (ucb *UnknownContentBlock) GetType() string {
316 | 	return ucb.Type
317 | }
318 |
319 | func (ucb *UnknownContentBlock) ToMap() map[string]interface{} {
320 | 	result := map[string]interface{}{
321 | 		"type": ucb.Type,
322 | 		"data": ucb.Data,
323 | 	}
324 | 	if ucb.ModelProviderName != nil {
325 | 		result["model_provider_name"] = *ucb.ModelProviderName
326 | 	}
327 | 	return result
328 | }
329 |
330 | // Message represents a message in a conversation
331 | type Message struct {
332 | 	Role    string         `json:"role"` // "user" or "assistant"
333 | 	Content []ContentBlock `json:"content"`
334 | }
335 |
336 | // TextChunk represents streaming text chunk
337 | type TextChunk struct {
338 | 	ID   string `json:"id"`
339 | 	Text string `json:"text"`
340 | 	Type string `json:"type"`
341 | }
342 |
343 | func (tc *TextChunk) GetType() string { return tc.Type }
344 | func (tc *TextChunk) GetID() string   { return tc.ID }
345 |
346 | // ToolCallChunk represents streaming tool call chunk
347 | type ToolCallChunk struct {
348 | 	ID           string `json:"id"`
349 | 	RawArguments string `json:"raw_arguments"`
350 | 	RawName      string `json:"raw_name"`
351 | 	Type         string `json:"type"`
352 | }
353 |
354 | func (tcc *ToolCallChunk) GetType() string { return tcc.Type }
355 | func (tcc *ToolCallChunk) GetID() string   { return tcc.ID }
356 |
357 | // ThoughtChunk represents streaming thought chunk
358 | type ThoughtChunk struct {
359 | 	ID        string  `json:"id"`
360 | 	Text      string  `json:"text"`
361 | 	Type      string  `json:"type"`
362 | 	Signature *string `json:"signature,omitempty"`
363 | }
364 |
365 | func (tc *ThoughtChunk) GetType() string { return tc.Type }
366 | func (tc *ThoughtChunk) GetID() string   { return tc.ID }
367 |
368 | // Tool represents a tool definition
369 | type Tool struct {
370 | 	Description string      `json:"description"`
371 | 	Parameters  interface{} `json:"parameters"`
372 | 	Name        string      `json:"name"`
373 | 	Strict      bool        `json:"strict"`
374 | }
375 |
376 | // ToolChoice represents tool choice options
377 | type ToolChoice interface{}
378 |
379 | // ToolParams represents tool parameters
380 | type ToolParams struct {
381 | 	ToolsAvailable    []Tool `json:"tools_available"`
382 | 	ToolChoice        string `json:"tool_choice"`
383 | 	ParallelToolCalls *bool  `json:"parallel_tool_calls,omitempty"`
384 | }
385 |
386 | // ExtraBody represents extra body content for requests
387 |
388 | // VariantExtraBody represents variant-specific extra body
389 | type VariantExtraBody struct {
390 | 	VariantName string      `json:"variant_name"`
391 | 	Pointer     string      `json:"pointer"`
392 | 	Value       interface{} `json:"value,omitempty"`
393 | 	Delete      *bool       `json:"delete,omitempty"`
394 | }
395 |
396 | // ProviderExtraBody represents provider-specific extra body
397 | type ProviderExtraBody struct {
398 | 	ModelProviderName string      `json:"model_provider_name"`
399 | 	Pointer           string      `json:"pointer"`
400 | 	Value             interface{} `json:"value,omitempty"`
401 | 	Delete            *bool       `json:"delete,omitempty"`
402 | }
403 |
404 | // TensorZeroError represents an error from TensorZero
405 | type TensorZeroError struct {
406 | 	StatusCode int
407 | 	Text       string
408 | }
409 |
410 | func (e *TensorZeroError) Error() string {
411 | 	return fmt.Sprintf("TensorZeroError (status code %d): %s", e.StatusCode, e.Text)
412 | }
413 |
414 | // TensorZeroInternalError represents an internal error
415 | type TensorZeroInternalError struct {
416 | 	Message string
417 | }
418 |
419 | func (e *TensorZeroInternalError) Error() string {
420 | 	return e.Message
421 | }
422 |
423 | // OrderBy specifies ordering for list inferences
424 | type OrderBy struct {
425 | 	By        string  `json:"by"`             // "timestamp" or "metric"
426 | 	Name      *string `json:"name,omitempty"` // metric name if by="metric"
427 | 	Direction string  `json:"direction"`      // "ASC" or "DESC"
428 | }
429 |
430 | // NewOrderByTimestamp creates ordering by timestamp
431 | func NewOrderByTimestamp(direction string) *OrderBy {
432 | 	return &OrderBy{
433 | 		By:        "timestamp",
434 | 		Direction: direction,
435 | 	}
436 | }
437 |
438 | // NewOrderByMetric creates ordering by metric
439 | func NewOrderByMetric(metricName, direction string) *OrderBy {
440 | 	return &OrderBy{
441 | 		By:        "metric",
442 | 		Name:      &metricName,
443 | 		Direction: direction,
444 | 	}
445 | }
"""
}

def extract_go_types_and_comments(file_content):
    types = []
    lines = file_content.split('\n')
    i = 0
    while i < len(lines):
        line = lines[i].strip()

        # Extract comments
        comments = []
        while line.startswith("//"):
            comments.append(line[2:].strip())
            i += 1
            if i < len(lines):
                line = lines[i].strip()
            else:
                break

        # Check for type definition
        match = re.match(r"type\s+(\w+)\s+(struct|interface|string|map|func)\s*({|\(.+\)|.*)", line)
        if match:
            type_name = match.group(1)
            type_kind = match.group(2)

            # For struct or interface, capture the full definition block
            if type_kind in ["struct", "interface"]:
                start_line = i
                brace_count = 0

                # Find the opening brace
                if '{' in line:
                    brace_count += 1

                # Find the closing brace
                while brace_count > 0 and i + 1 < len(lines):
                    i += 1
                    line = lines[i].strip()
                    brace_count += line.count('{')
                    brace_count -= line.count('}')

                full_definition = "\n".join(lines[start_line:i+1])
            else:
                # For other types (string, map, func), the definition is on one line
                full_definition = line

            types.append({
                "name": type_name,
                "kind": type_kind,
                "comments": "\n".join(comments),
                "definition": full_definition.strip()
            })
        i += 1
    return types

all_extracted_types = {}

for filename, content in go_files.items():
    extracted_types = extract_go_types_and_comments(content)
    for type_info in extracted_types:
        all_extracted_types[type_info["name"]] = type_info

# Print extracted types for verification
# The output will be a JSON string to avoid shell interpretation issues
print(json.dumps(all_extracted_types))
'] Result:Command executed in terminal  within working directory '/home/denkhaus/dev/gomodules/tensorzero-api'. Exit code: 0
Output:
{"Config": {"name": "Config", "kind": "struct", "comments": "Config represents TensorZero configuration", "definition": "type Config struct {\n\tFunctions FunctionsConfig `json:\"functions\"}`\n}"}, "FunctionsConfig": {"name": "FunctionsConfig", "kind": "map", "comments": "FunctionsConfig represents function configurations", "definition": "type FunctionsConfig map[string]FunctionConfig"}, "ChatFunctionConfig": {"name": "ChatFunctionConfig", "kind": "struct", "comments": "ChatFunctionConfig represents a chat function configuration", "definition": "type ChatFunctionConfig struct {\n\tType            string                 `json:\"type\"}`\n\tVariants        VariantsConfig         `json:\"variants\"}`\n\tSystemSchema    map[string]interface{} `json:\"system_schema,omitempty\"}`\n\tUserSchema      map[string]interface{} `json:\"user_schema,omitempty\"}`\n\tAssistantSchema map[string]interface{} `json:\"assistant_schema,omitempty\"}`\n}"}, "JsonFunctionConfig": {"name": "JsonFunctionConfig", "kind": "struct", "comments": "JsonFunctionConfig represents a JSON function configuration", "definition": "type JsonFunctionConfig struct {\n\tType            string                 `json:\"type\"}`\n\tVariants        VariantsConfig         `json:\"variants\"}`\n\tSystemSchema    map[string]interface{} `json:\"system_schema,omitempty\"}`\n\tUserSchema      map[string]interface{} `json:\"user_schema,omitempty\"}`\n\tAssistantSchema map[string]interface{} `json:\"assistant_schema,omitempty\"}`\n\tOutputSchema    map[string]interface{} `json:\"output_schema,omitempty\"}`\n}"}, "VariantsConfig": {"name": "VariantsConfig", "kind": "map", "comments": "VariantsConfig represents variant configurations", "definition": "type VariantsConfig map[string]VariantConfig"}, "ChatCompletionConfig": {"name": "ChatCompletionConfig", "kind": "struct", "comments": "ChatCompletionConfig represents a chat completion variant", "definition": "type ChatCompletionConfig struct {\n\tType              string  `json:\"type\"}`\n\tSystemTemplate    *string `json:\"system_template,omitempty\"}`\n\tUserTemplate      *string `json:\"user_template,omitempty\"}`\n\tAssistantTemplate *string `json:\"assistant_template,omitempty\"}`\n\tModel             string  `json:\"model\"}`\n}"}, "BestOfNSamplingConfig": {"name": "BestOfNSamplingConfig", "kind": "struct", "comments": "BestOfNSamplingConfig represents a best-of-n sampling variant", "definition": "type BestOfNSamplingConfig struct {\n\tType string `json:\"type\"}`\n\t// Add specific fields as needed\n}"}, "DiclConfig": {"name": "DiclConfig", "kind": "struct", "comments": "DiclConfig represents a DICL variant", "definition": "type DiclConfig struct {\n\tType string `json:\"type\"}`\n\t// Add specific fields as needed\n}"}, "MixtureOfNConfig": {"name": "MixtureOfNConfig", "kind": "struct", "comments": "MixtureOfNConfig represents a mixture-of-n variant", "definition": "type MixtureOfNConfig struct {\n\tType string `json:\"type\"}`\n\t// Add specific fields as needed\n}"}, "ChainOfThoughtConfig": {"name": "ChainOfThoughtConfig", "kind": "struct", "comments": "ChainOfThoughtConfig represents a chain-of-thought variant", "definition": "type ChainOfThoughtConfig struct {\n\tType string `json:\"type\"}`\n\t// Add specific fields as needed\n}"}, "OpenAISFTConfig": {"name": "OpenAISFTConfig", "kind": "struct", "comments": "OpenAISFTConfig represents OpenAI SFT optimization configuration", "definition": "type OpenAISFTConfig struct {\n\tModel                  string   `json:\"model\"}`\n\tBatchSize              *int     `json:\"batch_size,omitempty\"}`\n\tLearningRateMultiplier *float64 `json:\"learning_rate_multiplier,omitempty\"}`\n\tNEpochs                *int     `json:\"n_epochs,omitempty\"}`\n\tCredentials            *string  `json:\"credentials,omitempty\"}`\n\tAPIBase                *string  `json:\"api_base,omitempty\"}`\n\tSeed                   *int     `json:\"seed,omitempty\"}`\n\tSuffix                 *string  `json:\"suffix,omitempty\"}`\n}"}, "FireworksSFTConfig": {"name": "FireworksSFTConfig", "kind": "struct", "comments": "FireworksSFTConfig represents Fireworks SFT optimization configuration", "definition": "type FireworksSFTConfig struct {\n\tModel       string  `json:\"model\"}`\n\tCredentials *string `json:\"credentials,omitempty\"}`\n\tAccountID   string  `json:\"account_id\"}`\n\tAPIBase     *string `json:\"api_base,omitempty\"}`\n}"}, "GCPVertexGeminiSFTConfig": {"name": "GCPVertexGeminiSFTConfig", "kind": "struct", "comments": "GCPVertexGeminiSFTConfig represents GCP Vertex Gemini SFT optimization configuration", "definition": "type GCPVertexGeminiSFTConfig struct {\n\tModel                    string   `json:\"model\"}`\n\tBucketName               string   `json:\"bucket_name\"}`\n\tProjectID                string   `json:\"project_id\"}`\n\tRegion                   string   `json:\"region\"}`\n\tLearningRateMultiplier   *float64 `json:\"learning_rate_multiplier,omitempty\"}`\n\tAdapterSize              *int     `json:\"adapter_size,omitempty\"}`\n\tNEpochs                  *int     `json:\"n_epochs,omitempty\"}`\n\tExportLastCheckpointOnly *bool    `json:\"export_last_checkpoint_only,omitempty\"}`\n\tCredentials              *string  `json:\"credentials,omitempty\"}`\n\tAPIBase                  *string  `json:\"api_base,omitempty\"}`\n\tSeed                     *int     `json:\"seed,omitempty\"}`\n\tServiceAccount           *string  `json:\"service_account,omitempty\"}`\n\tKMSKeyName               *string  `json:\"kms_key_name,omitempty\"}`\n\tTunedModelDisplayName    *string  `json:\"tuned_model_display_name,omitempty\"}`\n\tBucketPathPrefix         *string  `json:\"bucket_path_prefix,omitempty\"}`\n}"}, "OptimizationJobStatus": {"name": "OptimizationJobStatus", "kind": "string", "comments": "OptimizationJobStatus represents optimization job status", "definition": "type OptimizationJobStatus string"}, "OptimizationJobInfo": {"name": "OptimizationJobInfo", "kind": "struct", "comments": "OptimizationJobInfo represents optimization job information", "definition": "type OptimizationJobInfo struct {\n\tMessage         string                `json:\"message\"}`\n\tStatus          OptimizationJobStatus `json:\"status\"}`\n\tOutput          interface{}           `json:\"output,omitempty\"}`\n\tEstimatedFinish *int64                `json:\"estimated_finish,omitempty\"}`\n}"}, "FloatMetricFilter": {"name": "FloatMetricFilter", "kind": "struct", "comments": "FloatMetricFilter filters inferences by float metric values", "definition": "type FloatMetricFilter struct {\n\tMetricName         string  `json:\"metric_name\"}`\n\tValue              float64 `json:\"value\"}`\n\tComparisonOperator string  `json:\"comparison_operator\"}` // \"<\", \"<=\", \"=\", \">\", \">=\", \"!=\"\n\tType               string  `json:\"type\"}`\n}"}, "BooleanMetricFilter": {"name": "BooleanMetricFilter", "kind": "struct", "comments": "BooleanMetricFilter filters inferences by boolean metric values", "definition": "type BooleanMetricFilter struct {\n\tMetricName string `json:\"metric_name\"}`\n\tValue      bool   `json:\"value\"}`\n\tType       string `json:\"type\"}`\n}"}, "TagFilter": {"name": "TagFilter", "kind": "struct", "comments": "TagFilter filters inferences by tag values", "definition": "type TagFilter struct {\n\tKey                string `json:\"key\"}`\n\tValue              string `json:\"value\"}`\n\tComparisonOperator string `json:\"comparison_operator\"}` // \"=\", \"!=\"\n\tType               string `json:\"type\"}`\n}"}, "TimeFilter": {"name": "TimeFilter", "kind": "struct", "comments": "TimeFilter filters inferences by timestamp", "definition": "type TimeFilter struct {\n\tTime               string `json:\"time\"}`                // RFC 3339 timestamp\n\tComparisonOperator string `json:\"comparison_operator\"}` // \"<\", \"<=\", \"=\", \">\", \">=\", \"!=\"\n\tType               string  `json:\"type\"}`\n}"}, "AndFilter": {"name": "AndFilter", "kind": "struct", "comments": "AndFilter combines multiple filters with AND logic", "definition": "type AndFilter struct {\n\tChildren []InferenceFilterTreeNode `json:\"children\"}`\n\tType     string                    `json:\"type\"}`\n}"}, "OrFilter": {"name": "OrFilter", "kind": "struct", "comments": "OrFilter combines multiple filters with OR logic", "definition": "type OrFilter struct {\n\tChildren []InferenceFilterTreeNode `json:\"children\"}`\n\tType     string                    `json:\"type\"}`\n}"}, "NotFilter": {"name": "NotFilter", "kind": "struct", "comments": "NotFilter negates a filter", "definition": "type NotFilter struct {\n\tChild InferenceFilterTreeNode `json:\"child\"}`\n\tType  string                  `json:\"type\"}`\n}"}, "InferenceRequestOption": {"name": "InferenceRequestOption", "kind": "func", "comments": "", "definition": "type InferenceRequestOption func(*InferenceRequest)"}, "ExtraBody": {"name": "ExtraBody", "kind": "struct", "comments": "ExtraBody represents a custom field to be added to the inference request body", "definition": "type ExtraBody struct {\n}"}, "InferenceInput": {"name": "InferenceInput", "kind": "struct", "comments": "InferenceInput represents input to an inference request", "definition": "type InferenceInput struct {\n\tMessages []Message `json:\"messages,omitempty\"}`\n\tSystem   System    `json:\"system,omitempty\"}`\n}"}, "ChatDatapointInsert": {"name": "ChatDatapointInsert", "kind": "struct", "comments": "ChatDatapointInsert represents chat datapoint insertion", "definition": "type ChatDatapointInsert struct {\n\tFunctionName      string            `json:\"function_name\"}`\n\tInput             InferenceInput    `json:\"input\"}`\n\tOutput            interface{}       `json:\"output,omitempty\"}`\n\tAllowedTools      []string          `json:\"allowed_tools,omitempty\"}`\n\tAdditionalTools   []interface{}     `json:\"additional_tools,omitempty\"}`\n\tToolChoice        *string           `json:\"tool_choice,omitempty\"}`\n\tParallelToolCalls *bool             `json:\"parallel_tool_calls,omitempty\"}`\n\tTags              map[string]string `json:\"tags,omitempty\"}`\n}"}, "JsonDatapointInsert": {"name": "JsonDatapointInsert", "kind": "struct", "comments": "JsonDatapointInsert represents JSON datapoint insertion", "definition": "type JsonDatapointInsert struct {\n\tFunctionName string            `json:\"function_name\"}`\n\tInput        InferenceInput    `json:\"input\"}`\n\tOutput       interface{}       `json:\"output,omitempty\"}`\n\tOutputSchema interface{}       `json:\"output_schema,omitempty\"}`\n\tTags         map[string]string `json:\"tags,omitempty\"}`\n}"}, "ListInferencesRequest": {"name": "ListInferencesRequest", "kind": "struct", "comments": "ListInferencesRequest represents a request to list inferences", "definition": "type ListInferencesRequest struct {\n\tFunctionName *string                 `json:\"function_name,omitempty\"}`\n\tEpisodeID    *uuid.UUID              `json:\"episode_id,omitempty\"}`\n\tVariantName  *string                 `json:\"variant_name,omitempty\"}`\n\tFilter       InferenceFilterTreeNode `json:\"filter,omitempty\"}`\n\tOrderBy      *OrderBy                `json:\"order_by,omitempty\"}`\n\tLimit        *int                    `json:\"limit,omitempty\"}`\n\tOffset       *int                    `json:\"offset,omitempty\"}`\n}"}, "InferenceRequest": {"name": "InferenceRequest", "kind": "struct", "comments": "InferenceRequest represents an inference request", "definition": "type InferenceRequest struct {\n\tInput                   InferenceInput           `json:\"input\"}`\n\tFunctionName            *string                  `json:\"function_name,omitempty\"}`\n\tModelName               *string                  `json:\"model_name,omitempty\"}`\n\tEpisodeID               *uuid.UUID               `json:\"episode_id,omitempty\"}`\n\tStream                  *bool                    `json:\"stream,omitempty\"}`\n\tParams                  map[string]interface{}   `json:\"params,omitempty\"}`\n\tVariantName             *string                  `json:\"variant_name,omitempty\"}`\n\tDryrun                  *bool                    `json:\"dryrun,omitempty\"}`\n\tOutputSchema            map[string]interface{}   `json:\"output_schema,omitempty\"}`\n\tAllowedTools            []string                 `json:\"allowed_tools,omitempty\"}`\n\tAdditionalTools         []map[string]interface{} `json:\"additional_tools,omitempty\"}`\n\tToolChoice              ToolChoice               `json:\"tool_choice,omitempty\"}`\n\tParallelToolCalls       *bool                    `json:\"parallel_tool_calls,omitempty\"}`\n\tInternal                *bool                    `json:\"internal,omitempty\"}`\n\tTags                    map[string]string        `json:\"tags,omitempty\"}`\n\tCredentials             map[string]string        `json:\"credentials,omitempty\"}`\n\tCacheOptions            map[string]interface{}   `json:\"cache_options,omitempty\"}`\n\tExtraBody               []ExtraBody              `json:\"extra_body,omitempty\"}`\n\tExtraHeaders            []map[string]interface{} `json:\"extra_headers,omitempty\"}`\n\tIncludeOriginalResponse *bool                    `json:\"include_original_response,omitempty\"}`\n}"}, "StoredInference": {"name": "StoredInference", "kind": "struct", "comments": "StoredInference represents a stored inference from the list API", "definition": "type StoredInference struct {\n\tID             uuid.UUID              `json:\"id\"}`\n\tEpisodeID      uuid.UUID              `json:\"episode_id\"}`\n\tFunctionName   string                 `json:\"function_name\"}`\n\tVariantName    string                 `json:\"variant_name\"}`\n\tInput          InferenceInput         `json:\"input\"}`\n\tOutput         interface{}            `json:\"output\"}`\n\tToolParams     *ToolParams            `json:\"tool_params,omitempty\"}`\n\tProcessingTime *float64               `json:\"processing_time,omitempty\"}`\n\tTimestamp      string                 `json:\"timestamp\"}` // RFC 3339\n\tTags           map[string]string      `json:\"tags,omitempty\"}`\n\tMetricValues   map[string]interface{} `json:\"metric_values,omitempty\"}`\n}"}, "Datapoint": {"name": "Datapoint", "kind": "struct", "comments": "Datapoint represents a datapoint", "definition": "type Datapoint struct {\n\tID           uuid.UUID      `json:\"id\"}`\n\tInput        InferenceInput `json:\"input\"}`\n\tOutput       interface{}    `json:\"output\"}`\n\tDatasetName  string         `json:\"dataset_name\"}`\n\tFunctionName string         `json:\"function_name\"}`\n\tToolParams   *ToolParams    `json:\"tool_params,omitempty\"}`\n\tOutputSchema interface{}    `json:\"output_schema,omitempty\"}`\n\tIsCustom     bool           `json:\"is_custom\"}`\n}"}, "FeedbackRequest": {"name": "FeedbackRequest", "kind": "struct", "comments": "FeedbackRequest represents a feedback request", "definition": "type FeedbackRequest struct {\n\tMetricName  string            `json:\"metric_name\"}`\n\tValue       interface{}       `json:\"value\"}`\n\tInferenceID *uuid.UUID        `json:\"inference_id,omitempty\"}`\n\tEpisodeID   *uuid.UUID        `json:\"episode_id,omitempty\"}`\n\tDryrun      *bool             `json:\"dryrun,omitempty\"}`\n\tInternal    *bool             `json:\"internal,omitempty\"}`\n\tTags        map[string]string `json:\"tags,omitempty\"}`\n}"}, "DynamicEvaluationRunRequest": {"name": "DynamicEvaluationRunRequest", "kind": "struct", "comments": "DynamicEvaluationRunRequest represents a dynamic evaluation run request", "definition": "type DynamicEvaluationRunRequest struct {\n\tVariants    map[string]string `json:\"variants\"}`\n\tTags        map[string]string `json:\"tags,omitempty\"}`\n\tProjectName *string           `json:\"project_name,omitempty\"}`\n\tDisplayName *string           `json:\"display_name,omitempty\"}`\n}"}, "DynamicEvaluationRunEpisodeRequest": {"name": "DynamicEvaluationRunEpisodeRequest", "kind": "struct", "comments": "DynamicEvaluationRunEpisodeRequest represents a dynamic evaluation run episode request", "definition": "type DynamicEvaluationRunEpisodeRequest {\n\tRunID         uuid.UUID         `json:\"run_id\"}`\n\tTaskName      *string           `json:\"task_name,omitempty\"}`\n\tDatapointName *string           `json:\"datapoint_name,omitempty\"}`\n\tTags          map[string]string `json:\"tags,omitempty\"}`\n}"}, "ListDatapointsRequest": {"name": "ListDatapointsRequest", "kind": "struct", "comments": "ListDatapointsRequest represents a list datapoints request", "definition": "type ListDatapointsRequest struct {\n\tDatasetName  string  `json:\"dataset_name\"}`\n\tFunctionName *string `json:\"function_name,omitempty\"}`\n\tLimit        *int    `json:\"limit,omitempty\"}`\n\tOffset       *int    `json:\"offset,omitempty\"}`\n}"}, "Gateway": {"name": "Gateway", "kind": "interface", "comments": "Gateway represents the base interface for TensorZero gateways", "definition": "type Gateway interface {\n\tInference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error)\n\tInferenceStream(ctx context.Context, req *InferenceRequest) (<--chan InferenceChunk, <--chan error)\n\tFeedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error)\n\tDynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error)\n\tDynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error)\n\tBulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error)\n\tDeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error\n\tListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error)\n\tListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error)\n\tClose() error\n}"}, "ContentBlock": {"name": "ContentBlock", "kind": "interface", "comments": "ContentBlock represents a piece of content in a message", "definition": "type ContentBlock interface {\n\tGetType() string\n\tToMap() map[string]interface{}\n}"}, "InferenceFilterTreeNode": {"name": "InferenceFilterTreeNode", "kind": "interface", "comments": "InferenceFilterTreeNode represents the base interface for filter nodes", "definition": "type InferenceFilterTreeNode interface {\n\tGetType() string\n}"}, "System": {"name": "System", "kind": "interface", "comments": "System represents system content", "definition": "type System interface{}"}, "ContentBlockChunk": {"name": "ContentBlockChunk", "kind": "interface", "comments": "ContentBlockChunk represents streaming content chunks", "definition": "type ContentBlockChunk interface {\n\tGetType() string\n\tGetID() string\n}"}, "InferenceResponse": {"name": "InferenceResponse", "kind": "interface", "comments": "InferenceResponse represents either chat or JSON inference response", "definition": "type InferenceResponse interface {\n\tGetInferenceID() uuid.UUID\n\tGetEpisodeID() uuid.UUID\n\tGetVariantName() string\n\tGetUsage() Usage\n\tGetFinishReason() *FinishReason\n\tGetOriginalResponse() *string\n}"}, "InferenceChunk": {"name": "InferenceChunk", "kind": "interface", "comments": "InferenceChunk represents either chat or JSON chunk", "definition": "type InferenceChunk interface {\n\tGetInferenceID() uuid.UUID\n\tGetEpisodeID() uuid.UUID\n\tGetVariantName() string\n}"}, "DatapointInsert": {"name": "DatapointInsert", "kind": "interface", "comments": "DatapointInsert represents a datapoint for insertion", "definition": "type DatapointInsert interface {\n\tGetFunctionName() string\n}"}, "FunctionConfig": {"name": "FunctionConfig", "kind": "interface", "comments": "FunctionConfig represents a function configuration", "definition": "type FunctionConfig interface {\n\tGetType() string\n\tGetVariants() VariantsConfig\n}"}, "VariantConfig": {"name": "VariantConfig", "kind": "interface", "comments": "", "definition": "type VariantConfig interface {\n\tGetType() string\n}"}, "OptimizationConfig": {"name": "OptimizationConfig", "kind": "interface", "comments": "OptimizationConfig represents optimization configurations", "definition": "type OptimizationConfig interface {\n\tGetType() string\n}"}, "OptimizationJobHandle": {"name": "OptimizationJobHandle", "kind": "interface", "comments": "OptimizationJobHandle represents an optimization job handle", "definition": "type OptimizationJobHandle interface {\n\tGetType() string\n\tGetJobID() string\n}"}, "JsonInferenceOutput": {"name": "JsonInferenceOutput", "kind": "struct", "comments": "JsonInferenceOutput represents JSON inference output", "definition": "type JsonInferenceOutput struct {\n\tRaw    *string                `json:\"raw,omitempty\"}`\n\tParsed map[string]interface{} `json:\"parsed,omitempty\"}`\n}"}, "ChatInferenceResponse": {"name": "ChatInferenceResponse", "kind": "struct", "comments": "ChatInferenceResponse represents a chat inference response", "definition": "type ChatInferenceResponse struct {\n\tInferenceID      uuid.UUID      `json:\"inference_id\"}`\n\tEpisodeID        uuid.UUID      `json:\"episode_id\"}`\n\tVariantName      string         `json:\"variant_name\"}`\n\tContent          []ContentBlock `json:\"content\"}`\n\tUsage            Usage          `json:\"usage\"}`\n\tFinishReason     *FinishReason  `json:\"finish_reason,omitempty\"}`\n\tOriginalResponse *string        `json:\"original_response,omitempty\"}`\n}"}, "JsonInferenceResponse": {"name": "JsonInferenceResponse", "kind": "struct", "comments": "JsonInferenceResponse represents a JSON inference response", "definition": "type JsonInferenceResponse struct {\n\tInferenceID      uuid.UUID           `json:\"inference_id\"}`\n\tEpisodeID        uuid.UUID           `json:\"episode_id\"}`\n\tVariantName      string              `json:\"variant_name\"}`\n\tOutput           JsonInferenceOutput `json:\"output\"}`\n\tUsage            Usage               `json:\"usage\"}`\n\tFinishReason     *FinishReason       `json:\"finish_reason,omitempty\"}`\n\tOriginalResponse *string             `json:\"original_response,omitempty\"}`\n}"}, "ChatChunk": {"name": "ChatChunk", "kind": "struct", "comments": "ChatChunk represents streaming chat chunk", "definition": "type ChatChunk struct {\n\tInferenceID  uuid.UUID           `json:\"inference_id\"}`\n\tEpisodeID    uuid.UUID           `json:\"episode_id\"}`\n\tVariantName  string              `json:\"variant_name\"}`\n\tContent      []ContentBlockChunk `json:\"content\"}`\n\tUsage        *Usage              `json:\"usage,omitempty\"}`\n\tFinishReason *FinishReason       `json:\"finish_reason,omitempty\"}`\n}"}, "JsonChunk": {"name": "JsonChunk", "kind": "struct", "comments": "JsonChunk represents streaming JSON chunk", "definition": "type JsonChunk struct {\n\tInferenceID  uuid.UUID     `json:\"inference_id\"}`\n\tEpisodeID    uuid.UUID     `json:\"episode_id\"}`\n\tVariantName  string        `json:\"variant_name\"}`\n\tRaw          string        `json:\"raw\"}`\n\tUsage        *Usage        `json:\"usage,omitempty\"}`\n\tFinishReason *FinishReason `json:\"finish_reason,omitempty\"}`\n}"}, "FeedbackResponse": {"name": "FeedbackResponse", "kind": "struct", "comments": "FeedbackResponse represents feedback response", "definition": "type FeedbackResponse struct {\n\tFeedbackID uuid.UUID `json:\"feedback_id\"}`\n}"}, "DynamicEvaluationRunResponse": {"name": "DynamicEvaluationRunResponse", "kind": "struct", "comments": "DynamicEvaluationRunResponse represents dynamic evaluation run response", "definition": "type DynamicEvaluationRunResponse struct {\n\tRunID uuid.UUID `json:\"run_id\"}`\n}"}, "DynamicEvaluationRunEpisodeResponse": {"name": "DynamicEvaluationRunEpisodeResponse", "kind": "struct", "comments": "DynamicEvaluationRunEpisodeResponse represents dynamic evaluation run episode response", "definition": "type DynamicEvaluationRunEpisodeResponse struct {\n\tEpisodeID uuid.UUID `json:\"episode_id\"}`\n}"}, "Usage": {"name": "Usage", "kind": "struct", "comments": "Usage represents token usage information", "definition": "type Usage struct {\n\tInputTokens  int `json:\"input_tokens\"}`\n\tOutputTokens int `json:\"output_tokens\"}`\n}"}, "FinishReason": {"name": "FinishReason", "kind": "string", "comments": "FinishReason represents the reason why inference finished", "definition": "type FinishReason string"}, "Text": {"name": "Text", "kind": "struct", "comments": "Text represents text content", "definition": "type Text struct {\n\tText      *string     `json:\"text,omitempty\"}`\n\tArguments interface{} `json:\"arguments,omitempty\"}`\n\tType      string      `json:\"type\"}`\n}"}, "RawText": {"name": "RawText", "kind": "struct", "comments": "RawText represents raw text content", "definition": "type RawText struct {\n\tValue string `json:\"value\"}`\n\tType  string `json:\"type\"}`\n}"}, "ImageBase64": {"name": "ImageBase64", "kind": "struct", "comments": "ImageBase64 represents base64-encoded image content", "definition": "type ImageBase64 struct {\n\tData     string `json:\"data\"}`\n\tMimeType string `json:\"mime_type\"}`\n\tType     string `json:\"type\"}`\n}"}, "ImageURL": {"name": "ImageURL", "kind": "struct", "comments": "ImageURL represents image content from URL", "definition": "type ImageURL struct {\n\tURL      string  `json:\"url\"}`\n\tMimeType *string `json:\"mime_type,omitempty\"}`\n\tType     string  `json:\"type\"}`\n}"}, "FileBase64": {"name": "FileBase64", "kind": "struct", "comments": "FileBase64 represents base64-encoded file content", "definition": "type FileBase64 struct {\n\tData     string `json:\"data\"}`\n\tMimeType string `json:\"mime_type\"}`\n\tType     string `json:\"type\"}`\n}"}, "FileURL": {"name": "FileURL", "kind": "struct", "comments": "FileURL represents file content from URL", "definition": "type FileURL struct {\n\tURL  string `json:\"url\"}`\n\tType string `json:\"type\"}`\n}"}, "ToolCall": {"name": "ToolCall", "kind": "struct", "comments": "ToolCall represents a tool call", "definition": "type ToolCall struct {\n\tID           string                 `json:\"id\"}`\n\tRawArguments string                 `json:\"raw_arguments\"}`\n\tRawName      string                 `json:\"raw_name\"}`\n\tArguments    map[string]interface{} `json:\"arguments,omitempty\"}`\n\tName         *string                `json:\"name,omitempty\"}`\n\tType         string                 `json:\"type\"}`\n}"}, "Thought": {"name": "Thought", "kind": "struct", "comments": "Thought represents a thought content block", "definition": "type Thought struct {\n\tText      *string `json:\"text,omitempty\"}`\n\tType      string  `json:\"type\"}`\n\tSignature *string `json:\"signature,omitempty\"}`\n}"}, "ToolResult": {"name": "ToolResult", "kind": "struct", "comments": "ToolResult represents a tool result", "definition": "type ToolResult struct {\n\tName   string `json:\"name\"}`\n\tResult string `json:\"result\"}`\n\tID     string `json:\"id\"}`\n\tType   string `json:\"type\"}`\n}"}, "UnknownContentBlock": {"name": "UnknownContentBlock", "kind": "struct", "comments": "UnknownContentBlock represents unknown content", "definition": "type UnknownContentBlock struct {\n\tData              interface{} `json:\"data\"}`\n\tModelProviderName *string     `json:\"model_provider_name,omitempty\"}`\n\tType              string      `json:\"type\"}`\n}"}, "Message": {"name": "Message", "kind": "struct", "comments": "Message represents a message in a conversation", "definition": "type Message struct {\n\tRole    string         `json:\"role\"}` // \"user\" or \"assistant\"\n\tContent []ContentBlock `json:\"content\"}`\n}"}, "TextChunk": {"name": "TextChunk", "kind": "struct", "comments": "TextChunk represents streaming text chunk", "definition": "type TextChunk struct {\n\tID   string `json:\"id\"}`\n\tText string `json:\"text\"}`\n\tType string `json:\"type\"}`\n}"}, "ToolCallChunk": {"name": "ToolCallChunk", "kind": "struct", "comments": "ToolCallChunk represents streaming tool call chunk", "definition": "type ToolCallChunk struct {\n\tID           string `json:\"id\"}`\n\tRawArguments string `json:\"raw_arguments\"}`\n\tRawName      string `json:\"raw_name\"}`\n\tType         string `json:\"type\"}`\n}"}, "ThoughtChunk": {"name": "ThoughtChunk", "kind": "struct", "comments": "ThoughtChunk represents streaming thought chunk", "definition": "type ThoughtChunk struct {\n\tID        string  `json:\"id\"}`\n\tText      string  `json:\"text\"}`\n\tType      string  `json:\"type\"}`\n\tSignature *string `json:\"signature,omitempty\"}`\n}"}, "Tool": {"name": "Tool", "kind": "struct", "comments": "Tool represents a tool definition", "definition": "type Tool struct {\n\tDescription string      `json:\"description\"}`\n\tParameters  interface{} `json:\"parameters\"}`\n\tName        string      `json:\"name\"}`\n\tStrict      bool        `json:\"strict\"}`\n}"}, "ToolChoice": {"name": "ToolChoice", "kind": "interface", "comments": "ToolChoice represents tool choice options", "definition": "type ToolChoice interface{}"}, "ToolParams": {"name": "ToolParams", "kind": "struct", "comments": "ToolParams represents tool parameters", "definition": "type ToolParams struct {\n\tToolsAvailable    []Tool `json:\"tools_available\"}`\n\tToolChoice        string `json:\"tool_choice\"}`\n\tParallelToolCalls *bool  `json:\"parallel_tool_calls,omitempty\"}`\n}"}, "VariantExtraBody": {"name": "VariantExtraBody", "kind": "struct", "comments": "VariantExtraBody represents variant-specific extra body", "definition": "type VariantExtraBody struct {\n\tVariantName string      `json:\"variant_name\"}`\n\tPointer     string      `json:\"pointer\"}`\n\tValue       interface{} `json:\"value,omitempty\"}`\n\tDelete      *bool       `json:\"delete,omitempty\"}`\n}"}, "ProviderExtraBody": {"name": "ProviderExtraBody", "kind": "struct", "comments": "ProviderExtraBody represents provider-specific extra body", "definition": "type ProviderExtraBody struct {\n\tModelProviderName string      `json:\"model_provider_name\"}`\n\tPointer           string      `json:\"pointer\"}`\n\tValue             interface{} `json:\"value,omitempty\"}`\n\tDelete            *bool       `json:\"delete,omitempty\"}`\n}"}, "TensorZeroError": {"name": "TensorZeroError", "kind": "struct", "comments": "TensorZeroError represents an error from TensorZero", "definition": "type TensorZeroError struct {\n\tStatusCode int\n\tText       string\n}"}, "TensorZeroInternalError": {"name": "TensorZeroInternalError", "kind": "struct", "comments": "TensorZeroInternalError represents an internal error", "definition": "type TensorZeroInternalError struct {\n\tMessage string\n}"}, "OrderBy": {"name": "OrderBy", "kind": "struct", "comments": "OrderBy specifies ordering for list inferences", "definition": "type OrderBy struct {\n\tBy        string  `json:\"by\"}`             // \"timestamp\" or \"metric\"\n\tName      *string `json:\"name,omitempty\"}` // metric name if by=\"metric\"\n\tDirection string  `json:\"direction\"}`      // \"ASC\" or \"DESC\"\n}"}}
<environment_details>
# VSCode Visible Files


# VSCode Open Tabs
docs/api/inference.md

# Current Time
Current time in ISO 8601 UTC format: 2025-07-29T17:22:00.672Z
User time zone: Europe/Berlin, UTC+2:00

# Current Cost
$0.00

# Current Mode
<slug>code</slug>
<name>Code</name>
<model>gemini-2.5-flash</model>
====

REMINDERS

Below is your current list of reminders for this task. Keep them updated as you progress.

| # | Content | Status |
|---|---------|--------|
| 1 | List all Go type definitions from the `types` directory. | Completed |
| 2 | Analyze and establish a mapping between Go type names and `docs/api` markdown files. | In Progress |
| 3 | Extract Go type definitions and comments from files. | In Progress |
| 4 | Map extracted Go types to existing `docs/api` markdown files. | Pending |
| 5 | For each Go type, read its corresponding documentation from `docs/api`. | Pending |
| 6 | Generate new documentation files for each Go type. | Pending |


IMPORTANT: When task status changes, remember to call the `update_todo_list` tool to update your progress.

</environment_details>
