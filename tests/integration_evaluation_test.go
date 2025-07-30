//go:build integration

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denkhaus/tensorzero/evaluation"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_DynamicEvaluation(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("CreateEvaluationRun", func(t *testing.T) {
		// Create a dynamic evaluation run
		runResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
			Variants: map[string]string{
				"variant_a": "basic_test",
				"variant_b": "basic_test_template_no_schema",
			},
			DisplayName: util.StringPtr("Integration Test Evaluation"),
			ProjectName: util.StringPtr("go_client_tests"),
			Tags: map[string]string{
				"test_type": "integration",
				"client":    "go",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, runResp.RunID)

		// Store run ID for subsequent tests
		runID := runResp.RunID

		t.Run("CreateEvaluationEpisodes", func(t *testing.T) {
			// Create multiple episodes for the evaluation run
			episodeTests := []struct {
				taskName      string
				datapointName string
				tags          map[string]string
			}{
				{
					taskName:      "question_answering",
					datapointName: "qa_sample_1",
					tags: map[string]string{
						"difficulty": "easy",
						"category":   "factual",
					},
				},
				{
					taskName:      "question_answering",
					datapointName: "qa_sample_2",
					tags: map[string]string{
						"difficulty": "medium",
						"category":   "reasoning",
					},
				},
				{
					taskName:      "creative_writing",
					datapointName: "story_prompt_1",
					tags: map[string]string{
						"difficulty": "hard",
						"category":   "creative",
					},
				},
			}

			var episodeIDs []uuid.UUID

			for i, test := range episodeTests {
				episodeResp, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
					RunID:         runID,
					TaskName:      util.StringPtr(test.taskName),
					DatapointName: util.StringPtr(test.datapointName),
					Tags:          test.tags,
				})
				require.NoError(t, err, "Failed to create episode %d", i)
				assert.NotEqual(t, uuid.Nil, episodeResp.EpisodeID)
				episodeIDs = append(episodeIDs, episodeResp.EpisodeID)
			}

			// Verify we created the expected number of episodes
			assert.Len(t, episodeIDs, len(episodeTests))

			// Verify all episode IDs are unique
			episodeIDSet := make(map[uuid.UUID]bool)
			for _, id := range episodeIDs {
				assert.False(t, episodeIDSet[id], "Episode ID should be unique")
				episodeIDSet[id] = true
			}
		})

		t.Run("CreateEvaluationRunWithMinimalData", func(t *testing.T) {
			// Test with minimal required data
			minimalRunResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
				Variants: map[string]string{
					"single_variant": "basic_test",
				},
			})
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, minimalRunResp.RunID)
			assert.NotEqual(t, runID, minimalRunResp.RunID, "Should create a new run with different ID")
		})

		t.Run("CreateEpisodeWithMinimalData", func(t *testing.T) {
			// Create episode with only required fields
			minimalEpisodeResp, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
				RunID: runID,
			})
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, minimalEpisodeResp.EpisodeID)
		})
	})

	t.Run("EvaluationRunWithMultipleVariants", func(t *testing.T) {
		// Test with multiple variants
		multiVariantResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
			Variants: map[string]string{
				"variant_1": "basic_test",
				"variant_2": "basic_test_template_no_schema",
				"variant_3": "json_success",
			},
			DisplayName: util.StringPtr("Multi-Variant Evaluation"),
			ProjectName: util.StringPtr("go_client_comprehensive_tests"),
			Tags: map[string]string{
				"test_type":     "multi_variant",
				"variant_count": "3",
				"purpose":       "comparison",
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, multiVariantResp.RunID)

		// Create episodes for each variant type
		variantTests := []struct {
			name string
			tags map[string]string
		}{
			{
				name: "chat_test",
				tags: map[string]string{"type": "chat", "complexity": "simple"},
			},
			{
				name: "json_test",
				tags: map[string]string{"type": "json", "complexity": "structured"},
			},
			{
				name: "template_test",
				tags: map[string]string{"type": "template", "complexity": "dynamic"},
			},
		}

		for _, test := range variantTests {
			episodeResp, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
				RunID:         multiVariantResp.RunID,
				TaskName:      util.StringPtr(test.name),
				DatapointName: util.StringPtr("test_datapoint_" + test.name),
				Tags:          test.tags,
			})
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, episodeResp.EpisodeID)
		}
	})

	t.Run("EvaluationErrorHandling", func(t *testing.T) {
		// Test with invalid variant names
		_, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
			Variants: map[string]string{
				"invalid_variant": "nonexistent_function",
			},
			DisplayName: util.StringPtr("Error Test"),
		})
		
		// This might succeed at creation time but fail during execution
		// The behavior depends on TensorZero's validation strategy
		if err != nil {
			t.Logf("Expected error for invalid variant: %v", err)
		}

		// Test creating episode with invalid run ID
		invalidRunID, _ := uuid.NewV7()
		_, err = client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
			RunID:         invalidRunID,
			TaskName:      util.StringPtr("test_task"),
			DatapointName: util.StringPtr("test_datapoint"),
		})
		assert.Error(t, err, "Should fail with invalid run ID")
	})

	t.Run("EvaluationRunWithComplexTags", func(t *testing.T) {
		// Test with complex tag structures
		complexTagsResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
			Variants: map[string]string{
				"baseline": "basic_test",
				"optimized": "basic_test_template_no_schema",
			},
			DisplayName: util.StringPtr("Complex Tags Evaluation"),
			ProjectName: util.StringPtr("advanced_testing"),
			Tags: map[string]string{
				"experiment_id":     func() string { id, _ := uuid.NewV7(); return id.String() }(),
				"researcher":        "go_integration_test",
				"hypothesis":        "template_performs_better",
				"expected_outcome":  "improved_accuracy",
				"dataset_version":   "v2.1",
				"model_temperature": "0.7",
				"max_tokens":        "150",
				"evaluation_type":   "comparative",
				"metrics":           "accuracy,latency,cost",
				"environment":       "test",
			},
		})
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, complexTagsResp.RunID)

		// Create episode with equally complex tags
		_, err = client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
			RunID:         complexTagsResp.RunID,
			TaskName:      util.StringPtr("comprehensive_qa"),
			DatapointName: util.StringPtr("complex_reasoning_sample"),
			Tags: map[string]string{
				"episode_id":        func() string { id, _ := uuid.NewV7(); return id.String() }(),
				"difficulty_level":  "expert",
				"domain":           "science",
				"subdomain":        "physics",
				"question_type":    "multi_step_reasoning",
				"expected_tokens":  "200-300",
				"reference_answer": "provided",
				"evaluation_criteria": "accuracy,completeness,clarity",
				"time_limit":       "30s",
				"retry_count":      "0",
			},
		})
		require.NoError(t, err)
	})
}

