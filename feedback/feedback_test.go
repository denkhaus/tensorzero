//go:build unit

package feedback

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFeedbackRequest(t *testing.T) {
	inferenceID := uuid.New()
	episodeID := uuid.New()
	dryrun := true
	internal := false
	tags := map[string]string{"feedback_type": "quality"}

	req := Request{
		MetricName:  "relevance",
		Value:       0.9,
		InferenceID: &inferenceID,
		EpisodeID:   &episodeID,
		Dryrun:      &dryrun,
		Internal:    &internal,
		Tags:        tags,
	}

	assert.Equal(t, "relevance", req.MetricName)
	assert.Equal(t, 0.9, req.Value)
	assert.Equal(t, inferenceID, *req.InferenceID)
	assert.Equal(t, episodeID, *req.EpisodeID)
	assert.True(t, *req.Dryrun)
	assert.False(t, *req.Internal)
	assert.Equal(t, tags, req.Tags)
}

func TestWithInferenceID(t *testing.T) {
	req := Request{}
	inferenceID := uuid.New()
	option := WithInferenceID(inferenceID)
	option(&req)
	assert.NotNil(t, req.InferenceID)
	assert.Equal(t, inferenceID, *req.InferenceID)
}

func TestWithEpisodeID(t *testing.T) {
	req := Request{}
	episodeID := uuid.New()
	option := WithEpisodeID(episodeID)
	option(&req)
	assert.NotNil(t, req.EpisodeID)
	assert.Equal(t, episodeID, *req.EpisodeID)
}

func TestWithTags(t *testing.T) {
	req := Request{}
	tags := map[string]string{"key": "value"}
	option := WithTags(tags)
	option(&req)
	assert.NotNil(t, req.Tags)
	assert.Equal(t, tags, req.Tags)
}

func TestWithInternal(t *testing.T) {
	req := Request{}
	option := WithInternal(true)
	option(&req)
	assert.NotNil(t, req.Internal)
	assert.True(t, *req.Internal)

	option = WithInternal(false)
	option(&req)
	assert.NotNil(t, req.Internal)
	assert.False(t, *req.Internal)
}

func TestFeedbackResponse(t *testing.T) {
	feedbackID := uuid.New()
	response := Response{FeedbackID: feedbackID}
	assert.Equal(t, feedbackID, response.FeedbackID)
}

func TestWithMetricValue(t *testing.T) {
	req := Request{}
	metricName := "rating"
	value := 4.5
	option := WithMetricValue(metricName, value)
	option(&req)
	assert.Equal(t, metricName, req.MetricName)
	assert.Equal(t, value, req.Value)
}

func TestWithComment(t *testing.T) {
	req := Request{}
	comment := "This is a test comment."
	option := WithComment(comment)
	option(&req)
	assert.Equal(t, "comment", req.MetricName)
	assert.Equal(t, comment, req.Value)
}

func TestWithRating(t *testing.T) {
	req := Request{}
	rating := 4.5
	option := WithRating(rating)
	option(&req)
	assert.Equal(t, "rating", req.MetricName)
	assert.Equal(t, rating, req.Value)
}

func TestWithBooleanMetric(t *testing.T) {
	req := Request{}
	metricName := "is_helpful"
	value := true
	option := WithBooleanMetric(metricName, value)
	option(&req)
	assert.Equal(t, metricName, req.MetricName)
	assert.Equal(t, value, req.Value)
}
