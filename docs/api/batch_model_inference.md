# BatchModelInference

The `BatchModelInference` table stores information about inferences made as part of a batch request.
Once the request succeeds, we use this information to populate the `ChatInference`, `JsonInference`, and `ModelInference` tables.

| Column | Type | Notes |
| --- | --- | --- |
| `inference_id` | UUID | Must be a UUIDv7 |
| `batch_id` | UUID | Must be a UUIDv7 |
| `function_name` | String | Name of the function being called |
| `variant_name` | String | Name of the function variant |
| `episode_id` | UUID | Must be a UUIDv7 |
| `input` | String (JSON) | `input` field in the `/inference` request body |
| `system` | String | The `system` input to the model |
| `input_messages` | Array(RequestMessage) | The user and assistant messages input to the model |
| `tool_params` | String (JSON) | Object with any tool parameters (e.g. `tool_choice`, `tools_available`) used for the inference |
| `inference_params` | String (JSON) | Object with any inference parameters per variant type (e.g. `{"chat_completion": {"temperature": 0.5}}`) |
| `raw_request` | String | Raw request sent to the model provider |
| `model_name` | String | Name of the model used |
| `model_provider_name` | String | Name of the model provider |
| `output_schema` | String | Optional schema for JSON outputs |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
