package stateful

import "encoding/json"

//go:generate easyjson -all types.go

// Request is shared by stateful resource creation operations.
type Request struct {
	Model        string            `json:"model,omitempty"`
	AssistantID  string            `json:"assistant_id,omitempty"`
	ThreadID     string            `json:"thread_id,omitempty"`
	Instructions string            `json:"instructions,omitempty"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Data         json.RawMessage   `json:"data,omitempty"`
}

// Resource is a typed common stateful resource envelope.
type Resource struct {
	ID          string            `json:"id"`
	Object      string            `json:"object"`
	CreatedAt   int64             `json:"created_at,omitempty"`
	Status      string            `json:"status,omitempty"`
	Model       string            `json:"model,omitempty"`
	AssistantID string            `json:"assistant_id,omitempty"`
	ThreadID    string            `json:"thread_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListResponse is a cursor list of stateful resources.
type ListResponse struct {
	Object  string     `json:"object"`
	Data    []Resource `json:"data"`
	FirstID string     `json:"first_id,omitempty"`
	LastID  string     `json:"last_id,omitempty"`
	HasMore bool       `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
