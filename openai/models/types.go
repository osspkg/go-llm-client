package models

//go:generate easyjson -all types.go

// Model describes a model available to the account.
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ListResponse is a model list.
type ListResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// DeleteResponse describes a deleted model.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
