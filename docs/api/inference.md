# POST /inference

## Description
The inference endpoint is the core of the TensorZero Gateway API. Under the hood, the gateway validates the request, samples a variant from the function, handles templating when applicable, and routes the inference to the appropriate model provider. If a problem occurs, it attempts to gracefully fallback to a different model provider or variant. After a successful inference, it returns the data to the client and asynchronously stores structured information in the database.

## Request Parameters

| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| `additional_tools` | `a list of tools` | `no` | `[]` | A list of tools defined at inference time that the model is allowed to call. This field allows for dynamic tool use, i.e. defining tools at runtime. You should prefer to define tools in the configuration file if possible. Only use this field if dynamic tool use is necessary for your use case. Each tool is an object with the following fields: `description`, `name`, `parameters`, and `strict`. The fields are identical to those in the configuration file, except that the `parameters` field should contain the JSON schema itself rather than a path to it. |
| `allowed_tools` | `list of strings` | `no` | | A list of tool names that the model is allowed to call. The tools must be defined in the configuration file. Any tools provided in `additional_tools` are always allowed, irrespective of this field. |
| `cache_options` | `object` | `no` | `{"enabled": "write_only"}` | Options for controlling inference caching behavior. |
| `credentials` | `object (a map from dynamic credential names to API keys)` | `no` | `no credentials` | Each model provider in your TensorZero configuration can be configured to accept credentials at inference time by using the `dynamic` location. The gateway expects the credentials to be provided in the `credentials` field of the request body. The gateway will return a 400 error if the credentials are not provided and the model provider has been configured with dynamic credentials. |
| `dryrun` | `boolean` | `no` | | If `true`, the inference request will be executed but won’t be stored to the database. The gateway will still call the downstream model providers. This field is primarily for debugging and testing, and you should generally not use it in production. |
| `episode_id` | `UUID` | `no` | | The ID of an existing episode to associate the inference with. For the first inference of a new episode, you should not provide an `episode_id`. If null, the gateway will generate a new episode ID and return it in the response. Only use episode IDs that were returned by the TensorZero gateway. |
| `extra_body` | `array of objects` | `no` | | The `extra_body` field allows you to modify the request body that TensorZero sends to a model provider. This advanced feature is an “escape hatch” that lets you use provider-specific functionality that TensorZero hasn’t implemented yet. Each object in the array must have three fields: `variant_name` or `model_provider_name`: The modification will only be applied to the specified variant or model provider, `pointer`: A [JSON Pointer](https://datatracker.ietf.org/doc/html/rfc6901) string specifying where to modify the request body, One of the following: `value`: The value to insert at that location; it can be of any type including nested types, `delete = true`: Deletes the field at the specified location, if present. |
| `extra_headers` | `array of objects` | `no` | | The `extra_headers` field allows you to modify the request headers that TensorZero sends to a model provider. This advanced feature is an “escape hatch” that lets you use provider-specific functionality that TensorZero hasn’t implemented yet. Each object in the array must have three fields: `variant_name` or `model_provider_name`: The modification will only be applied to the specified variant or model provider, `name`: The name of the header to modify, `value`: The value to set the header to. |
| `function_name` | `string` | `either `function_name` or `model_name` must be provided` | | The name of the function to call. The function must be defined in the configuration file. Alternatively, you can use the `model_name` field to call a model directly, without the need to define a function. |
| `include_original_response` | `boolean` | `no` | | If `true`, the original response from the model will be included in the response in the `original_response` field as a string. |
| `input` | `varies` | `yes` | | The input to the function. The type of the input depends on the function type. |
| `model_name` | `string` | `either `model_name` or `function_name` must be provided` | | The name of the model to call. Under the hood, the gateway will use a built-in passthrough chat function called `tensorzero::default`. |
| `output_schema` | `object (valid JSON Schema)` | `no` | | If set, this schema will override the `output_schema` defined in the function configuration for a JSON function. This dynamic output schema is used for validating the output of the function, and sent to providers which support structured outputs. |
| `parallel_tool_calls` | `boolean` | `no` | | If `true`, the function will be allowed to request multiple tool calls in a single conversation turn. If not set, we default to the configuration value for the function being called. Most model providers do not support parallel tool calls. In those cases, the gateway ignores this field. At the moment, only Fireworks AI and OpenAI support parallel tool calls. |
| `params` | `object` | `no` | `{}` | Override inference-time parameters for a particular variant type. This fields allows for dynamic inference parameters, i.e. defining parameters at runtime. This field’s format is `{ variant_type: { param: value, ... }, ... }`. You should prefer to set these parameters in the configuration file if possible. Only use this field if you need to set these parameters dynamically at runtime. Note that the parameters will apply to every variant of the specified type. Currently, we support the following: `chat_completion`, `frequency_penalty`, `json_mode`, `max_tokens`, `presence_penalty`, `seed`, `stop_sequences`, `temperature`, `top_p`. |
| `stream` | `boolean` | `no` | | If `true`, the gateway will stream the response from the model provider. |
| `tags` | `flat JSON object with string keys and values` | `no` | | User-provided tags to associate with the inference. For example, `{"user_id": "123"}` or `{"author": "Alice"}`. |
| `tool_choice` | `string` | `no` | | If set, overrides the tool choice strategy for the request. The supported tool choice strategies are: `none`: The function should not use any tools. `auto`: The model decides whether or not to use a tool. If it decides to use a tool, it also decides which tools to use. `required`: The model should use a tool. If multiple tools are available, the model decides which tool to use. `{ specific = "tool_name" }`: The model should use a specific tool. The tool must be defined in the `tools` section of the configuration file or provided in `additional_tools`. |
| `variant_name` | `string` | `no` | | If set, pins the inference request to a particular variant (not recommended). You should generally not set this field, and instead let the TensorZero gateway assign a variant. This field is primarily used for testing or debugging purposes. |

### `cache_options` Sub-parameters

| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| `enabled` | `string` | `no` | `"write_only"` | The cache mode to use. Must be one of: `"write_only"` (default): Only write to cache but don’t serve cached responses, `"read_only"`: Only read from cache but don’t write new entries, `"on"`: Both read from and write to cache, `"off"`: Disable caching completely. Note: When using `dryrun=true`, the gateway never writes to the cache. |
| `max_age_s` | `integer` | `no` | `null` | Maximum age in seconds for cache entries. If set, cached responses older than this value will not be used. |

### `input` Sub-parameters

| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| `messages` | `list of messages` | `no` | `[]` | A list of messages to provide to the model. Each message is an object with the following fields: `role`: The role of the message (`assistant` or `user`). `content`: The content of the message (string or list of content blocks). |
| `system` | `string or object` | `no` | | The input for the system message. If the function does not have a system schema, this field should be a string. If the function has a system schema, this field should be an object that matches the schema. |

#### `messages` Content Block Types

| Type | Description | Fields |
|---|---|---|
| `text` | Text for a text message. | `text` (string) or `arguments` (JSON object). |
| `tool_call` | Tool call. | `arguments`, `id`, `name`. |
| `tool_result` | Tool result. | `id`, `name`, `result`. |
| `file` | File. | `url` or `mime_type` and `data` (base64-encoded). |
| `raw_text` | Raw text. | `value`. |
| `thought` | Thought. | `text`. |
| `unknown` | Unknown content block. | `data`, `model_provider_name` (optional). |

## Response Structure

### Chat Function Response
When the function type is `chat`, the response is structured as follows.

| Field | Type | Description |
|---|---|---|
| `content` | `a list of content blocks` | The content blocks generated by the model. A content block can have `type` equal to `text` and `tool_call`. Reasoning models (e.g. DeepSeek R1) might also include `thought` content blocks. |
| `episode_id` | `UUID` | The ID of the episode associated with the inference. |
| `inference_id` | `UUID` | The ID assigned to the inference. |
| `original_response` | `string (optional)` | The original response from the model provider (only available when `include_original_response` is `true`). The returned data depends on the variant type: `chat_completion`, `experimental_best_of_n_sampling`, `experimental_mixture_of_n_sampling`, `experimental_dynamic_in_context_learning`, `experimental_chain_of_thought`. |
| `variant_name` | `string` | The name of the variant used for the inference. |
| `usage` | `object (optional)` | The usage metrics for the inference. Fields: `input_tokens` (integer), `output_tokens` (integer). |

#### `content` Content Block Types

| Type | Fields |
|---|---|
| `text` | `text` (string) |
| `tool_call` | `arguments` (object), `id` (string), `name` (string), `raw_arguments` (string), `raw_name` (string) |
| `thought` | `text` (string) |
| `unknown` | `data` (object), `model_provider_name` (string (optional)) |

### JSON Function Response
When the function type is `json`, the response is structured as follows.

| Field | Type | Description |
|---|---|---|
| `inference_id` | `UUID` | The ID assigned to the inference. |
| `episode_id` | `UUID` | The ID of the episode associated with the inference. |
| `original_response` | `string (optional)` | The original response from the model provider (only available when `include_original_response` is `true`). The returned data depends on the variant type: `chat_completion`, `experimental_best_of_n_sampling`, `experimental_mixture_of_n_sampling`, `experimental_dynamic_in_context_learning`, `experimental_chain_of_thought`. |
| `output` | `object` | The output object contains the following fields: |
| `variant_name` | `string` | The name of the variant used for the inference. |
| `usage` | `object (optional)` | The usage metrics for the inference. Fields: `input_tokens` (integer), `output_tokens` (integer). |

#### `output` Sub-fields

| Name | Type | Description |
|---|---|---|
| `raw` | `string` | The raw response from the model provider (which might be invalid JSON). |
| `parsed` | `object` | The parsed response from the model provider (`null` if invalid JSON). |
