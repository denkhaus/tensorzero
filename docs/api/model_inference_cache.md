# ModelInferenceCache

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
