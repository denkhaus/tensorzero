# InferenceById

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
