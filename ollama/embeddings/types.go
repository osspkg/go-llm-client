package embeddings

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates one or more embeddings.
type Request struct {
	Model      string          `json:"model"`
	Input      json.RawMessage `json:"input"`
	Truncate   *bool           `json:"truncate,omitempty"`
	Dimensions int             `json:"dimensions,omitempty"`
	KeepAlive  string          `json:"keep_alive,omitempty"`
	Options    json.RawMessage `json:"options,omitempty"`
}

// Response contains embeddings.
type Response struct {
	Model           string      `json:"model"`
	Embeddings      [][]float64 `json:"embeddings"`
	TotalDuration   int64       `json:"total_duration,omitempty"`
	LoadDuration    int64       `json:"load_duration,omitempty"`
	PromptEvalCount int         `json:"prompt_eval_count,omitempty"`
}
