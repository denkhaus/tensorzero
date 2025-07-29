# CommentFeedback

The `CommentFeedback` table stores feedback provided with `metric_name` of `"comment"`.
Comments are free-form text feedbacks.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 that is either `inference_id` or `episode_id` depending on `level` in metric config |
| `target_type` | `"inference"` or `"episode"` |  |
| `value` | String |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |
