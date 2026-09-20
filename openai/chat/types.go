package chat

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a chat completion.
type Request struct {
	Model          string            `json:"model"`
	Messages       []Message         `json:"messages"`
	Stream         bool              `json:"stream,omitempty"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
	Temperature    float64           `json:"temperature,omitempty"`
	TopP           float64           `json:"top_p,omitempty"`
	Tools          []Tool            `json:"tools,omitempty"`
	ToolChoice     json.RawMessage   `json:"tool_choice,omitempty"`
	ResponseFormat json.RawMessage   `json:"response_format,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// Message is a chat message.
type Message struct {
	Role       string          `json:"role"`
	Content    json.RawMessage `json:"content,omitempty"`
	Name       string          `json:"name,omitempty"`
	ToolCalls  []ToolCall      `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
}

// Tool describes a function tool.
type Tool struct {
	Type     string       `json:"type"`
	Function FunctionTool `json:"function"`
}

// FunctionTool describes a callable function.
type FunctionTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall is a model-generated tool call.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Index    int          `json:"index,omitempty"`
	Function FunctionCall `json:"function"`
}

// FunctionCall contains generated function arguments.
type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// Response is a chat completion.
type Response struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	Created     int64    `json:"created"`
	Model       string   `json:"model"`
	Choices     []Choice `json:"choices"`
	Usage       *Usage   `json:"usage,omitempty"`
	ServiceTier string   `json:"service_tier,omitempty"`
}

// Choice is a completion choice.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        Message `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

// Usage contains token accounting.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk is one streamed chat event.
type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}
