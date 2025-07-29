package tensorzero

import (
	"encoding/json"
	"fmt"

	"github.com/denkhaus/tensorzero/types"
	"github.com/google/uuid"
)

// parseInferenceResponse parses an inference response from JSON
func parseInferenceResponse(data []byte) (types.InferenceResponse, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if content, hasContent := raw["content"]; hasContent {
		if contentList, ok := content.([]interface{}); ok {
			// Parse as ChatInferenceResponse but handle content separately
			var resp types.ChatInferenceResponse

			// Parse basic fields
			if inferenceID, ok := raw["inference_id"].(string); ok {
				resp.InferenceID = uuid.MustParse(inferenceID)
			}
			if episodeID, ok := raw["episode_id"].(string); ok {
				resp.EpisodeID = uuid.MustParse(episodeID)
			}
			if variantName, ok := raw["variant_name"].(string); ok {
				resp.VariantName = variantName
			}
			if usage, ok := raw["usage"].(map[string]interface{}); ok {
				if inputTokens, ok := usage["input_tokens"].(float64); ok {
					resp.Usage.InputTokens = int(inputTokens)
				}
				if outputTokens, ok := usage["output_tokens"].(float64); ok {
					resp.Usage.OutputTokens = int(outputTokens)
				}
			}
			if finishReason, ok := raw["finish_reason"].(string); ok {
				fr := types.FinishReason(finishReason)
				resp.FinishReason = &fr
			}
			if originalResponse, ok := raw["original_response"].(string); ok {
				resp.OriginalResponse = &originalResponse
			}

			// Parse content blocks
			resp.Content = make([]types.ContentBlock, len(contentList))
			for i, block := range contentList {
				if blockMap, ok := block.(map[string]interface{}); ok {
					contentBlock, err := parseContentBlock(blockMap)
					if err != nil {
						return nil, fmt.Errorf("failed to parse content block %d: %w", i, err)
					}
					resp.Content[i] = contentBlock
				}
			}

			return &resp, nil
		}
	}

	if output, hasOutput := raw["output"]; hasOutput {
		if _, ok := output.(map[string]interface{}); ok {
			var resp types.JsonInferenceResponse
			if err := json.Unmarshal(data, &resp); err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
			}
			return &resp, nil
		}
	}

	return nil, fmt.Errorf("unable to determine response type")
}

// parseContentBlock parses a content block from a map
func parseContentBlock(block map[string]interface{}) (types.ContentBlock, error) {
	blockType, ok := block["type"].(string)
	if !ok {
		return nil, fmt.Errorf("content block missing type field")
	}

	switch blockType {
	case "text":
		text := &types.Text{Type: blockType}
		if textVal, ok := block["text"].(string); ok {
			text.Text = &textVal
		}
		if args, ok := block["arguments"]; ok {
			text.Arguments = args
		}
		return text, nil

	case "raw_text":
		value, ok := block["value"].(string)
		if !ok {
			return nil, fmt.Errorf("raw_text block missing value field")
		}
		return types.NewRawText(value), nil

	case "image":
		if data, hasData := block["data"].(string); hasData {
			mimeType, ok := block["mime_type"].(string)
			if !ok {
				return nil, fmt.Errorf("image block missing mime_type field")
			}
			return types.NewImageBase64(data, mimeType), nil
		}
		if url, hasURL := block["url"].(string); hasURL {
			img := types.NewImageURL(url)
			if mimeType, ok := block["mime_type"].(string); ok {
				img.MimeType = &mimeType
			}
			return img, nil
		}
		return nil, fmt.Errorf("image block missing data or url field")

	case "file":
		if data, hasData := block["data"].(string); hasData {
			mimeType, ok := block["mime_type"].(string)
			if !ok {
				return nil, fmt.Errorf("file block missing mime_type field")
			}
			return types.NewFileBase64(data, mimeType), nil
		}
		if url, hasURL := block["url"].(string); hasURL {
			return types.NewFileURL(url), nil
		}
		return nil, fmt.Errorf("file block missing data or url field")

	case "tool_call":
		id, ok := block["id"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_call block missing id field")
		}
		rawArgs, ok := block["raw_arguments"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_call block missing raw_arguments field")
		}
		rawName, ok := block["raw_name"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_call block missing raw_name field")
		}

		toolCall := types.NewToolCall(id, rawArgs, rawName)
		if args, ok := block["arguments"].(map[string]interface{}); ok {
			toolCall.Arguments = args
		}
		if name, ok := block["name"].(string); ok {
			toolCall.Name = &name
		}
		return toolCall, nil

	case "thought":
		thought := types.NewThought("")
		if text, ok := block["text"].(string); ok {
			thought.Text = &text
		}
		if sig, ok := block["signature"].(string); ok {
			thought.Signature = &sig
		}
		return thought, nil

	case "tool_result":
		name, ok := block["name"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_result block missing name field")
		}
		result, ok := block["result"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_result block missing result field")
		}
		id, ok := block["id"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_result block missing id field")
		}
		return types.NewToolResult(name, result, id), nil

	case "unknown":
		data, ok := block["data"]
		if !ok {
			return nil, fmt.Errorf("unknown block missing data field")
		}
		unknown := types.NewUnknownContentBlock(data)
		if provider, ok := block["model_provider_name"].(string); ok {
			unknown.ModelProviderName = &provider
		}
		return unknown, nil

	default:
		return nil, fmt.Errorf("unknown content block type: %s", blockType)
	}
}

