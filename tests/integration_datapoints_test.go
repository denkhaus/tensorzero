//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/denkhaus/tensorzero/datapoint"
	"github.com/denkhaus/tensorzero/inference"
	"github.com/denkhaus/tensorzero/shared"
	"github.com/denkhaus/tensorzero/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_DatapointOperations(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()
	datasetID, _ := uuid.NewV7()
	datasetName := "test_dataset_" + datasetID.String()[:8]

	t.Run("BulkInsertChatDatapoints", func(t *testing.T) {
		// Create chat datapoints for insertion
		chatDatapoints := []datapoint.DatapointInsert{
			&inference.ChatDatapointInsert{
				FunctionName: "basic_test",
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr("What is the capital of France?")},
							},
						},
					},
				},
				Output: []shared.ContentBlock{
					&shared.Text{Text: util.StringPtr("The capital of France is Paris.")},
				},
				Tags: map[string]string{
					"type": "qa",
					"difficulty": "easy",
				},
			},
			&inference.ChatDatapointInsert{
				FunctionName: "basic_test",
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr("Explain quantum computing")},
							},
						},
					},
				},
				Output: []shared.ContentBlock{
					&shared.Text{Text: util.StringPtr("Quantum computing uses quantum mechanical phenomena...")},
				},
				Tags: map[string]string{
					"type": "explanation",
					"difficulty": "hard",
				},
			},
		}

		// Insert datapoints
		ids, err := client.BulkInsertDatapoints(ctx, datasetName, chatDatapoints)
		require.NoError(t, err)
		assert.Len(t, ids, 2)

		// Verify all IDs are valid UUIDs
		for _, id := range ids {
			assert.NotEqual(t, uuid.Nil, id)
		}
	})

	t.Run("BulkInsertJsonDatapoints", func(t *testing.T) {
		// Create JSON datapoints for insertion
		jsonDatapoints := []datapoint.DatapointInsert{
			&inference.JsonDatapointInsert{
				FunctionName: "json_success",
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr("Extract entities from: John works at Google")},
							},
						},
					},
				},
				Output: map[string]interface{}{
					"entities": []map[string]interface{}{
						{"name": "John", "type": "person"},
						{"name": "Google", "type": "organization"},
					},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entities": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"name": map[string]interface{}{"type": "string"},
									"type": map[string]interface{}{"type": "string"},
								},
							},
						},
					},
				},
				Tags: map[string]string{
					"task": "entity_extraction",
					"domain": "business",
				},
			},
		}

		// Insert JSON datapoints
		ids, err := client.BulkInsertDatapoints(ctx, datasetName, jsonDatapoints)
		require.NoError(t, err)
		assert.Len(t, ids, 1)
		assert.NotEqual(t, uuid.Nil, ids[0])
	})

	t.Run("ListDatapoints", func(t *testing.T) {
		// List all datapoints in the dataset
		datapoints, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: datasetName,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(datapoints), 3) // At least 3 from previous tests

		// Verify datapoint structure
		for _, dp := range datapoints {
			assert.NotEqual(t, uuid.Nil, dp.ID)
			assert.Equal(t, datasetName, dp.DatasetName)
			assert.NotEmpty(t, dp.FunctionName)
			assert.NotNil(t, dp.Input)
			assert.NotNil(t, dp.Output)
		}
	})

	t.Run("ListDatapointsWithFilters", func(t *testing.T) {
		// List datapoints filtered by function name
		datapoints, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName:  datasetName,
			FunctionName: util.StringPtr("basic_test"),
			Limit:        util.IntPtr(10),
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(datapoints), 2) // Should have at least 2 chat datapoints

		// Verify all returned datapoints have the correct function name
		for _, dp := range datapoints {
			assert.Equal(t, "basic_test", dp.FunctionName)
		}
	})

	t.Run("ListDatapointsWithPagination", func(t *testing.T) {
		// Test pagination
		firstPage, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: datasetName,
			Limit:       util.IntPtr(2),
			Offset:      util.IntPtr(0),
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(firstPage), 2)

		secondPage, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: datasetName,
			Limit:       util.IntPtr(2),
			Offset:      util.IntPtr(2),
		})
		require.NoError(t, err)

		// Verify pages don't overlap
		if len(firstPage) > 0 && len(secondPage) > 0 {
			assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID)
		}
	})

	t.Run("DeleteDatapoint", func(t *testing.T) {
		// First, get a datapoint to delete
		datapoints, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: datasetName,
			Limit:       util.IntPtr(1),
		})
		require.NoError(t, err)
		require.NotEmpty(t, datapoints)

		datapointToDelete := datapoints[0]

		// Delete the datapoint
		err = client.DeleteDatapoint(ctx, datasetName, datapointToDelete.ID)
		require.NoError(t, err)

		// Verify it's deleted by trying to list with a specific function
		// (Note: The actual verification depends on TensorZero's delete behavior)
		remainingDatapoints, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: datasetName,
		})
		require.NoError(t, err)

		// Check that the deleted datapoint is not in the list
		for _, dp := range remainingDatapoints {
			assert.NotEqual(t, datapointToDelete.ID, dp.ID)
		}
	})
}

func TestIntegration_DatapointErrorHandling(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	t.Run("BulkInsertInvalidDatapoints", func(t *testing.T) {
		// Try to insert datapoints with invalid function name
		invalidDatapoints := []datapoint.DatapointInsert{
			&inference.ChatDatapointInsert{
				FunctionName: "nonexistent_function",
				Input: inference.InferenceInput{
					Messages: []shared.Message{
						{
							Role: "user",
							Content: []shared.ContentBlock{
								&shared.Text{Text: util.StringPtr("Test message")},
							},
						},
					},
				},
				Output: []shared.ContentBlock{
					&shared.Text{Text: util.StringPtr("Test response")},
				},
			},
		}

		_, err := client.BulkInsertDatapoints(ctx, "test_dataset", invalidDatapoints)
		assert.Error(t, err)
		// Should be a TensorZero error
		var tzErr *shared.TensorZeroError
		assert.ErrorAs(t, err, &tzErr)
	})

	t.Run("ListNonexistentDataset", func(t *testing.T) {
		// Try to list datapoints from a nonexistent dataset
		datapoints, err := client.ListDatapoints(ctx, &datapoint.ListDatapointsRequest{
			DatasetName: "nonexistent_dataset_" + func() string { id, _ := uuid.NewV7(); return id.String() }(),
		})

		// This might return empty list or error depending on TensorZero behavior
		if err != nil {
			var tzErr *shared.TensorZeroError
			assert.ErrorAs(t, err, &tzErr)
		} else {
			assert.Empty(t, datapoints)
		}
	})

	t.Run("DeleteNonexistentDatapoint", func(t *testing.T) {
		// Try to delete a nonexistent datapoint
		nonexistentID, _ := uuid.NewV7()
		err := client.DeleteDatapoint(ctx, "test_dataset", nonexistentID)
		
		// This should either succeed (idempotent) or return a specific error
		if err != nil {
			var tzErr *shared.TensorZeroError
			assert.ErrorAs(t, err, &tzErr)
		}
	})
}