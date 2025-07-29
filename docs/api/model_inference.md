# ModelInference

The `ModelInference` table stores information about each inference request to a model provider.
This is the inference request you’d make if you had called the model provider directly.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `inference_id` | UUID | Must be a UUIDv7 |
| `raw_request` | String | Raw request as sent to the model provider (varies) |
| `raw_response` | String | Raw response from the model provider (varies) |
| `model_name` | String | Name of the model used for the inference |
| `model_provider_name` | String | Name of the model provider used for the inference |
| `input_tokens` | Nullable(UInt32) |  |
| `output_tokens` | Nullable(UInt32) |  |
| `response_time_ms` | Nullable(UInt32) |  |
| `ttft_ms` | Nullable(UInt32) | Only available in streaming inferences |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `system` | Nullable(String) | The `system` input to the model |
| `input_messages` | Array(RequestMessage) | The user and assistant messages input to the model |
| `output` | Array(ContentBlock) | The output of the model |

A `RequestMessage` is an object with shape `{role: "user" | "assistant", content: List[ContentBlock]}` (content blocks are defined [here](about:/docs/gateway/api-reference/inference/#content-block)).
