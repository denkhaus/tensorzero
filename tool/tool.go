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

// Tool represents a tool definition
type Tool struct {
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
	Name        string      `json:"name"`
	Strict      bool        `json:"strict"`
}

// ToolParams represents tool parameters
type ToolParams struct {
	ToolsAvailable    []Tool `json:"tools_available"`
	ToolChoice        string `json:"tool_choice"`
	ParallelToolCalls *bool  `json:"parallel_tool_calls,omitempty"`
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

// ToolResult represents a tool result
type ToolResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
	ID     string `json:"id"`
	Type   string `json:"type"`
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
