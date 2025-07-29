package types

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWithDryRun(t *testing.T) {
	req := &InferenceRequest{}
	WithDryRun(true)(req)
	assert.NotNil(t, req.Dryrun)
	assert.True(t, *req.Dryrun)

	WithDryRun(false)(req)
	assert.NotNil(t, req.Dryrun)
	assert.False(t, *req.Dryrun)
}

func TestWithStream(t *testing.T) {
	req := &InferenceRequest{}
	WithStream(true)(req)
	assert.NotNil(t, req.Stream)
	assert.True(t, *req.Stream)
}

func TestWithFunctionName(t *testing.T) {
	req := &InferenceRequest{}
	name := "test_function"
	WithFunctionName(name)(req)
	assert.NotNil(t, req.FunctionName)
	assert.Equal(t, name, *req.FunctionName)
}

func TestWithModelName(t *testing.T) {
	req := &InferenceRequest{}
	name := "gpt-3.5-turbo"
	WithModelName(name)(req)
	assert.NotNil(t, req.ModelName)
	assert.Equal(t, name, *req.ModelName)
}

func TestWithEpisodeID(t *testing.T) {
	req := &InferenceRequest{}
	id := uuid.New()
	WithEpisodeID(id)(req)
	assert.NotNil(t, req.EpisodeID)
	assert.Equal(t, id, *req.EpisodeID)
}

func TestWithOutputSchema(t *testing.T) {
	req := &InferenceRequest{}
	schema := map[string]interface{}{"type": "object"}
	WithOutputSchema(schema)(req)
	assert.NotNil(t, req.OutputSchema)
	assert.Equal(t, schema, req.OutputSchema)
}

func TestWithAllowedTools(t *testing.T) {
	req := &InferenceRequest{}
	tools := []string{"tool1", "tool2"}
	WithAllowedTools(tools)(req)
	assert.NotNil(t, req.AllowedTools)
	assert.Equal(t, tools, req.AllowedTools)
}

func TestWithAdditionalTools(t *testing.T) {
	req := &InferenceRequest{}
	tools := []map[string]interface{}{{"name": "tool3"}}
	WithAdditionalTools(tools)(req)
	assert.NotNil(t, req.AdditionalTools)
	assert.Equal(t, tools, req.AdditionalTools)
}

func TestWithParams(t *testing.T) {
	req := &InferenceRequest{}
	params := map[string]interface{}{"temp": 0.7}
	WithParams(params)(req)
	assert.NotNil(t, req.Params)
	assert.Equal(t, params, req.Params)
}

func TestWithVariantName(t *testing.T) {
	req := &InferenceRequest{}
	name := "variant_A"
	WithVariantName(name)(req)
	assert.NotNil(t, req.VariantName)
	assert.Equal(t, name, *req.VariantName)
}

func TestWithToolChoice(t *testing.T) {
	req := &InferenceRequest{}
	choice := ToolChoice("auto")
	WithToolChoice(choice)(req)
	assert.Equal(t, choice, req.ToolChoice)
}

func TestWithParallelToolCalls(t *testing.T) {
	req := &InferenceRequest{}
	WithParallelToolCalls(true)(req)
	assert.NotNil(t, req.ParallelToolCalls)
	assert.True(t, *req.ParallelToolCalls)
}

func TestWithInternal(t *testing.T) {
	req := &InferenceRequest{}
	WithInternal(true)(req)
	assert.NotNil(t, req.Internal)
	assert.True(t, *req.Internal)
}

func TestWithTags(t *testing.T) {
	req := &InferenceRequest{}
	tags := map[string]string{"env": "dev"}
	WithTags(tags)(req)
	assert.NotNil(t, req.Tags)
	assert.Equal(t, tags, req.Tags)
}

func TestWithCredentials(t *testing.T) {
	req := &InferenceRequest{}
	creds := map[string]string{"api_key": "123"}
	WithCredentials(creds)(req)
	assert.NotNil(t, req.Credentials)
	assert.Equal(t, creds, req.Credentials)
}

func TestWithCacheOptions(t *testing.T) {
	req := &InferenceRequest{}
	options := map[string]interface{}{"ttl": 3600}
	WithCacheOptions(options)(req)
	assert.NotNil(t, req.CacheOptions)
	assert.Equal(t, options, req.CacheOptions)
}

func TestWithExtraHeaders(t *testing.T) {
	req := &InferenceRequest{}
	headers := []map[string]interface{}{{"Auth": "Bearer token"}}
	WithExtraHeaders(headers)(req)
	assert.NotNil(t, req.ExtraHeaders)
	assert.Equal(t, headers, req.ExtraHeaders)
}

func TestWithIncludeOriginalResponse(t *testing.T) {
	req := &InferenceRequest{}
	WithIncludeOriginalResponse(true)(req)
	assert.NotNil(t, req.IncludeOriginalResponse)
	assert.True(t, *req.IncludeOriginalResponse)
}
