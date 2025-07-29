package types

import (
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
)

type InferenceRequestOption func(*InferenceRequest)

// WithDryRun sets the dry run option for the inference request
func WithDryRun(dryRun bool) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Dryrun = util.BoolPtr(dryRun)
	}
}

// WithStream sets the stream option for the inference request
func WithStream(stream bool) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Stream = util.BoolPtr(stream)
	}
}

// WithFunctionName sets the function name for the inference request
func WithFunctionName(functionName string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.FunctionName = util.StringPtr(functionName)
	}
}

// WithModelName sets the model name for the inference request
func WithModelName(modelName string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.ModelName = util.StringPtr(modelName)
	}
}

// WithEpisodeID sets the episode ID for the inference request
func WithEpisodeID(episodeID uuid.UUID) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.EpisodeID = &episodeID
	}

}

// WithOutputSchema sets the output schema for the inference request
func WithOutputSchema(outputSchema map[string]interface{}) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.OutputSchema = outputSchema
	}
}

// WithAllowedTools sets the allowed tools for the inference request
func WithAllowedTools(allowedTools []string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.AllowedTools = allowedTools
	}

}

// WithAdditionalTools sets the additional tools for the inference request
func WithAdditionalTools(additionalTools []map[string]interface{}) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.AdditionalTools = additionalTools
	}
}

// WithParams sets the params for the inference request
func WithParams(params map[string]interface{}) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Params = params
	}
}

// WithVariantName sets the variant name for the inference request
func WithVariantName(variantName string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.VariantName = util.StringPtr(variantName)
	}
}

// WithToolChoice sets the tool choice for the inference request
func WithToolChoice(toolChoice ToolChoice) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.ToolChoice = toolChoice
	}
}

// WithParallelToolCalls sets the parallel tool calls option for the inference request
func WithParallelToolCalls(parallelToolCalls bool) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.ParallelToolCalls = util.BoolPtr(parallelToolCalls)
	}
}

// WithInternal sets the internal option for the inference request
func WithInternal(internal bool) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Internal = util.BoolPtr(internal)
	}
}

// WithTags sets the tags for the inference request
func WithTags(tags map[string]string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Tags = tags
	}
}

// WithCredentials sets the credentials for the inference request
func WithCredentials(credentials map[string]string) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.Credentials = credentials
	}
}

// WithCacheOptions sets the cache options for the inference request
func WithCacheOptions(cacheOptions map[string]interface{}) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.CacheOptions = cacheOptions
	}
}

// WithExtraBody sets the extra body for the inference request
func WithExtraBody(extraBody []ExtraBody) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.ExtraBody = extraBody
	}
}

// WithExtraHeaders sets the extra headers for the inference request
func WithExtraHeaders(extraHeaders []map[string]interface{}) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.ExtraHeaders = extraHeaders
	}
}

// WithIncludeOriginalResponse sets the include original response option for the inference request
func WithIncludeOriginalResponse(includeOriginalResponse bool) InferenceRequestOption {
	return func(g *InferenceRequest) {
		g.IncludeOriginalResponse = util.BoolPtr(includeOriginalResponse)
	}
}
