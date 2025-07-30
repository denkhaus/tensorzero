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

// WithMetricValue sets both metric name and value in a single option
func WithMetricValue(metricName string, value interface{}) FeedbackRequestOption {
    return func(req *Request) {
        req.MetricName = metricName
        req.Value = value
    }
}

// WithComment is a convenience method for comment feedback
func WithComment(comment string) FeedbackRequestOption {
    return WithMetricValue("comment", comment)
}

// WithRating is a convenience method for rating feedback
func WithRating(rating float64) FeedbackRequestOption {
    return WithMetricValue("rating", rating)
}

// WithBooleanMetric is a convenience method for boolean feedback
func WithBooleanMetric(metricName string, value bool) FeedbackRequestOption {
    return WithMetricValue(metricName, value)
}
