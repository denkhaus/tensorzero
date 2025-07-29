# ChatInferenceDataset

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
