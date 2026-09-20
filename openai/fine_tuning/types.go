package fine_tuning //nolint:revive // package path follows the provider domain name.

import "encoding/json"

//go:generate easyjson -all types.go

// CreateRequest starts a fine-tuning job.
type CreateRequest struct {
	Model           string            `json:"model"`
	TrainingFile    string            `json:"training_file"`
	ValidationFile  string            `json:"validation_file,omitempty"`
	Hyperparameters json.RawMessage   `json:"hyperparameters,omitempty"`
	Suffix          string            `json:"suffix,omitempty"`
	Integrations    json.RawMessage   `json:"integrations,omitempty"`
	Seed            int               `json:"seed,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// Job describes a fine-tuning job.
type Job struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	Model          string `json:"model"`
	FineTunedModel string `json:"fine_tuned_model,omitempty"`
	Status         string `json:"status"`
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file,omitempty"`
	CreatedAt      int64  `json:"created_at"`
	FinishedAt     int64  `json:"finished_at,omitempty"`
	Epochs         int    `json:"epochs,omitempty"`
	Error          *Error `json:"error,omitempty"`
}

// Error describes a fine-tuning failure.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   string `json:"param,omitempty"`
}

// ListResponse lists fine-tuning jobs.
type ListResponse struct {
	Object  string `json:"object"`
	Data    []Job  `json:"data"`
	HasMore bool   `json:"has_more"`
	FirstID string `json:"first_id,omitempty"`
	LastID  string `json:"last_id,omitempty"`
}

// Event is one fine-tuning event.
type Event struct {
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
}
