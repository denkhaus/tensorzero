# BooleanMetricFeedback

The `BooleanMetricFeedback` table stores feedback for metrics of `type = "boolean"`.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 that is either `inference_id` or `episode_id` depending on `level` in metric config |
| `metric_name` | String |  |
| `value` | Bool |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |
