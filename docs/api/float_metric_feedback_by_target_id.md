# FloatMetricFeedbackByTargetId

The `FloatMetricFeedbackByTargetId` table indexes float metric feedback by target ID, enabling efficient lookup of feedback for a specific target.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 |
| `metric_name` | String | Name of the metric (stored as LowCardinality) |
| `value` | Float32 | The float feedback value |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |
