# BatchIdByInferenceId

The `BatchIdByInferenceId` table maps inference IDs to batch IDs, allowing for efficient lookup of which batch an inference belongs to.

| Column | Type | Notes |
| --- | --- | --- |
| `inference_id` | UUID | Must be a UUIDv7 |
| `batch_id` | UUID | Must be a UUIDv7 |
