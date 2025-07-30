//go:build integration

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/feedback"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ComprehensiveFeedback(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Create test inferences first
	var testInferences []inference.InferenceResponse
	episodeID, _ := uuid.NewV7()

	t.Run("SetupTestInferences", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			resp, err := client.Inference(ctx, &inference.InferenceRequest{
				FunctionName: util.StringPtr("basic_test"),
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr(fmt.Sprintf("Test question %d: What is AI?", i+1))},
							},
						},
					},
				},
				EpisodeID: util.UUIDPtr(episodeID),
				Tags: map[string]string{
					"test_batch": "feedback_comprehensive",
					"question_id": fmt.Sprintf("q_%d", i+1),
				},
			})
			require.NoError(t, err)
			testInferences = append(testInferences, resp)
		}
		assert.Len(t, testInferences, 3)
	})

	t.Run("BooleanMetricFeedback", func(t *testing.T) {
		// Test various boolean metrics
		booleanTests := []struct {
			metricName string
			value      bool
			tags       map[string]string
		}{
			{
				metricName: "thumbs_up",
				value:      true,
				tags: map[string]string{
					"evaluator": "human",
					"confidence": "high",
				},
			},
			{
				metricName: "is_helpful",
				value:      true,
				tags: map[string]string{
					"evaluator": "automated",
					"model": "classifier_v1",
				},
			},
			{
				metricName: "is_accurate",
				value:      false,
				tags: map[string]string{
					"evaluator": "expert",
					"domain": "technical",
				},
			},
			{
				metricName: "follows_guidelines",
				value:      true,
				tags: map[string]string{
					"guideline_version": "v2.1",
					"checker": "automated",
				},
			},
		}

		for i, test := range booleanTests {
			feedbackResp, err := client.Feedback(ctx, &feedback.Request{
				MetricName:  test.metricName,
				Value:       test.value,
				InferenceID: util.UUIDPtr(testInferences[i%len(testInferences)].GetInferenceID()),
				Tags:        test.tags,
				Internal:    util.BoolPtr(false),
			})
			require.NoError(t, err, "Failed to submit boolean feedback for %s", test.metricName)
			assert.NotEqual(t, uuid.Nil, feedbackResp.FeedbackID)
		}
	})

	t.Run("FloatMetricFeedback", func(t *testing.T) {
		// Test various float metrics
		floatTests := []struct {
			metricName string
			value      float64
			tags       map[string]string
		}{
			{
				metricName: "quality_score",
				value:      8.5,
				tags: map[string]string{
					"scale": "1-10",
					"evaluator": "human_expert",
				},
			},
			{
				metricName: "relevance_score",
				value:      7.2,
				tags: map[string]string{
					"scale": "0-10",
					"context": "question_answering",
				},
			},
			{
				metricName: "creativity_score",
				value:      6.8,
				tags: map[string]string{
					"scale": "1-10",
					"domain": "creative_writing",
				},
			},
			{
				metricName: "response_time",
				value:      1.45,
				tags: map[string]string{
					"unit": "seconds",
					"measurement_type": "end_to_end",
				},
			},
			{
				metricName: "cost_efficiency",
				value:      0.023,
				tags: map[string]string{
					"unit": "dollars_per_request",
					"model": "gpt-4",
				},
			},
		}

		for i, test := range floatTests {
			feedbackResp, err := client.Feedback(ctx, &feedback.Request{
				MetricName:  test.metricName,
				Value:       test.value,
				InferenceID: util.UUIDPtr(testInferences[i%len(testInferences)].GetInferenceID()),
				Tags:        test.tags,
				Internal:    util.BoolPtr(true), // Mark as internal metric
			})
			require.NoError(t, err, "Failed to submit float feedback for %s", test.metricName)
			assert.NotEqual(t, uuid.Nil, feedbackResp.FeedbackID)
		}
	})

	t.Run("CommentFeedback", func(t *testing.T) {
		// Test comment feedback
		commentTests := []struct {
			comment string
			tags    map[string]string
		}{
			{
				comment: "The response was accurate but could be more concise. The explanation was helpful for understanding the concept.",
				tags: map[string]string{
					"feedback_type": "detailed_review",
					"reviewer": "domain_expert",
					"sentiment": "positive",
				},
			},
			{
				comment: "Missing key information about edge cases. Response seems incomplete.",
				tags: map[string]string{
					"feedback_type": "improvement_suggestion",
					"reviewer": "qa_specialist",
					"sentiment": "constructive",
				},
			},
			{
				comment: "Excellent response! Clear, comprehensive, and well-structured.",
				tags: map[string]string{
					"feedback_type": "praise",
					"reviewer": "end_user",
					"sentiment": "very_positive",
				},
			},
		}

		for i, test := range commentTests {
			feedbackResp, err := client.Feedback(ctx, &feedback.Request{
				MetricName:  "comment",
				Value:       test.comment,
				InferenceID: util.UUIDPtr(testInferences[i%len(testInferences)].GetInferenceID()),
				Tags:        test.tags,
				Internal:    util.BoolPtr(false),
			})
			require.NoError(t, err, "Failed to submit comment feedback")
			assert.NotEqual(t, uuid.Nil, feedbackResp.FeedbackID)
		}
	})

	t.Run("EpisodeLevelFeedback", func(t *testing.T) {
		// Test episode-level feedback
		episodeTests := []struct {
			metricName string
			value      interface{}
			tags       map[string]string
		}{
			{
				metricName: "conversation_quality",
				value:      8.7,
				tags: map[string]string{
					"conversation_length": "3_turns",
					"topic": "ai_explanation",
				},
			},
			{
				metricName: "user_satisfaction",
				value:      true,
				tags: map[string]string{
					"survey_response": "yes",
					"follow_up": "requested_more_info",
				},
			},
			{
				metricName: "episode_comment",
				value:      "The conversation flow was natural and the AI provided consistent, helpful responses throughout.",
				tags: map[string]string{
					"evaluation_type": "holistic",
					"evaluator": "conversation_analyst",
				},
			},
		}

		for _, test := range episodeTests {
			feedbackResp, err := client.Feedback(ctx, &feedback.Request{
				MetricName: test.metricName,
				Value:      test.value,
				EpisodeID:  util.UUIDPtr(episodeID),
				Tags:       test.tags,
				Internal:   util.BoolPtr(false),
			})
			require.NoError(t, err, "Failed to submit episode-level feedback for %s", test.metricName)
			assert.NotEqual(t, uuid.Nil, feedbackResp.FeedbackID)
		}
	})

	t.Run("DemonstrationFeedback", func(t *testing.T) {
		// Test demonstration feedback (providing example outputs)
		demonstrationResp, err := client.Feedback(ctx, &feedback.Request{
			MetricName: "demonstration",
			Value: map[string]interface{}{
				"improved_response": "Artificial Intelligence (AI) is a field of computer science focused on creating systems that can perform tasks typically requiring human intelligence. This includes learning, reasoning, problem-solving, and understanding natural language.",
				"explanation": "This response is more structured and provides a clear definition with key components.",
				"improvements": []string{
					"Added clear definition",
					"Mentioned key capabilities",
					"More concise structure",
				},
			},
			InferenceID: util.UUIDPtr(testInferences[0].GetInferenceID()),
			Tags: map[string]string{
				"demonstration_type": "improved_output",
				"improvement_focus": "clarity_and_structure",
				"provided_by": "content_expert",
			},
			Internal: util.BoolPtr(false),
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, demonstrationResp.FeedbackID)
	})

	t.Run("BatchFeedbackSubmission", func(t *testing.T) {
		// Submit multiple feedback items quickly to test batch handling
		var feedbackResponses []*feedback.Response
		
		for i := 0; i < 10; i++ {
			feedbackResp, err := client.Feedback(ctx, &feedback.Request{
				MetricName: "batch_test_metric",
				Value:      float64(i + 1),
				InferenceID: util.UUIDPtr(testInferences[i%len(testInferences)].GetInferenceID()),
				Tags: map[string]string{
					"batch_id": "batch_test_1",
					"item_index": fmt.Sprintf("%d", i),
					"timestamp": time.Now().Format(time.RFC3339),
				},
				Internal: util.BoolPtr(true),
			})
			require.NoError(t, err, "Failed to submit batch feedback item %d", i)
			feedbackResponses = append(feedbackResponses, feedbackResp)
		}

		// Verify all feedback was submitted successfully
		assert.Len(t, feedbackResponses, 10)
		
		// Verify all feedback IDs are unique
		feedbackIDSet := make(map[uuid.UUID]bool)
		for i, resp := range feedbackResponses {
			assert.NotEqual(t, uuid.Nil, resp.FeedbackID, "Feedback %d should have valid ID", i)
			assert.False(t, feedbackIDSet[resp.FeedbackID], "Feedback ID should be unique")
			feedbackIDSet[resp.FeedbackID] = true
		}
	})

	t.Run("FeedbackErrorHandling", func(t *testing.T) {
		// Test feedback with invalid inference ID
		invalidInferenceID, _ := uuid.NewV7()
		_, err := client.Feedback(ctx, &feedback.Request{
			MetricName:  "test_metric",
			Value:       true,
			InferenceID: util.UUIDPtr(invalidInferenceID),
		})
		assert.Error(t, err, "Should fail with invalid inference ID")
		
		var tzErr *shared.TensorZeroError
		assert.ErrorAs(t, err, &tzErr)

		// Test feedback with invalid episode ID
		invalidEpisodeID, _ := uuid.NewV7()
		_, err = client.Feedback(ctx, &feedback.Request{
			MetricName: "test_metric",
			Value:      true,
			EpisodeID:  util.UUIDPtr(invalidEpisodeID),
		})
		assert.Error(t, err, "Should fail with invalid episode ID")
		assert.ErrorAs(t, err, &tzErr)

		// Test feedback with both inference and episode ID (should be invalid)
		_, err = client.Feedback(ctx, &feedback.Request{
			MetricName:  "test_metric",
			Value:       true,
			InferenceID: util.UUIDPtr(testInferences[0].GetInferenceID()),
			EpisodeID:   util.UUIDPtr(episodeID),
		})
		assert.Error(t, err, "Should fail when both inference and episode ID are provided")

		// Test feedback with neither inference nor episode ID
		_, err = client.Feedback(ctx, &feedback.Request{
			MetricName: "test_metric",
			Value:      true,
		})
		assert.Error(t, err, "Should fail when neither inference nor episode ID is provided")
	})

	t.Run("FeedbackWithComplexValues", func(t *testing.T) {
		// Test feedback with complex structured values
		complexValue := map[string]interface{}{
			"overall_score": 8.5,
			"dimensions": map[string]interface{}{
				"accuracy":    9.0,
				"helpfulness": 8.0,
				"clarity":     8.5,
				"completeness": 7.5,
			},
			"strengths": []string{
				"Accurate information",
				"Clear explanation",
				"Good examples",
			},
			"weaknesses": []string{
				"Could be more comprehensive",
				"Missing edge cases",
			},
			"metadata": map[string]interface{}{
				"evaluation_time": time.Now().Unix(),
				"evaluator_confidence": 0.85,
				"evaluation_method": "structured_rubric",
			},
		}

		feedbackResp, err := client.Feedback(ctx, &feedback.Request{
			MetricName:  "detailed_evaluation",
			Value:       complexValue,
			InferenceID: util.UUIDPtr(testInferences[0].GetInferenceID()),
			Tags: map[string]string{
				"evaluation_type": "comprehensive",
				"rubric_version": "v3.2",
				"evaluator_id": "expert_001",
			},
			Internal: util.BoolPtr(false),
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, feedbackResp.FeedbackID)
	})
}