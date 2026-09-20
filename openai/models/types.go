package models

//go:generate easyjson -all types.go

// Model describes a model available to the account.
type Model struct {
	// ID is the model identifier accepted by API requests.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Created is the Unix registration time.
	Created int64 `json:"created"`
	// OwnedBy identifies the organization that owns the model.
	OwnedBy string `json:"owned_by"`
}

// ListResponse is a model list.
type ListResponse struct {
	// Object identifies the provider list kind.
	Object string `json:"object"`
	// Data contains the available models.
	Data []Model `json:"data"`
}

// DeleteResponse describes a deleted model.
type DeleteResponse struct {
	// ID identifies the deleted model.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Deleted confirms that deletion succeeded.
	Deleted bool `json:"deleted"`
}
