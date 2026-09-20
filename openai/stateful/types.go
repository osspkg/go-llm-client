package stateful

import "encoding/json"

//go:generate easyjson -all types.go

// Request is shared by stateful resource creation operations.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model,omitempty"`
	// AssistantID contains the provider value used for the `AssistantID` field.
	AssistantID string `json:"assistant_id,omitempty"`
	// ThreadID contains the provider value used for the `ThreadID` field.
	ThreadID string `json:"thread_id,omitempty"`
	// Instructions contains system-level guidance for response generation.
	Instructions string `json:"instructions,omitempty"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Description explains a tool or resource so the model can use it correctly.
	Description string `json:"description,omitempty"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Data contains the returned records.
	Data json.RawMessage `json:"data,omitempty"`
}

// Resource is a typed common stateful resource envelope.
type Resource struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at,omitempty"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status,omitempty"`
	// Model selects the model used for the operation.
	Model string `json:"model,omitempty"`
	// AssistantID contains the provider value used for the `AssistantID` field.
	AssistantID string `json:"assistant_id,omitempty"`
	// ThreadID contains the provider value used for the `ThreadID` field.
	ThreadID string `json:"thread_id,omitempty"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListResponse is a cursor list of stateful resources.
type ListResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Resource `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
}
