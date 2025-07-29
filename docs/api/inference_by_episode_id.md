# InferenceByEpisodeId

The `InferenceByEpisodeId` table is a materialized view that indexes inferences by their episode ID, enabling efficient lookup of all inferences within an episode.
We store `episode_id_uint` as a `UInt128` so that they are sorted in the natural order by time as ClickHouse sorts UUIDs in little-endian order.

| Column | Type | Notes |
| --- | --- | --- |
| `episode_id_uint` | UInt128 | Integer representation of UUIDv7 for sorting order |
| `id_uint` | UInt128 | Integer representation of UUIDv7 for sorting order |
| `function_name` | String | Name of the function being called |
| `variant_name` | String | Name of the function variant |
| `function_type` | Enum(‘chat’, ‘json’) | Type of function (chat or json) |
