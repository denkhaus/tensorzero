# BatchRequest

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
