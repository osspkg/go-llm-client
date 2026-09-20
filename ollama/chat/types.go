package chat

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a chat response.
type Request struct {
	Model     string          `json:"model"`
	Messages  []Message       `json:"messages"`
	Stream    bool            `json:"stream,omitempty"`
	Tools     []Tool          `json:"tools,omitempty"`
	Think     json.RawMessage `json:"think,omitempty"`
	Format    json.RawMessage `json:"format,omitempty"`
	Options   json.RawMessage `json:"options,omitempty"`
	KeepAlive string          `json:"keep_alive,omitempty"`
}

// Message is an Ollama chat message.
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	Thinking  string     `json:"thinking,omitempty"`
	Images    []string   `json:"images,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Tool is an Ollama tool declaration.
type Tool struct {
	Type     string       `json:"type"`
	Function FunctionTool `json:"function"`
}

// FunctionTool describes a tool.
type FunctionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall is a model tool call.
type ToolCall struct {
	Function FunctionCall `json:"function"`
}

// FunctionCall contains tool arguments.
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// Response is one chat response event.
type Response struct {
	Model           string  `json:"model"`
	CreatedAt       string  `json:"created_at"`
	Message         Message `json:"message"`
	Done            bool    `json:"done"`
	TotalDuration   int64   `json:"total_duration,omitempty"`
	LoadDuration    int64   `json:"load_duration,omitempty"`
	PromptEvalCount int     `json:"prompt_eval_count,omitempty"`
	EvalCount       int     `json:"eval_count,omitempty"`
	DoneReason      string  `json:"done_reason,omitempty"`
}
