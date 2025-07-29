package tensorzero

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Gateway represents the base interface for TensorZero gateways
type Gateway interface {
	Inference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error)
	InferenceStream(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error)
	Feedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error)
	DynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error)
	DynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error)
	BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error)
	ListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error)
	Close() error
}

// HTTPGateway implements Gateway using HTTP requests
type HTTPGateway struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewHTTPGateway creates a new HTTP gateway client
func NewHTTPGateway(baseURL string, options ...HTTPGatewayOption) *HTTPGateway {
	gateway := &HTTPGateway{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{},
		timeout:    30 * time.Second,
	}

	for _, option := range options {
		option(gateway)
	}

	if gateway.httpClient.Timeout == 0 {
		gateway.httpClient.Timeout = gateway.timeout
	}

	return gateway
}

// HTTPGatewayOption represents configuration options for HTTPGateway
type HTTPGatewayOption func(*HTTPGateway)

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) HTTPGatewayOption {
	return func(g *HTTPGateway) {
		g.timeout = timeout
		g.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) HTTPGatewayOption {
	return func(g *HTTPGateway) {
		g.httpClient = client
	}
}

// InferenceRequest represents an inference request
type InferenceRequest struct {
	Input                   InferenceInput         `json:"input"`
	FunctionName            *string                `json:"function_name,omitempty"`
	ModelName               *string                `json:"model_name,omitempty"`
	EpisodeID               *uuid.UUID             `json:"episode_id,omitempty"`
	Stream                  *bool                  `json:"stream,omitempty"`
	Params                  map[string]interface{} `json:"params,omitempty"`
	VariantName             *string                `json:"variant_name,omitempty"`
	Dryrun                  *bool                  `json:"dryrun,omitempty"`
	OutputSchema            map[string]interface{} `json:"output_schema,omitempty"`
	AllowedTools            []string               `json:"allowed_tools,omitempty"`
	AdditionalTools         []map[string]interface{} `json:"additional_tools,omitempty"`
	ToolChoice              ToolChoice             `json:"tool_choice,omitempty"`
	ParallelToolCalls       *bool                  `json:"parallel_tool_calls,omitempty"`
	Internal                *bool                  `json:"internal,omitempty"`
	Tags                    map[string]string      `json:"tags,omitempty"`
	Credentials             map[string]string      `json:"credentials,omitempty"`
	CacheOptions            map[string]interface{} `json:"cache_options,omitempty"`
	ExtraBody               []ExtraBody            `json:"extra_body,omitempty"`
	ExtraHeaders            []map[string]interface{} `json:"extra_headers,omitempty"`
	IncludeOriginalResponse *bool                  `json:"include_original_response,omitempty"`
}

// FeedbackRequest represents a feedback request
type FeedbackRequest struct {
	MetricName   string            `json:"metric_name"`
	Value        interface{}       `json:"value"`
	InferenceID  *uuid.UUID        `json:"inference_id,omitempty"`
	EpisodeID    *uuid.UUID        `json:"episode_id,omitempty"`
	Dryrun       *bool             `json:"dryrun,omitempty"`
	Internal     *bool             `json:"internal,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// DynamicEvaluationRunRequest represents a dynamic evaluation run request
type DynamicEvaluationRunRequest struct {
	Variants    map[string]string `json:"variants"`
	Tags        map[string]string `json:"tags,omitempty"`
	ProjectName *string           `json:"project_name,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
}

// DynamicEvaluationRunEpisodeRequest represents a dynamic evaluation run episode request
type DynamicEvaluationRunEpisodeRequest struct {
	RunID         uuid.UUID         `json:"run_id"`
	TaskName      *string           `json:"task_name,omitempty"`
	DatapointName *string           `json:"datapoint_name,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// ListDatapointsRequest represents a list datapoints request
type ListDatapointsRequest struct {
	DatasetName  string  `json:"dataset_name"`
	FunctionName *string `json:"function_name,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Offset       *int    `json:"offset,omitempty"`
}

// DatapointInsert represents a datapoint for insertion
type DatapointInsert interface {
	GetFunctionName() string
}

func (c *ChatDatapointInsert) GetFunctionName() string { return c.FunctionName }
func (j *JsonDatapointInsert) GetFunctionName() string { return j.FunctionName }

// Datapoint represents a datapoint
type Datapoint struct {
	ID           uuid.UUID      `json:"id"`
	Input        InferenceInput `json:"input"`
	Output       interface{}    `json:"output"`
	DatasetName  string         `json:"dataset_name"`
	FunctionName string         `json:"function_name"`
	ToolParams   *ToolParams    `json:"tool_params,omitempty"`
	OutputSchema interface{}    `json:"output_schema,omitempty"`
	IsCustom     bool           `json:"is_custom"`
}

// Inference makes an inference request
func (g *HTTPGateway) Inference(ctx context.Context, req *InferenceRequest) (InferenceResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/inference", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return parseInferenceResponse(body)
}

// InferenceStream makes a streaming inference request
func (g *HTTPGateway) InferenceStream(ctx context.Context, req *InferenceRequest) (<-chan InferenceChunk, <-chan error) {
	chunkCh := make(chan InferenceChunk, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(chunkCh)
		defer close(errCh)

		// Set stream to true
		streamReq := *req
		streamTrue := true
		streamReq.Stream = &streamTrue

		reqBody, err := json.Marshal(streamReq)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/inference", bytes.NewBuffer(reqBody))
		if err != nil {
			errCh <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "text/event-stream")

		resp, err := g.httpClient.Do(httpReq)
		if err != nil {
			errCh <- fmt.Errorf("failed to make request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errCh <- &TensorZeroError{
				StatusCode: resp.StatusCode,
				Text:       string(body),
			}
			return
		}

		// Parse SSE stream
		scanner := NewSSEScanner(resp.Body)
		for scanner.Scan() {
			event := scanner.Event()
			if event.Data == "" {
				continue
			}

			chunk, err := parseInferenceChunk([]byte(event.Data))
			if err != nil {
				errCh <- fmt.Errorf("failed to parse chunk: %w", err)
				return
			}

			select {
			case chunkCh <- chunk:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("failed to read stream: %w", err)
		}
	}()

	return chunkCh, errCh
}

// Feedback sends feedback
func (g *HTTPGateway) Feedback(ctx context.Context, req *FeedbackRequest) (*FeedbackResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/feedback", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var feedbackResp FeedbackResponse
	if err := json.NewDecoder(resp.Body).Decode(&feedbackResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &feedbackResp, nil
}

// DynamicEvaluationRun creates a dynamic evaluation run
func (g *HTTPGateway) DynamicEvaluationRun(ctx context.Context, req *DynamicEvaluationRunRequest) (*DynamicEvaluationRunResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/dynamic_evaluation_run", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var evalResp DynamicEvaluationRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&evalResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &evalResp, nil
}

// DynamicEvaluationRunEpisode creates a dynamic evaluation run episode
func (g *HTTPGateway) DynamicEvaluationRunEpisode(ctx context.Context, req *DynamicEvaluationRunEpisodeRequest) (*DynamicEvaluationRunEpisodeResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/dynamic_evaluation_run_episode", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var episodeResp DynamicEvaluationRunEpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&episodeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &episodeResp, nil
}

// BulkInsertDatapoints inserts multiple datapoints
func (g *HTTPGateway) BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []DatapointInsert) ([]uuid.UUID, error) {
	reqBody, err := json.Marshal(map[string]interface{}{
		"datapoints": datapoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("/datasets/%s/datapoints/bulk", url.PathEscape(datasetName))
	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var ids []uuid.UUID
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return ids, nil
}

// DeleteDatapoint deletes a datapoint
func (g *HTTPGateway) DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error {
	endpoint := fmt.Sprintf("/datasets/%s/datapoints/%s", url.PathEscape(datasetName), datapointID.String())
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", g.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	return nil
}

// ListDatapoints lists datapoints
func (g *HTTPGateway) ListDatapoints(ctx context.Context, req *ListDatapointsRequest) ([]Datapoint, error) {
	endpoint := fmt.Sprintf("/datasets/%s/datapoints", url.PathEscape(req.DatasetName))
	
	u, err := url.Parse(g.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	if req.FunctionName != nil {
		q.Set("function_name", *req.FunctionName)
	}
	if req.Limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *req.Limit))
	}
	if req.Offset != nil {
		q.Set("offset", fmt.Sprintf("%d", *req.Offset))
	}
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var datapoints []Datapoint
	if err := json.NewDecoder(resp.Body).Decode(&datapoints); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return datapoints, nil
}

// Close closes the gateway
func (g *HTTPGateway) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}
// ListInferences lists stored inferences with filtering and ordering
func (g *HTTPGateway) ListInferences(ctx context.Context, req *ListInferencesRequest) ([]StoredInference, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/inferences/list", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var inferences []StoredInference
	if err := json.NewDecoder(resp.Body).Decode(&inferences); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return inferences, nil
}
