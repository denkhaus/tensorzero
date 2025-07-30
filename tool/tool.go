package tool

// ToolCallChunk represents streaming tool call chunk
type ToolCallChunk struct {
	ID           string `json:"id"`
	RawArguments string `json:"raw_arguments"`
	RawName      string `json:"raw_name"`
	Type         string `json:"type"`
}

func (tcc *ToolCallChunk) GetType() string { return tcc.Type }
func (tcc *ToolCallChunk) GetID() string   { return tcc.ID }

// ToolChoice represents tool choice options
type ToolChoice interface{}

// Tool represents a tool definition that can be called by AI models during inference.
// Tools enable models to perform actions, retrieve information, or interact with external systems.
type Tool struct {
	// Description provides a clear explanation of what this tool does and when to use it.
	// This helps the model understand the tool's purpose and make appropriate calls.
	Description string `json:"description"`

	// Parameters defines the JSON schema for the tool's input parameters.
	// This schema validates the arguments the model provides when calling the tool.
	Parameters interface{} `json:"parameters"`

	// Name is the unique identifier for this tool. Must be unique within the function scope.
	// The model will use this name when requesting to call the tool.
	Name string `json:"name"`

	// Strict enforces strict parameter validation according to the schema.
	// When true, ensures the model provides exactly the required parameters in the correct format.
	Strict bool `json:"strict"`
}

// ToolParams represents tool-related parameters for an inference request.
// This configures which tools are available and how the model should use them.
type ToolParams struct {
	// ToolsAvailable lists all tools that are available for the model to call during inference.
	// These can include both predefined tools from configuration and additional dynamic tools.
	ToolsAvailable []Tool `json:"tools_available"`

	// ToolChoice specifies the strategy for tool selection during inference.
	// Options include "none", "auto", "required", or specific tool selection.
	ToolChoice string `json:"tool_choice"`

	// ParallelToolCalls, when true, allows the model to request multiple tool calls
	// in a single response. Only supported by certain model providers.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`
}

// ToolCodeContent represents content of type "tool_code"
type ToolCodeContent struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Type     string `json:"type"`
}

func (tcc *ToolCodeContent) GetType() string { return tcc.Type }
func (tcc *ToolCodeContent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":     tcc.Type,
		"language": tcc.Language,
		"code":     tcc.Code,
	}
}

// ToolOutputContent represents content of type "tool_output"
type ToolOutputContent struct {
	Output string `json:"output"`
	Type   string `json:"type"`
}

func (toc *ToolOutputContent) GetType() string { return toc.Type }
func (toc *ToolOutputContent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   toc.Type,
		"output": toc.Output,
	}
}

// ToolResult represents the result of a tool call execution.
// This contains the output from executing a tool that was called by the model.
type ToolResult struct {
	// Name is the name of the tool that was executed, matching the tool definition.
	Name string `json:"name"`

	// Result contains the actual output or return value from executing the tool.
	// This is typically a string representation of the tool's output.
	Result string `json:"result"`

	// ID is the unique identifier for this tool call, matching the ID from the tool call request.
	// This links the result back to the specific tool call that generated it.
	ID string `json:"id"`

	// Type indicates the content type, typically "tool_result" for tool execution results.
	Type string `json:"type"`
}

func NewToolResult(name, result, id string) *ToolResult {
	return &ToolResult{
		Name:   name,
		Result: result,
		ID:     id,
		Type:   "tool_result",
	}
}

func (tr *ToolResult) GetType() string {
	return tr.Type
}

func (tr *ToolResult) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   tr.Type,
		"name":   tr.Name,
		"result": tr.Result,
		"id":     tr.ID,
	}
}
