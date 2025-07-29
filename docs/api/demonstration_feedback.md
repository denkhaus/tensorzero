# DemonstrationFeedback

The `DemonstrationFeedback` table stores feedback in the form of demonstrations.
Demonstrations are examples of good behaviors.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `inference_id` | UUID | Must be a UUIDv7 |
| `value` | String | The demonstration or example provided as feedback (must match function output) |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |
