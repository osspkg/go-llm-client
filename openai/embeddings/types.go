package embeddings

//go:generate easyjson -all types.go

// Request creates embeddings.
type Request struct {
	// Model selects the embedding model.
	Model string `json:"model"`
	// Input is the text to embed.
	Input string `json:"input"`
	// EncodingFormat selects the vector encoding returned by the provider.
	EncodingFormat string `json:"encoding_format,omitempty"`
	// Dimensions requests a reduced embedding dimension when supported.
	Dimensions int `json:"dimensions,omitempty"`
}

// Response is an embeddings response.
type Response struct {
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Data contains the vectors in input order.
	Data []Embedding `json:"data"`
	// Model identifies the model that produced the vectors.
	Model string `json:"model"`
	// Usage reports consumed tokens.
	Usage Usage `json:"usage"`
}

// Embedding is one vector.
type Embedding struct {
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Embedding is the numeric vector.
	Embedding []float64 `json:"embedding"`
	// Index maps this vector to its input position.
	Index int `json:"index"`
}

// Usage contains embedding token accounting.
type Usage struct {
	// PromptTokens counts tokens in the embedding input.
	PromptTokens int `json:"prompt_tokens"`
	// TotalTokens is the total billable token count.
	TotalTokens int `json:"total_tokens"`
}
