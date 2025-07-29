# InferenceTag

The `InferenceTag` table stores tags associated with inferences. Tags are used to categorize and add metadata to inferences, allowing for user-defined filtering later on. Data is inserted into this table by materialized views reading from the `ChatInference` and `JsonInference` tables.

| Column | Type | Notes |
| --- | --- | --- |
| `function_name` | String | Name of the function the tag is associated with. |
| `key` | String | Key of the tag. |
| `value` | String | Value of the tag. |
| `inference_id` | UUID | UUID referencing the related inference (e.g., `ChatInference.id`). |
