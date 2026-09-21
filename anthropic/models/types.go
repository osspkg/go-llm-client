// Package models provides typed models API operations.
package models

//go:generate easyjson -all types.go

// Model describes a model available to the Anthropic account.
type Model struct {
	// ID is the stable model identifier accepted by message requests.
	ID string `json:"id"`
	// Type identifies the provider resource type.
	Type string `json:"type"`
	// DisplayName is the human-readable model name.
	DisplayName string `json:"display_name"`
	// CreatedAt is the model release timestamp in ISO-8601 form.
	CreatedAt string `json:"created_at"`
	// MaxInputTokens is the maximum input context accepted by the model.
	MaxInputTokens int `json:"max_input_tokens,omitempty"`
	// MaxTokens is the maximum output token count supported by the model.
	MaxTokens int `json:"max_tokens,omitempty"`
}

// ListResponse is a cursor page of Anthropic models.
type ListResponse struct {
	// Data contains models returned for the current page.
	Data []Model `json:"data"`
	// FirstID is the first model ID in the page.
	FirstID string `json:"first_id,omitempty"`
	// LastID is the last model ID in the page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}
