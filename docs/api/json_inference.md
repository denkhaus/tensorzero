# JsonInference

The `JsonInference` table stores information about inference requests for JSON Functions made to the TensorZero Gateway.

A `JsonInference` row can be associated with one or more `ModelInference` rows, depending on the variant’s `type`.
For `chat_completion`, there will be a one-to-one relationship between rows in the two tables.
For other variant types, there might be more associated rows.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `function_name` | String |  |
| `variant_name` | String |  |
| `episode_id` | UUID | Must be a UUIDv7 |
| `input` | String (JSON) | `input` field in the `/inference` request body |
| `output` | String (JSON) | Object with `parsed` and `raw` fields |
| `output_schema` | String (JSON) | Schema that the output must conform to |
| `inference_params` | String (JSON) | Object with any inference parameters per variant type (e.g. `{"chat_completion": {"temperature": 0.5}}`) |
| `processing_time_ms` | UInt32 |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"user_id": "123"}`) |
