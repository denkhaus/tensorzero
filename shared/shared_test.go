//go:build unit

package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	text := NewText("hello world")
	assert.Equal(t, "hello world", *text.Text)
	assert.Equal(t, "text", text.Type)
	assert.Equal(t, "text", text.GetType())
	assert.Equal(t, map[string]interface{}{"type": "text", "text": "hello world"}, text.ToMap())

	args := map[string]interface{}{"param1": "value1"}
	textWithArgs := NewTextWithArguments(args)
	assert.Nil(t, textWithArgs.Text)
	assert.Equal(t, args, textWithArgs.Arguments)
	assert.Equal(t, "text", textWithArgs.Type)
	assert.Equal(t, "text", textWithArgs.GetType())
	assert.Equal(t, map[string]interface{}{"type": "text", "arguments": args}, textWithArgs.ToMap())
}

func TestRawText(t *testing.T) {
	rawText := NewRawText("raw string")
	assert.Equal(t, "raw string", rawText.Value)
	assert.Equal(t, "raw_text", rawText.Type)
	assert.Equal(t, "raw_text", rawText.GetType())
	assert.Equal(t, map[string]interface{}{"type": "raw_text", "value": "raw string"}, rawText.ToMap())
}

func TestImageBase64(t *testing.T) {
	img := NewImageBase64("base64data", "image/png")
	assert.Equal(t, "base64data", img.Data)
	assert.Equal(t, "image/png", img.MimeType)
	assert.Equal(t, "image", img.Type)
	assert.Equal(t, "image", img.GetType())
	assert.Equal(t, map[string]interface{}{"type": "image", "data": "base64data", "mime_type": "image/png"}, img.ToMap())
}

func TestImageURL(t *testing.T) {
	img := NewImageURL("http://example.com/image.jpg")
	assert.Equal(t, "http://example.com/image.jpg", img.URL)
	assert.Nil(t, img.MimeType)
	assert.Equal(t, "image", img.Type)
	assert.Equal(t, "image", img.GetType())
	assert.Equal(t, map[string]interface{}{"type": "image", "url": "http://example.com/image.jpg"}, img.ToMap())

	imgWithMime := NewImageURLWithMimeType("http://example.com/image.gif", "image/gif")
	assert.Equal(t, "http://example.com/image.gif", imgWithMime.URL)
	assert.NotNil(t, imgWithMime.MimeType)
	assert.Equal(t, "image/gif", *imgWithMime.MimeType)
	assert.Equal(t, "image", imgWithMime.Type)
	assert.Equal(t, "image", imgWithMime.GetType())
	assert.Equal(t, map[string]interface{}{"type": "image", "url": "http://example.com/image.gif", "mime_type": "image/gif"}, imgWithMime.ToMap())
}

func TestFileBase64(t *testing.T) {
	file := NewFileBase64("base64filedata", "application/pdf")
	assert.Equal(t, "base64filedata", file.Data)
	assert.Equal(t, "application/pdf", file.MimeType)
	assert.Equal(t, "file", file.Type)
	assert.Equal(t, "file", file.GetType())
	assert.Equal(t, map[string]interface{}{"type": "file", "data": "base64filedata", "mime_type": "application/pdf"}, file.ToMap())
}

func TestFileURL(t *testing.T) {
	file := NewFileURL("http://example.com/document.docx")
	assert.Equal(t, "http://example.com/document.docx", file.URL)
	assert.Equal(t, "file", file.Type)
	assert.Equal(t, "file", file.GetType())
	assert.Equal(t, map[string]interface{}{"type": "file", "url": "http://example.com/document.docx"}, file.ToMap())
}

func TestToolCall(t *testing.T) {
	toolCall := NewToolCall("tool123", `{"arg1":"val1"}`, "my_tool")
	assert.Equal(t, "tool123", toolCall.ID)
	assert.Equal(t, `{"arg1":"val1"}`, toolCall.RawArguments)
	assert.Equal(t, "my_tool", toolCall.RawName)
	assert.Equal(t, "tool_call", toolCall.Type)
	assert.Equal(t, "tool_call", toolCall.GetType())

	// Test ToMap with optional fields
	toolCall.Arguments = map[string]interface{}{"arg1": "val1"}
	name := "processed_tool"
	toolCall.Name = &name
	expectedMap := map[string]interface{}{
		"type":          "tool_call",
		"id":            "tool123",
		"raw_arguments": `{"arg1":"val1"}`,
		"raw_name":      "my_tool",
		"arguments":     map[string]interface{}{"arg1": "val1"},
		"name":          "processed_tool",
	}
	assert.Equal(t, expectedMap, toolCall.ToMap())
}

