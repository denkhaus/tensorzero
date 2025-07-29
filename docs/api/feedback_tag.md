# FeedbackTag

The `FeedbackTag` table stores tags associated with various feedback types. Tags are used to categorize and add metadata to feedback entries, allowing for user-defined filtering later on. Data is inserted into this table by materialized views reading from the `BooleanMetricFeedback`, `CommentFeedback`, `DemonstrationFeedback`, and `FloatMetricFeedback` tables.

| Column | Type | Notes |
| --- | --- | --- |
| `metric_name` | String | Name of the metric the tag is associated with. |
| `key` | String | Key of the tag. |
| `value` | String | Value of the tag. |
| `feedback_id` | UUID | UUID referencing the related feedback entry (e.g., `BooleanMetricFeedback.id`). |
