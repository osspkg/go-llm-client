package responses

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a response.
type Request struct {
	Model           string            `json:"model"`
	Input           json.RawMessage   `json:"input,omitempty"`
	Instructions    string            `json:"instructions,omitempty"`
	Stream          bool              `json:"stream,omitempty"`
	MaxOutputTokens int               `json:"max_output_tokens,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Tools           []Tool            `json:"tools,omitempty"`
}

// Tool describes a Responses tool declaration.
type Tool struct {
	Type        string          `json:"type"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// Response is the typed response object.
type Response struct {
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	CreatedAt  int64        `json:"created_at"`
	Model      string       `json:"model"`
	Status     string       `json:"status"`
	Output     []OutputItem `json:"output,omitempty"`
	OutputText string       `json:"output_text,omitempty"`
	Usage      *Usage       `json:"usage,omitempty"`
	Error      *Error       `json:"error,omitempty"`
}

// OutputItem is a response output item.
type OutputItem struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Role    string          `json:"role,omitempty"`
	Status  string          `json:"status,omitempty"`
	Content []OutputContent `json:"content,omitempty"`
	Text    string          `json:"text,omitempty"`
}

// OutputContent is a typed output content item.
type OutputContent struct {
	Type        string       `json:"type"`
	Text        string       `json:"text,omitempty"`
	Refusal     string       `json:"refusal,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

// Annotation is an output annotation.
type Annotation struct {
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
}

// Usage contains token accounting.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Error is a provider response error.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// StreamEvent is one Responses SSE event.
type StreamEvent struct {
	Type           string    `json:"type"`
	SequenceNumber int       `json:"sequence_number,omitempty"`
	ResponseID     string    `json:"response_id,omitempty"`
	Delta          string    `json:"delta,omitempty"`
	Text           string    `json:"text,omitempty"`
	Response       *Response `json:"response,omitempty"`
}

// DeleteResponse confirms deletion of a stored response.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// InputItem is an input item returned by the Responses API.
type InputItem struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Role    string          `json:"role,omitempty"`
	Status  string          `json:"status,omitempty"`
	Content json.RawMessage `json:"content,omitempty"`
}

// InputItemsResponse is a cursor page of response input items.
type InputItemsResponse struct {
	Object  string      `json:"object"`
	Data    []InputItem `json:"data"`
	FirstID string      `json:"first_id,omitempty"`
	LastID  string      `json:"last_id,omitempty"`
	HasMore bool        `json:"has_more"`
}
