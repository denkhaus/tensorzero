## `POST /feedback`

The `/feedback` endpoint assigns feedback to a particular inference or episode.

Each feedback is associated with a metric that is defined in the configuration file.

### Request

#### `dryrun`

* **Type:** boolean
* **Required:** no

If `true`, the feedback request will be executed but won’t be stored to the database (i.e. no-op).

This field is primarily for debugging and testing, and you should ignore it in production.

#### `episode_id`

* **Type:** UUID
* **Required:** when the metric level is `episode`

The episode ID to provide feedback for.

You should use this field when the metric level is `episode`.

Only use episode IDs that were returned by the TensorZero gateway.

#### `inference_id`

* **Type:** UUID
* **Required:** when the metric level is `inference`

The inference ID to provide feedback for.

You should use this field when the metric level is `inference`.

Only use inference IDs that were returned by the TensorZero gateway.

#### `metric_name`

* **Type:** string
* **Required:** yes

The name of the metric to provide feedback.

For example, if your metric is defined as `[metrics.draft_accepted]` in your configuration file, then you would set `metric_name: "draft_accepted"`.

The metric names `comment` and `demonstration` are reserved for special types of feedback.
A `comment` is free-form text (string) that can be assigned to either an inference or an episode.
The `demonstration` metric accepts values that would be a valid output.
See [Metrics & Feedback](/docs/gateway/guides/metrics-feedback/) for more details.

#### `tags`

* **Type:** flat JSON object with string keys and values
* **Required:** no

User-provided tags to associate with the feedback.

For example, `{"user_id": "123"}` or `{"author": "Alice"}`.

#### `value`

* **Type:** varies
* **Required:** yes

The value of the feedback.

The type of the value depends on the metric type (e.g. boolean for a metric with `type = "boolean"`).

### Response

#### `feedback_id`

* **Type:** UUID

The ID assigned to the feedback.

### Examples

#### Inference-Level Boolean Metric

Inference-Level Boolean Metric

##### Configuration

```
# ...

[metrics.draft_accepted]
type = "boolean"
level = "inference"
# ...
```

##### Request

* [Python](#tab-panel-234)
* [HTTP](#tab-panel-235)

```
from tensorzero import AsyncTensorZeroGateway

async with await AsyncTensorZeroGateway.build_http(gateway_url="http://localhost:3000") as client:
    result = await client.feedback(
        inference_id="00000000-0000-0000-0000-000000000000",
        metric_name="draft_accepted",
        value=True,
    )
```

##### Response

```
{ "feedback_id": "11111111-1111-1111-1111-111111111111" }
```

#### Episode-Level Float Metric

Episode-Level Float Metric

##### Configuration

```
# ...

[metrics.user_rating]
type = "float"
level = "episode"
# ...
```

##### Request

* [Python](#tab-panel-236)
* [HTTP](#tab-panel-237)

```
from tensorzero import AsyncTensorZeroGateway

async with await AsyncTensorZeroGateway.build_http(gateway_url="http://localhost:3000") as client:
    result = await client.feedback(
        episode_id="00000000-0000-0000-0000-000000000000",
        metric_name="user_rating",
        value=10,
    )
```

##### Response

```
{ "feedback_id": "11111111-1111-1111-1111-111111111111" }
