# DemonstrationFeedbackByInferenceId

The `DemonstrationFeedbackByInferenceId` table stores demonstration feedback associated with inferences, enabling efficient lookup of demonstrations by inference ID.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `inference_id` | UUID | Must be a UUIDv7 |
| `value` | String | The demonstration feedback content |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |
