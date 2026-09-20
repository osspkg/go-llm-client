package generate

import "encoding/json"

//go:generate easyjson -all types.go

// Request generates a completion.
type Request struct {
	Model     string          `json:"model"`
	Prompt    string          `json:"prompt,omitempty"`
	Suffix    string          `json:"suffix,omitempty"`
	Stream    bool            `json:"stream,omitempty"`
	Format    json.RawMessage `json:"format,omitempty"`
	Options   json.RawMessage `json:"options,omitempty"`
	KeepAlive string          `json:"keep_alive,omitempty"`
}

// Response is one generate response or a complete non-stream response.
type Response struct {
	Model           string `json:"model"`
	CreatedAt       string `json:"created_at"`
	Response        string `json:"response,omitempty"`
	Done            bool   `json:"done"`
	TotalDuration   int64  `json:"total_duration,omitempty"`
	LoadDuration    int64  `json:"load_duration,omitempty"`
	PromptEvalCount int    `json:"prompt_eval_count,omitempty"`
	EvalCount       int    `json:"eval_count,omitempty"`
	DoneReason      string `json:"done_reason,omitempty"`
	Thinking        string `json:"thinking,omitempty"`
}