func TestThought(t *testing.T) {
	thought := NewThought("thinking process")
	assert.Equal(t, "thinking process", *thought.Text)
	assert.Equal(t, "thought", thought.Type)
	assert.Equal(t, "thought", thought.GetType())
	assert.Equal(t, map[string]interface{}{"type": "thought", "text": "thinking process"}, thought.ToMap())

	signature := "sig123"
	thought.Signature = &signature
	assert.Equal(t, map[string]interface{}{"type": "thought", "text": "thinking process", "signature": "sig123"}, thought.ToMap())
}

func TestUnknownContentBlock(t *testing.T) {
	unknown := NewUnknownContentBlock(map[string]interface{}{"error": "unknown type"})
	assert.Equal(t, map[string]interface{}{"error": "unknown type"}, unknown.Data)
	assert.Equal(t, "unknown", unknown.Type)
	assert.Equal(t, "unknown", unknown.GetType())
	assert.Equal(t, map[string]interface{}{"type": "unknown", "data": map[string]interface{}{"error": "unknown type"}}, unknown.ToMap())

	modelProvider := "some_provider"
	unknown.ModelProviderName = &modelProvider
	assert.Equal(t, map[string]interface{}{"type": "unknown", "data": map[string]interface{}{"error": "unknown type"}, "model_provider_name": "some_provider"}, unknown.ToMap())
}

func TestMessage(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: []ContentBlock{NewText("Hello"), NewRawText("raw")},
	}
	assert.Equal(t, "user", msg.Role)
	assert.Len(t, msg.Content, 2)
	assert.Equal(t, "text", msg.Content[0].GetType())
	assert.Equal(t, "raw_text", msg.Content[1].GetType())
}

func TestTextChunk(t *testing.T) {
	chunk := TextChunk{ID: "chunk1", Text: "partial text", Type: "text_chunk"}
	assert.Equal(t, "chunk1", chunk.ID)
	assert.Equal(t, "partial text", chunk.Text)
	assert.Equal(t, "text_chunk", chunk.Type)
	assert.Equal(t, "text_chunk", chunk.GetType())
	assert.Equal(t, "chunk1", chunk.GetID())
}

func TestThoughtChunk(t *testing.T) {
	signature := "sig456"
	chunk := ThoughtChunk{ID: "thoughtchunk1", Text: "thinking...", Type: "thought_chunk", Signature: &signature}
	assert.Equal(t, "thoughtchunk1", chunk.ID)
	assert.Equal(t, "thinking...", chunk.Text)
	assert.Equal(t, "thought_chunk", chunk.Type)
	assert.Equal(t, "thought_chunk", chunk.GetType())
	assert.Equal(t, "thoughtchunk1", chunk.GetID())
	assert.Equal(t, &signature, chunk.Signature)
}

func TestVariantExtraBody(t *testing.T) {
	del := true
	body := VariantExtraBody{
		VariantName: "variant_a",
		Pointer:     "/path/to/value",
		Value:       "new_value",
		Delete:      &del,
	}
	assert.Equal(t, "variant_a", body.VariantName)
	assert.Equal(t, "/path/to/value", body.Pointer)
	assert.Equal(t, "new_value", body.Value)
	assert.True(t, *body.Delete)
}

func TestProviderExtraBody(t *testing.T) {
	body := ProviderExtraBody{
		ModelProviderName: "openai",
		Pointer:           "/config/param",
		Value:             123,
	}
	assert.Equal(t, "openai", body.ModelProviderName)
	assert.Equal(t, "/config/param", body.Pointer)
	assert.Equal(t, 123, body.Value)
	assert.Nil(t, body.Delete)
}

func TestTensorZeroError(t *testing.T) {
	err := &TensorZeroError{StatusCode: 400, Text: "Bad Request"}
	assert.Equal(t, 400, err.StatusCode)
	assert.Equal(t, "Bad Request", err.Text)
	assert.Equal(t, "TensorZeroError (status code 400): Bad Request", err.Error())
}

func TestTensorZeroInternalError(t *testing.T) {
	err := &TensorZeroInternalError{Message: "Internal Server Error"}
	assert.Equal(t, "Internal Server Error", err.Message)
	assert.Equal(t, "Internal Server Error", err.Error())
}

func TestOrderBy(t *testing.T) {
	orderBy := NewOrderByTimestamp("ASC")
	assert.Equal(t, "timestamp", orderBy.By)
	assert.Nil(t, orderBy.Name)
	assert.Equal(t, "ASC", orderBy.Direction)

	metricName := "accuracy"
	orderBy = NewOrderByMetric(metricName, "DESC")
	assert.Equal(t, "metric", orderBy.By)
	assert.NotNil(t, orderBy.Name)
	assert.Equal(t, metricName, *orderBy.Name)
	assert.Equal(t, "DESC", orderBy.Direction)
}
