# CommentFeedbackByTargetId

The `CommentFeedbackByTargetId` table stores text feedback associated with inferences or episodes, enabling efficient lookup of comments by their target ID.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | UUID | Must be a UUIDv7 |
| `target_id` | UUID | Must be a UUIDv7 |
| `target_type` | Enum(‘inference’, ‘episode’) | Type of entity this feedback is for |
| `value` | String | The text feedback content |
| `tags` | Map(String, String) | Key-value pairs of tags associated with the feedback |
