// Package batches provides typed Anthropic Message Batches API operations.
package batches

import "go.osspkg.com/llm-client/anthropic/messages"

//go:generate easyjson -all types.go

// Request creates a batch of independent message requests.
type Request struct {
	// Requests contains message parameters and stable IDs for each batch item.
	Requests []RequestItem `json:"requests"`
	// Metadata contains caller-defined labels attached to the batch.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RequestItem associates one message request with a caller-defined result ID.
type RequestItem struct {
	// CustomID is returned unchanged in the corresponding JSONL result record.
	CustomID string `json:"custom_id"`
	// Params contains the model, messages, token limit, and other message parameters.
	Params messages.Request `json:"params"`
}

// Batch describes asynchronous batch processing state.
type Batch struct {
	// ID uniquely identifies the batch.
	ID string `json:"id"`
	// Type identifies the provider resource type.
	Type string `json:"type"`
	// ProcessingStatus is the current batch lifecycle state.
	ProcessingStatus string `json:"processing_status"`
	// RequestCounts contains submitted, processing, succeeded, and errored counts.
	RequestCounts RequestCounts `json:"request_counts,omitempty"`
	// CreatedAt is the batch creation timestamp.
	CreatedAt string `json:"created_at,omitempty"`
	// EndedAt is the completion timestamp, if processing ended.
	EndedAt string `json:"ended_at,omitempty"`
	// ExpiresAt is the timestamp after which results are no longer available.
	ExpiresAt string `json:"expires_at,omitempty"`
	// ResultsURL is the provider URL for JSONL batch results.
	ResultsURL string `json:"results_url,omitempty"`
}

// RequestCounts reports how many requests reached each batch lifecycle state.
type RequestCounts struct {
	// Processing is the number of requests currently being processed.
	Processing int `json:"processing,omitempty"`
	// Succeeded is the number of requests completed successfully.
	Succeeded int `json:"succeeded,omitempty"`
	// Errored is the number of requests that failed during processing.
	Errored int `json:"errored,omitempty"`
	// Canceled is the number of requests canceled before completion.
	Canceled int `json:"canceled,omitempty"`
	// Expired is the number of requests that expired before completion.
	Expired int `json:"expired,omitempty"`
}

// ListResponse is a cursor page of message batches.
type ListResponse struct {
	// Data contains batches returned for the current page.
	Data []Batch `json:"data"`
	// HasMore reports whether another page can be requested.
	HasMore bool `json:"has_more"`
	// FirstID is the first batch ID in the page.
	FirstID string `json:"first_id,omitempty"`
	// LastID is the last batch ID in the page.
	LastID string `json:"last_id,omitempty"`
}

// Result is one JSONL result record returned by a completed batch.
type Result struct {
	// CustomID links the result to the custom_id in the submitted request.
	CustomID string `json:"custom_id"`
	// Result contains the success or error variant for this request.
	Result ResultPayload `json:"result"`
}

// ResultPayload is the succeeded-or-errored result union for one batch item.
type ResultPayload struct {
	// Type identifies the result variant, normally succeeded or errored.
	Type string `json:"type"`
	// Message contains the generated message for a succeeded result.
	Message *messages.Response `json:"message,omitempty"`
	// Error contains the provider error for an errored result.
	Error *ErrorPayload `json:"error,omitempty"`
}

// ErrorPayload describes a failed request inside a completed batch.
type ErrorPayload struct {
	// Type identifies the provider error category.
	Type string `json:"type"`
	// Message explains why the individual request failed.
	Message string `json:"message"`
}
