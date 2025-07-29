# DynamicInContextLearningExample

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
