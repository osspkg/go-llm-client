package batches

//go:generate easyjson -all types.go

// Request creates a batch.
type Request struct {
	InputFileID      string            `json:"input_file_id"`
	Endpoint         string            `json:"endpoint"`
	CompletionWindow string            `json:"completion_window"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// Batch describes a batch job.
type Batch struct {
	ID           string  `json:"id"`
	Object       string  `json:"object"`
	Endpoint     string  `json:"endpoint"`
	Errors       []Error `json:"errors,omitempty"`
	InputFileID  string  `json:"input_file_id"`
	OutputFileID string  `json:"output_file_id,omitempty"`
	ErrorFileID  string  `json:"error_file_id,omitempty"`
	Status       string  `json:"status"`
	CreatedAt    int64   `json:"created_at"`
	CompletedAt  int64   `json:"completed_at,omitempty"`
}

// Error is a batch error.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ListResponse is a cursor page of batches.
type ListResponse struct {
	Object  string  `json:"object"`
	Data    []Batch `json:"data"`
	FirstID string  `json:"first_id,omitempty"`
	LastID  string  `json:"last_id,omitempty"`
	HasMore bool    `json:"has_more"`
}
