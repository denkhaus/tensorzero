//go:build unit

package tool

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestToolCallChunk(t *testing.T) {
	chunk := ToolCallChunk{ID: "toolchunk1", RawArguments: "{}", RawName: "tool_name", Type: "tool_call_chunk"}
	assert.Equal(t, "toolchunk1", chunk.ID)
	assert.Equal(t, "{}", chunk.RawArguments)
	assert.Equal(t, "tool_name", chunk.RawName)
	assert.Equal(t, "tool_call_chunk", chunk.Type)
	assert.Equal(t, "tool_call_chunk", chunk.GetType())
	assert.Equal(t, "toolchunk1", chunk.GetID())
}

func TestTool(t *testing.T) {
	tool := Tool{
		Description: "A test tool",
		Parameters:  map[string]interface{}{"type": "object"},
		Name:        "test_tool",
		Strict:      true,
	}
	assert.Equal(t, "A test tool", tool.Description)
	assert.Equal(t, map[string]interface{}{"type": "object"}, tool.Parameters)
	assert.Equal(t, "test_tool", tool.Name)
	assert.True(t, tool.Strict)
}

func TestToolParams(t *testing.T) {
	tool1 := Tool{Name: "tool1"}
	tool2 := Tool{Name: "tool2"}
	parallel := true
	params := ToolParams{
		ToolsAvailable:    []Tool{tool1, tool2},
		ToolChoice:        "auto",
		ParallelToolCalls: &parallel,
	}
	assert.Len(t, params.ToolsAvailable, 2)
	assert.Equal(t, "auto", params.ToolChoice)
	assert.NotNil(t, params.ParallelToolCalls)
	assert.True(t, *params.ParallelToolCalls)
}

func TestToolResult(t *testing.T) {
	toolResult := NewToolResult("weather_tool", "sunny", "res456")
	assert.Equal(t, "weather_tool", toolResult.Name)
	assert.Equal(t, "sunny", toolResult.Result)
	assert.Equal(t, "res456", toolResult.ID)
	assert.Equal(t, "tool_result", toolResult.Type)
	assert.Equal(t, "tool_result", toolResult.GetType())
	assert.Equal(t, map[string]interface{}{"type": "tool_result", "name": "weather_tool", "result": "sunny", "id": "res456"}, toolResult.ToMap())
}