func TestIntegration_EvaluationWorkflow(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("EndToEndEvaluationWorkflow", func(t *testing.T) {
		// Step 1: Create evaluation run
		runResp, err := client.DynamicEvaluationRun(ctx, &evaluation.RunRequest{
			Variants: map[string]string{
				"control":     "basic_test",
				"experiment":  "basic_test_template_no_schema",
			},
			DisplayName: util.StringPtr("End-to-End Workflow Test"),
			ProjectName: util.StringPtr("workflow_validation"),
			Tags: map[string]string{
				"workflow_step": "1_create_run",
				"test_phase":    "end_to_end",
			},
		})
		require.NoError(t, err)
		runID := runResp.RunID

		// Step 2: Create multiple episodes
		episodeCount := 3
		var episodeIDs []uuid.UUID

		for i := 0; i < episodeCount; i++ {
			episodeResp, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
				RunID:         runID,
				TaskName:      util.StringPtr("workflow_task"),
				DatapointName: util.StringPtr(fmt.Sprintf("datapoint_%d", i+1)),
				Tags: map[string]string{
					"workflow_step": "2_create_episodes",
					"episode_index": fmt.Sprintf("%d", i+1),
					"batch_id":      "workflow_batch_1",
				},
			})
			require.NoError(t, err)
			episodeIDs = append(episodeIDs, episodeResp.EpisodeID)
		}

		// Step 3: Verify all episodes were created
		assert.Len(t, episodeIDs, episodeCount)
		for i, episodeID := range episodeIDs {
			assert.NotEqual(t, uuid.Nil, episodeID, "Episode %d should have valid ID", i+1)
		}

		// Step 4: Create additional episodes with different task types
		taskTypes := []string{"reasoning", "creativity", "factual_recall"}
		for _, taskType := range taskTypes {
			_, err := client.DynamicEvaluationRunEpisode(ctx, &evaluation.EpisodeRequest{
				RunID:         runID,
				TaskName:      util.StringPtr(taskType),
				DatapointName: util.StringPtr("specialized_" + taskType),
				Tags: map[string]string{
					"workflow_step": "3_specialized_episodes",
					"task_category": taskType,
					"specialization": "true",
				},
			})
			require.NoError(t, err)
		}

		// The workflow is complete - in a real scenario, you would now:
		// - Run inferences for each episode
		// - Collect feedback/metrics
		// - Analyze results
		// This demonstrates the evaluation setup is working correctly
	})
}