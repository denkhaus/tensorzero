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

	"github.com/denkhaus/tensorzero/types"
	"github.com/google/uuid"
)

// Gateway represents the base interface for TensorZero gateways
type Gateway interface {
	Inference(ctx context.Context, req *types.InferenceRequest) (types.InferenceResponse, error)
	InferenceStream(ctx context.Context, req *types.InferenceRequest) (<-chan types.InferenceChunk, <-chan error)
	Feedback(ctx context.Context, req *types.FeedbackRequest) (*types.FeedbackResponse, error)
	DynamicEvaluationRun(ctx context.Context, req *types.DynamicEvaluationRunRequest) (*types.DynamicEvaluationRunResponse, error)
	DynamicEvaluationRunEpisode(ctx context.Context, req *types.DynamicEvaluationRunEpisodeRequest) (*types.DynamicEvaluationRunEpisodeResponse, error)
	BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []types.DatapointInsert) ([]uuid.UUID, error)
	DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error
	ListDatapoints(ctx context.Context, req *types.ListDatapointsRequest) ([]types.Datapoint, error)
	ListInferences(ctx context.Context, req *types.ListInferencesRequest) ([]types.StoredInference, error)
	Close() error
}

// HTTPGateway implements Gateway using HTTP requests
type httpGateway struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewHTTPGateway creates a new HTTP gateway client
func NewHTTPGateway(baseURL string, options ...HTTPGatewayOption) Gateway {
	gateway := &httpGateway{
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
type HTTPGatewayOption func(*httpGateway)

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) HTTPGatewayOption {
	return func(g *httpGateway) {
		g.timeout = timeout
		g.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) HTTPGatewayOption {
	return func(g *httpGateway) {
		g.httpClient = client
	}
}

// Inference makes an inference request
func (g *httpGateway) Inference(ctx context.Context, req *types.InferenceRequest) (types.InferenceResponse, error) {
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
		return nil, &types.TensorZeroError{
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
func (g *httpGateway) InferenceStream(ctx context.Context, req *types.InferenceRequest) (<-chan types.InferenceChunk, <-chan error) {
	chunkCh := make(chan types.InferenceChunk, 10)
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
			errCh <- &types.TensorZeroError{
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
func (g *httpGateway) Feedback(ctx context.Context, req *types.FeedbackRequest) (*types.FeedbackResponse, error) {
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
		return nil, &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var feedbackResp types.FeedbackResponse
	if err := json.NewDecoder(resp.Body).Decode(&feedbackResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &feedbackResp, nil
}

// DynamicEvaluationRun creates a dynamic evaluation run
func (g *httpGateway) DynamicEvaluationRun(ctx context.Context, req *types.DynamicEvaluationRunRequest) (*types.DynamicEvaluationRunResponse, error) {
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
		return nil, &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var evalResp types.DynamicEvaluationRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&evalResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &evalResp, nil
}

// DynamicEvaluationRunEpisode creates a dynamic evaluation run episode
func (g *httpGateway) DynamicEvaluationRunEpisode(ctx context.Context, req *types.DynamicEvaluationRunEpisodeRequest) (*types.DynamicEvaluationRunEpisodeResponse, error) {
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
		return nil, &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var episodeResp types.DynamicEvaluationRunEpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&episodeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &episodeResp, nil
}

// BulkInsertDatapoints inserts multiple datapoints
func (g *httpGateway) BulkInsertDatapoints(ctx context.Context, datasetName string, datapoints []types.DatapointInsert) ([]uuid.UUID, error) {
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
		return nil, &types.TensorZeroError{
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
func (g *httpGateway) DeleteDatapoint(ctx context.Context, datasetName string, datapointID uuid.UUID) error {
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
		return &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	return nil
}

// ListDatapoints lists datapoints
func (g *httpGateway) ListDatapoints(ctx context.Context, req *types.ListDatapointsRequest) ([]types.Datapoint, error) {
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
		return nil, &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var datapoints []types.Datapoint
	if err := json.NewDecoder(resp.Body).Decode(&datapoints); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return datapoints, nil
}

// Close closes the gateway
func (g *httpGateway) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}

// ListInferences lists stored inferences with filtering and ordering
func (g *httpGateway) ListInferences(ctx context.Context, req *types.ListInferencesRequest) ([]types.StoredInference, error) {
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
		return nil, &types.TensorZeroError{
			StatusCode: resp.StatusCode,
			Text:       string(body),
		}
	}

	var inferences []types.StoredInference
	if err := json.NewDecoder(resp.Body).Decode(&inferences); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return inferences, nil
}
