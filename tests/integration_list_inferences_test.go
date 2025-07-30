//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/feedback"
	"github.com/denkhaus/tensorzero/filter"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ListInferences(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Create some test inferences first
	episodeID, _ := uuid.NewV7()
	var inferenceIDs []uuid.UUID

	t.Run("SetupTestInferences", func(t *testing.T) {
		// Create multiple inferences for testing
		for i := 0; i < 5; i++ {
			resp, err := client.Inference(ctx, &inference.InferenceRequest{
				FunctionName: util.StringPtr("basic_test"),
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr("Test message " + string(rune(i+'1')))},
							},
						},
					},
				},
				EpisodeID: util.UUIDPtr(episodeID),
				Tags: map[string]string{
					"test_batch": "list_inferences",
					"index":      string(rune(i + '1')),
				},
			})
			require.NoError(t, err)
			inferenceIDs = append(inferenceIDs, resp.GetInferenceID())
		}

		// Add some feedback to test metric filtering
		for i, infID := range inferenceIDs {
			// Add boolean feedback
			_, err := client.Feedback(ctx, &feedback.Request{
				MetricName:  "thumbs_up",
				Value:       i%2 == 0, // Alternate true/false
				InferenceID: util.UUIDPtr(infID),
				Tags: map[string]string{
					"source": "test",
				},
			})
			require.NoError(t, err)

			// Add float feedback
			_, err = client.Feedback(ctx, &feedback.Request{
				MetricName:  "rating",
				Value:       float64(i + 1), // 1.0, 2.0, 3.0, 4.0, 5.0
				InferenceID: util.UUIDPtr(infID),
				Tags: map[string]string{
					"source": "test",
				},
			})
			require.NoError(t, err)
		}

		// Wait a bit for feedback to be processed
		time.Sleep(2 * time.Second)
	})

	t.Run("SimpleListInferences", func(t *testing.T) {
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(inferences), 5)

		// Verify structure
		for _, inf := range inferences {
			assert.NotEqual(t, uuid.Nil, inf.ID)
			assert.NotEqual(t, uuid.Nil, inf.EpisodeID)
			assert.Equal(t, "basic_test", inf.FunctionName)
			assert.NotEmpty(t, inf.VariantName)
			assert.NotNil(t, inf.Input)
			assert.NotNil(t, inf.Output)
		}
	})

	t.Run("ListInferencesByEpisode", func(t *testing.T) {
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			EpisodeID: util.UUIDPtr(episodeID),
			OrderBy:   &shared.OrderBy{By: "timestamp", Direction: "asc"},
		})
		require.NoError(t, err)
		assert.Len(t, inferences, 5)

		// Verify all belong to the same episode
		for _, inf := range inferences {
			assert.Equal(t, episodeID, inf.EpisodeID)
		}
	})

	t.Run("ListInferencesWithPagination", func(t *testing.T) {
		// First page
		firstPage, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Limit:        util.IntPtr(3),
			Offset:       util.IntPtr(0),
			OrderBy: &shared.OrderBy{By: "timestamp", Direction: "desc"},
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(firstPage), 3)

		// Second page
		secondPage, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Limit:        util.IntPtr(3),
			Offset:       util.IntPtr(3),
			OrderBy: &shared.OrderBy{By: "timestamp", Direction: "desc"},
		})
		require.NoError(t, err)

		// Verify no overlap
		if len(firstPage) > 0 && len(secondPage) > 0 {
			firstIDs := make(map[uuid.UUID]bool)
			for _, inf := range firstPage {
				firstIDs[inf.ID] = true
			}
			for _, inf := range secondPage {
				assert.False(t, firstIDs[inf.ID], "Found overlapping inference ID between pages")
			}
		}
	})
}

