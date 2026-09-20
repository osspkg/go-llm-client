package batches

//go:generate easyjson -all types.go

// Request creates a batch.
type Request struct {
	// InputFileID identifies the batch input file.
	InputFileID string `json:"input_file_id"`
	// Endpoint identifies the API endpoint used by a batch.
	Endpoint string `json:"endpoint"`
	// CompletionWindow selects the provider processing window for a batch.
	CompletionWindow string `json:"completion_window"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Batch describes a batch job.
type Batch struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Endpoint identifies the API endpoint used by a batch.
	Endpoint string `json:"endpoint"`
	// Errors contains the provider value used for the `Errors` field.
	Errors []Error `json:"errors,omitempty"`
	// InputFileID identifies the batch input file.
	InputFileID string `json:"input_file_id"`
	// OutputFileID identifies the batch output file.
	OutputFileID string `json:"output_file_id,omitempty"`
	// ErrorFileID identifies the file containing batch errors.
	ErrorFileID string `json:"error_file_id,omitempty"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// CompletedAt contains the provider value used for the `CompletedAt` field.
	CompletedAt int64 `json:"completed_at,omitempty"`
}

// Error is a batch error.
type Error struct {
	// Code contains the provider error code.
	Code string `json:"code,omitempty"`
	// Message contains the provider error message.
	Message string `json:"message,omitempty"`
}

// ListResponse is a cursor page of batches.
type ListResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Batch `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}