// parseInferenceChunk parses an inference chunk from JSON
func parseInferenceChunk(data []byte) (types.InferenceChunk, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
	}

	if content, hasContent := raw["content"]; hasContent {
		if contentList, ok := content.([]interface{}); ok {
			// Parse as ChatChunk but handle content separately
			var chunk types.ChatChunk

			// Parse basic fields
			if inferenceID, ok := raw["inference_id"].(string); ok {
				chunk.InferenceID = uuid.MustParse(inferenceID)
			}
			if episodeID, ok := raw["episode_id"].(string); ok {
				chunk.EpisodeID = uuid.MustParse(episodeID)
			}
			if variantName, ok := raw["variant_name"].(string); ok {
				chunk.VariantName = variantName
			}
			if usage, ok := raw["usage"].(map[string]interface{}); ok {
				chunk.Usage = &types.Usage{}
				if inputTokens, ok := usage["input_tokens"].(float64); ok {
					chunk.Usage.InputTokens = int(inputTokens)
				}
				if outputTokens, ok := usage["output_tokens"].(float64); ok {
					chunk.Usage.OutputTokens = int(outputTokens)
				}
			}
			if finishReason, ok := raw["finish_reason"].(string); ok {
				fr := types.FinishReason(finishReason)
				chunk.FinishReason = &fr
			}

			// Parse content block chunks
			chunk.Content = make([]types.ContentBlockChunk, len(contentList))
			for i, block := range contentList {
				if blockMap, ok := block.(map[string]interface{}); ok {
					contentChunk, err := parseContentBlockChunk(blockMap)
					if err != nil {
						return nil, fmt.Errorf("failed to parse content block chunk %d: %w", i, err)
					}
					chunk.Content[i] = contentChunk
				}
			}

			return &chunk, nil
		}
	}

	if _, hasRaw := raw["raw"]; hasRaw {
		var chunk types.JsonChunk
		if err := json.Unmarshal(data, &chunk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON chunk: %w", err)
		}
		return &chunk, nil
	}

	return nil, fmt.Errorf("unable to determine chunk type")
}

// parseContentBlockChunk parses a content block chunk from a map
func parseContentBlockChunk(block map[string]interface{}) (types.ContentBlockChunk, error) {
	blockType, ok := block["type"].(string)
	if !ok {
		return nil, fmt.Errorf("content block chunk missing type field")
	}

	id, ok := block["id"].(string)
	if !ok {
		return nil, fmt.Errorf("content block chunk missing id field")
	}

	switch blockType {
	case "text":
		text, ok := block["text"].(string)
		if !ok {
			return nil, fmt.Errorf("text chunk missing text field")
		}
		return &types.TextChunk{
			ID:   id,
			Text: text,
			Type: blockType,
		}, nil

	case "tool_call":
		rawArgs, ok := block["raw_arguments"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_call chunk missing raw_arguments field")
		}
		rawName, ok := block["raw_name"].(string)
		if !ok {
			return nil, fmt.Errorf("tool_call chunk missing raw_name field")
		}
		return &types.ToolCallChunk{
			ID:           id,
			RawArguments: rawArgs,
			RawName:      rawName,
			Type:         blockType,
		}, nil

	case "thought":
		text, ok := block["text"].(string)
		if !ok {
			return nil, fmt.Errorf("thought chunk missing text field")
		}
		chunk := &types.ThoughtChunk{
			ID:   id,
			Text: text,
			Type: blockType,
		}
		if sig, ok := block["signature"].(string); ok {
			chunk.Signature = &sig
		}
		return chunk, nil

	default:
		return nil, fmt.Errorf("unknown content block chunk type: %s", blockType)
	}
}