func TestIntegration_InferenceFiltering(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("BooleanMetricFilter", func(t *testing.T) {
		// Filter for inferences with thumbs_up = true
		boolFilter := filter.NewBooleanMetricFilter("thumbs_up", true)
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       boolFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Should have some results (from our setup)
		assert.Greater(t, len(inferences), 0)

		// Verify all have the expected metric value
		for _, inf := range inferences {
			if thumbsUp, exists := inf.MetricValues["thumbs_up"]; exists {
				assert.Equal(t, true, thumbsUp)
			}
		}
	})

	t.Run("FloatMetricFilter", func(t *testing.T) {
		// Filter for inferences with rating >= 3.0
		floatFilter := filter.NewFloatMetricFilter("rating", 3.0, ">=")
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       floatFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Should have some results
		assert.Greater(t, len(inferences), 0)

		// Verify all have rating >= 3.0
		for _, inf := range inferences {
			if rating, exists := inf.MetricValues["rating"]; exists {
				if ratingFloat, ok := rating.(float64); ok {
					assert.GreaterOrEqual(t, ratingFloat, 3.0)
				}
			}
		}
	})

	t.Run("TagFilter", func(t *testing.T) {
		// Filter by tag
		tagFilter := filter.NewTagFilter("test_batch", "list_inferences", "=")
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			Filter: tagFilter,
			Limit:  util.IntPtr(10),
		})
		require.NoError(t, err)

		// Should have our test inferences
		assert.GreaterOrEqual(t, len(inferences), 5)

		// Verify all have the expected tag
		for _, inf := range inferences {
			if testBatch, exists := inf.Tags["test_batch"]; exists {
				assert.Equal(t, "list_inferences", testBatch)
			}
		}
	})

	t.Run("AndFilter", func(t *testing.T) {
		// Combine boolean and float filters with AND
		boolFilter := filter.NewBooleanMetricFilter("thumbs_up", true)
		floatFilter := filter.NewFloatMetricFilter("rating", 3.0, ">=")
		andFilter := filter.NewAndFilter(boolFilter, floatFilter)
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       andFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Verify results match both conditions
		for _, inf := range inferences {
			if thumbsUp, exists := inf.MetricValues["thumbs_up"]; exists {
				assert.Equal(t, true, thumbsUp)
			}
			if rating, exists := inf.MetricValues["rating"]; exists {
				if ratingFloat, ok := rating.(float64); ok {
					assert.GreaterOrEqual(t, ratingFloat, 3.0)
				}
			}
		}
	})

	t.Run("OrFilter", func(t *testing.T) {
		// Combine different conditions with OR
		highRating := filter.NewFloatMetricFilter("rating", 4.0, ">=")
		thumbsUp := filter.NewBooleanMetricFilter("thumbs_up", true)
		orFilter := filter.NewOrFilter(highRating, thumbsUp)
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       orFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Should have results that match either condition
		assert.Greater(t, len(inferences), 0)

		// Verify each result matches at least one condition
		for _, inf := range inferences {
			matchesCondition := false
			
			if rating, exists := inf.MetricValues["rating"]; exists {
				if ratingFloat, ok := rating.(float64); ok && ratingFloat >= 4.0 {
					matchesCondition = true
				}
			}
			
			if thumbsUp, exists := inf.MetricValues["thumbs_up"]; exists {
				if thumbsUpBool, ok := thumbsUp.(bool); ok && thumbsUpBool {
					matchesCondition = true
				}
			}
			
			assert.True(t, matchesCondition, "Inference should match at least one OR condition")
		}
	})

	t.Run("NotFilter", func(t *testing.T) {
		// Filter for NOT thumbs_up
		thumbsUpFilter := filter.NewBooleanMetricFilter("thumbs_up", true)
		notFilter := filter.NewNotFilter(thumbsUpFilter)
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       notFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Verify none have thumbs_up = true
		for _, inf := range inferences {
			if thumbsUp, exists := inf.MetricValues["thumbs_up"]; exists {
				if thumbsUpBool, ok := thumbsUp.(bool); ok {
					assert.False(t, thumbsUpBool, "Should not have thumbs_up = true with NOT filter")
				}
			}
		}
	})

	t.Run("ComplexNestedFilter", func(t *testing.T) {
		// Complex filter: (rating >= 3 AND thumbs_up = true) OR (rating = 1)
		highRatingAndThumbsUp := filter.NewAndFilter(
			filter.NewFloatMetricFilter("rating", 3.0, ">="),
			filter.NewBooleanMetricFilter("thumbs_up", true),
		)
		lowRating := filter.NewFloatMetricFilter("rating", 1.0, "=")
		complexFilter := filter.NewOrFilter(highRatingAndThumbsUp, lowRating)
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       complexFilter,
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)

		// Should have some results
		assert.Greater(t, len(inferences), 0)

		// Verify each result matches the complex condition
		for _, inf := range inferences {
			matchesCondition := false
			
			// Check first condition: rating >= 3 AND thumbs_up = true
			if rating, ratingExists := inf.MetricValues["rating"]; ratingExists {
				if thumbsUp, thumbsExists := inf.MetricValues["thumbs_up"]; thumbsExists {
					if ratingFloat, ok := rating.(float64); ok {
						if thumbsUpBool, ok := thumbsUp.(bool); ok {
							if ratingFloat >= 3.0 && thumbsUpBool {
								matchesCondition = true
							}
						}
					}
				}
			}
			
			// Check second condition: rating = 1
			if rating, exists := inf.MetricValues["rating"]; exists {
				if ratingFloat, ok := rating.(float64); ok && ratingFloat == 1.0 {
					matchesCondition = true
				}
			}
			
			assert.True(t, matchesCondition, "Inference should match complex filter condition")
		}
	})

	t.Run("TimeFilter", func(t *testing.T) {
		// Filter for inferences from the last hour
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		timeFilter := filter.NewTimeFilter(oneHourAgo.Format(time.RFC3339), ">=")
		
		inferences, err := client.ListInferences(ctx, &inference.ListInferencesRequest{
			FunctionName: util.StringPtr("basic_test"),
			Filter:       timeFilter,
			Limit:        util.IntPtr(20),
		})
		require.NoError(t, err)

		// Should have our recent test inferences
		assert.GreaterOrEqual(t, len(inferences), 5)

		// Verify timestamps are recent
		for _, inf := range inferences {
			// Parse timestamp and verify it's recent
			timestamp, err := time.Parse(time.RFC3339, inf.Timestamp)
			if err == nil {
				assert.True(t, timestamp.After(oneHourAgo), "Inference timestamp should be within the last hour")
			}
		}
	})
}