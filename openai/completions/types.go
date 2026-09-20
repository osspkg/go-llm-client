package completions

//go:generate easyjson -all types.go

// Request creates a legacy text completion.
type Request struct {
	// Model selects the completion model.
	Model string `json:"model"`
	// Prompt supplies the text to complete.
	Prompt string `json:"prompt"`
	// MaxTokens limits generated tokens.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Stream requests SSE chunks instead of one response.
	Stream bool `json:"stream,omitempty"`
	// Temperature controls sampling randomness.
	Temperature float64 `json:"temperature,omitempty"`
}

// Response is a text completion response.
type Response struct {
	// ID identifies this completion.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Created is the Unix creation time.
	Created int64 `json:"created"`
	// Model identifies the model that produced the completion.
	Model string `json:"model"`
	// Choices contains generated alternatives.
	Choices []Choice `json:"choices"`
	// Usage reports consumed tokens.
	Usage Usage `json:"usage"`
}

// Choice is one text completion choice.
type Choice struct {
	// Text is the generated completion text.
	Text string `json:"text"`
	// Index is this choice's zero-based position.
	Index int `json:"index"`
	// FinishReason explains why generation stopped.
	FinishReason string `json:"finish_reason,omitempty"`
}

// Usage contains token accounting.
type Usage struct {
	// PromptTokens counts tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens counts generated tokens.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens is the total billable token count.
	TotalTokens int `json:"total_tokens"`
}
