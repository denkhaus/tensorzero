import re
import os

doc_content = """
The TensorZero Gateway stores inference and feedback data in ClickHouse.
This data can be used for observability, experimentation, and optimization.

## `ChatInference`

The `ChatInference` table stores information about inference requests for Chat Functions made to the TensorZero Gateway.

A `ChatInference` row can be associated with one or more `ModelInference` rows, depending on the variant’s `type`.
For `chat_completion`, there will be a one-to-one relationship between rows in the two tables.
For other variant types, there might be more associated rows.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `function_name` | String |  |
| `variant_name` | String |  |
| `episode_id` | UUID | Must be a UUIDv7 |
| `input` | String (JSON) | `input` field in the `/inference` request body |
| `output` | String (JSON) | Array of content blocks |
| `tool_params` | String (JSON) | Object with any tool parameters (e.g. `tool_choice`, `tools_available`) used for the inference |
| `inference_params` | String (JSON) | Object with any inference parameters per variant type (e.g. `{"chat_completion": {"temperature": 0.5}}`) |
| `processing_time_ms` | UInt32 |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"user_id": "123"}`) |

## `JsonInference`

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

## `ModelInference`

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

## `DynamicInContextLearningExample`

The `DynamicInContextLearningExample` table stores examples for dynamic in-context learning variants.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `function_name` | String |  |
| `variant_name` | String |  |
| `namespace` | String |  |
| `input` | String (JSON) |  |
| `output` | String |  |
| `embedding` | Array(Float32) |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |

## `BooleanMetricFeedback`

The `BooleanMetricFeedback` table stores feedback for metrics of `type = "boolean"`.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 that is either `inference_id` or `episode_id` depending on `level` in metric config |
| `metric_name` | String |  |
| `value` | Bool |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |

## `FloatMetricFeedback`

table stores feedback for metrics of `type = "float"`.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 that is either `inference_id` or `episode_id` depending on `level` in metric config |
| `metric_name` | String |  |
| `value` | Float32 |  |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |

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

## `DemonstrationFeedback`

The `DemonstrationFeedback` table stores feedback in the form of demonstrations.
Demonstrations are examples of good behaviors.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `inference_id` | UUID | Must be a UUIDv7 |
| `value` | String | The demonstration or example provided as feedback (must match function output) |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"author": "Alice"}`) |

## `ModelInferenceCache`

The `ModelInferenceCache` table stores cached model inference results to avoid duplicate requests.

| Column | Type | Notes |
| --- | --- | --- |
| `short_cache_key` | UInt64 | First part of composite key for fast lookups |
| `long_cache_key` | FixedString(64) | Hex-encoded 256-bit key for full cache validation |
| `timestamp` | DateTime | When this cache entry was created, defaults to now() |
| `output` | String | The cached model output |
| `raw_request` | String | Raw request that was sent to the model provider |
| `raw_response` | String | Raw response received from the model provider |
| `is_deleted` | Bool | Soft deletion flag, defaults to false |

The table uses the `ReplacingMergeTree` engine with `timestamp` and `is_deleted` columns for deduplication.
It is partitioned by month and ordered by the composite cache key `(short_cache_key, long_cache_key)`.
The `short_cache_key` serves as the primary key for performance, while a bloom filter index on `long_cache_key`
helps optimize point queries.

## `ChatInferenceDataset`

The `ChatInferenceDataset` table stores chat inference examples organized into datasets.

| Column | Type | Notes |
| --- | --- | --- |
| `dataset_name` | LowCardinality(String) | Name of the dataset this example belongs to |
| `function_name` | LowCardinality(String) | Name of the function this example is for |
| `id` | UUID | Must be a UUIDv7, often the inference ID if generated from an inference |
| `episode_id` | UUID | Must be a UUIDv7 |
| `input` | String (JSON) | `input` field in the `/inference` request body |
| `output` | Nullable(String) (JSON) | Array of content blocks |
| `tool_params` | String (JSON) | Object with any tool parameters (e.g. `tool_choice`, `tools_available`) used for the inference |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"user_id": "123"}`) |
| `auxiliary` | String | Additional JSON data (unstructured) |
| `is_deleted` | Bool | Soft deletion flag, defaults to false |
| `updated_at` | DateTime | When this dataset entry was updated, defaults to now() |

The table uses the `ReplacingMergeTree` engine with `updated_at` and `is_deleted` columns for deduplication.
It is ordered by `dataset_name`, `function_name`, and `id` to optimize queries filtering by dataset and function.

## `JsonInferenceDataset`

The `JsonInferenceDataset` table stores JSON inference examples organized into datasets.

| Column | Type | Notes |
| --- | --- | --- |
| `dataset_name` | LowCardinality(String) | Name of the dataset this example belongs to |
| `function_name` | LowCardinality(String) | Name of the function this example is for |
| `id` | UUID | Must be a UUIDv7, often the inference ID if generated from an inference |
| `episode_id` | UUID | Must be a UUIDv7 |
| `input` | String (JSON) | `input` field in the `/inference` request body |
| `output` | String (JSON) | Object with `parsed` and `raw` fields |
| `output_schema` | String (JSON) | Schema that the output must conform to |
| `tags` | Map(String, String) | User-assigned tags (e.g. `{"user_id": "123"}`) |
| `auxiliary` | String | Additional JSON data (unstructured) |
| `is_deleted` | Bool | Soft deletion flag, defaults to false |
| `updated_at` | DateTime | When this dataset entry was updated, defaults to now() |

The table uses the `ReplacingMergeTree` engine with `updated_at` and `is_deleted` columns for deduplication.
It is ordered by `dataset_name`, `function_name`, and `id` to optimize queries filtering by dataset and function.

## `BatchRequest`

The `BatchRequest` table stores information about batch requests made to model providers. We update it every time a particular `batch_id` is created or polled.

| Column | Type | Notes |
| --- | --- | --- |
| `batch_id` | UUID | Must be a UUIDv7 |
| `id` | UUID | Must be a UUIDv7 |
| `batch_params` | String | Parameters used for the batch request |
| `model_name` | String | Name of the model used |
| `model_provider_name` | String | Name of the model provider |
| `status` | String | One of: ‘pending’, ‘completed’, ‘failed’ |
| `errors` | Array(String) | Array of error messages if status is ‘failed’ |
| `timestamp` | DateTime | Materialized from `id` (using `UUIDv7ToDateTime` function) |
| `raw_request` | String | Raw request sent to the model provider |
| `raw_response` | String | Raw response received from the model provider |
| `function_name` | String | Name of the function being called |
| `variant_name` | String | Name of the function variant |

## `BatchModelInference`

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

Materialized View Tables

[Materialized views](https://clickhouse.com/docs/en/materialized-view) in columnar databases like ClickHouse pre-compute alternative indexings of data, dramatically improving query performance compared to computing results on-the-fly.
In TensorZero’s case, we store denormalized data about inferences and feedback in the materialized views below to support efficient queries for common downstream use cases.

## `FeedbackTag`

The `FeedbackTag` table stores tags associated with various feedback types. Tags are used to categorize and add metadata to feedback entries, allowing for user-defined filtering later on. Data is inserted into this table by materialized views reading from the `BooleanMetricFeedback`, `CommentFeedback`, `DemonstrationFeedback`, and `FloatMetricFeedback` tables.

| Column | Type | Notes |
| --- | --- | --- |
| `metric_name` | String | Name of the metric the tag is associated with. |
| `key` | String | Key of the tag. |
| `value` | String | Value of the tag. |
| `feedback_id` | UUID | UUID referencing the related feedback entry (e.g., `BooleanMetricFeedback.id`). |

## `InferenceById`

The `InferenceById` table is a materialized view that combines data from `ChatInference` and `JSONInference`.
Notably, it indexes the table by `id_uint` for fast lookup by the gateway to validate feedback requests.
We store `id_uint` as a UInt128 so that they are sorted in the natural order by time as ClickHouse sorts UUIDs in little-endian order.

| Column | Type | Notes |
| --- | --- | --- |
| `id_uint` | UInt128 | Integer representation of UUIDv7 for sorting order |
| `function_name` | String |  |
| `variant_name` | String |  |
| `episode_id` | UUID | Must be a UUIDv7 |
| `function_type` | String | Either `'chat'` or `'json'` |

## `InferenceByEpisodeId`

The `InferenceByEpisodeId` table is a materialized view that indexes inferences by their episode ID, enabling efficient lookup of all inferences within an episode.
We store `episode_id_uint` as a `UInt128` so that they are sorted in the natural order by time as ClickHouse sorts UUIDs in little-endian order.

| Column | Type | Notes |
| --- | --- | --- |
| `episode_id_uint` | UInt128 | Integer representation of UUIDv7 for sorting order |
| `id_uint` | UInt128 | Integer representation of UUIDv7 for sorting order |
| `function_name` | String | Name of the function being called |
| `variant_name` | String | Name of the function variant |
| `function_type` | Enum(‘chat’, ‘json’) | Type of function (chat or json) |

## `InferenceTag`

The `InferenceTag` table stores tags associated with inferences. Tags are used to categorize and add metadata to inferences, allowing for user-defined filtering later on. Data is inserted into this table by materialized views reading from the `ChatInference` and `JsonInference` tables.

| Column | Type | Notes |
| --- | --- | --- |
| `function_name` | String | Name of the function the tag is associated with. |
| `key` | String | Key of the tag. |
| `value` | String | Value of the tag. |
| `inference_id` | UUID | UUID referencing the related inference (e.g., `ChatInference.id`). |

## `BatchIdByInferenceId`

The `BatchIdByInferenceId` table maps inference IDs to batch IDs, allowing for efficient lookup of which batch an inference belongs to.

| Column | Type | Notes |
| --- | --- | --- |
| `inference_id` | UUID | Must be a UUIDv7 |
| `batch_id` | UUID | Must be a UUIDv7 |

## `BooleanMetricFeedbackByTargetId`

The `BooleanMetricFeedbackByTargetId` table indexes boolean metric feedback by target ID, enabling efficient lookup of feedback for a specific target.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 |
| `metric_name` | String | Name of the metric (stored as LowCardinality) |
| `value` | Bool | The boolean feedback value |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |

The `CommentFeedbackByTargetId` table stores text feedback associated with inferences or episodes, enabling efficient lookup of comments by their target ID.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 |
| `target_type` | Enum(‘inference’, ‘episode’) | Type of entity this feedback is for |
| `value` | String | The text feedback content |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |

## `DemonstrationFeedbackByInferenceId`

The `DemonstrationFeedbackByInferenceId` table stores demonstration feedback associated with inferences, enabling efficient lookup of demonstrations by inference ID.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `inference_id` | UUID | Must be a UUIDv7 |
| `value` | String | The demonstration feedback content |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |

## `FloatMetricFeedbackByTargetId`

The `FloatMetricFeedbackByTargetId` table indexes float metric feedback by target ID, enabling efficient lookup of feedback for a specific target.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 |
| `metric_name` | String | Name of the metric (stored as LowCardinality) |
| `value` | Float32 | The float feedback value |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |<environment_details>
# VSCode Visible Files
types/request.go

# VSCode Open Tabs
client.go,types/inference-request-options_test.go,types/interfaces_test.go,tests/openai_test.go,types/shared.go,types/request.go,types/shared_test.go,parsers.go,util/util.go,go.mod,types/filters.go,types/response.go,types/interfaces.go,types/config.go

# Current Time
Current time in ISO 8601 UTC format: 2025-07-29T17:02:37.387Z
User time zone: Europe/Berlin, UTC+2:00

# Current Cost
$0.00

# Current Mode
<slug>code</slug>
<name>Code</name>
<model>gemini-2.5-flash</model>
====

REMINDERS

Below is your current list of reminders for this task. Keep them updated as you progress.

| # | Content | Status |
|---|---------|--------|
| 1 | Access and read the documentation from the provided URL. | Completed |
| 2 | Parse the documentation content and identify key sections for markdown files. | In Progress |
| 3 | Create multiple markdown files in `docs/api` based on the parsed content. | Pending |
| 4 | Confirm all information has been successfully written to markdown files. | Pending |


IMPORTANT: When task status changes, remember to call the `update_todo_list` tool to update your progress.

</environment_details>
