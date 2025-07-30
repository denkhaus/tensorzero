package feedback

import (
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
)

type FeedbackRequestOption func(*Request)

// WithDryRun sets the dry run option for the inference request
func WithDryRun(dryRun bool) FeedbackRequestOption {
	return func(g *Request) {
		g.Dryrun = util.BoolPtr(dryRun)
	}
}

// WithTags sets the tags for the feedback request
func WithTags(tags map[string]string) FeedbackRequestOption {
	return func(g *Request) {
		g.Tags = tags
	}
}

// WithInternal sets the internal flag for the feedback request
func WithInternal(internal bool) FeedbackRequestOption {
	return func(g *Request) {
		g.Internal = util.BoolPtr(internal)
	}
}

// WithInferenceID sets the inference ID for the feedback request
func WithInferenceID(inferenceID uuid.UUID) FeedbackRequestOption {
	return func(g *Request) {
		g.InferenceID = &inferenceID
	}
}

// WithEpisodeID sets the episode ID for the feedback request
func WithEpisodeID(episodeID uuid.UUID) FeedbackRequestOption {
	return func(g *Request) {
		g.EpisodeID = &episodeID
	}
}
