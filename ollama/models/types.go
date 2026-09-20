package models

//go:generate easyjson -all types.go

// Model is an installed or running model.
type Model struct {
	Name       string  `json:"name"`
	Model      string  `json:"model,omitempty"`
	ModifiedAt string  `json:"modified_at,omitempty"`
	Size       int64   `json:"size,omitempty"`
	Digest     string  `json:"digest,omitempty"`
	Details    Details `json:"details,omitempty"`
	ExpiresAt  string  `json:"expires_at,omitempty"`
	SizeVRAM   int64   `json:"size_vram,omitempty"`
}

// Details contains model metadata.
type Details struct {
	ParentModel       string   `json:"parent_model,omitempty"`
	Format            string   `json:"format,omitempty"`
	Family            string   `json:"family,omitempty"`
	Families          []string `json:"families,omitempty"`
	ParameterSize     string   `json:"parameter_size,omitempty"`
	QuantizationLevel string   `json:"quantization_level,omitempty"`
}

// ListResponse lists local models.
type ListResponse struct {
	Models []Model `json:"models"`
}

// Request is shared by model lifecycle operations.
type Request struct {
	Model    string            `json:"model"`
	From     string            `json:"from,omitempty"`
	Files    map[string]string `json:"files,omitempty"`
	Template string            `json:"template,omitempty"`
	System   string            `json:"system,omitempty"`
	License  []string          `json:"license,omitempty"`
	Stream   bool              `json:"stream,omitempty"`
	Insecure bool              `json:"insecure,omitempty"`
	Quantize string            `json:"quantize,omitempty"`
}

// Progress is a model lifecycle event.
type Progress struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// ShowResponse contains model information.
type ShowResponse struct {
	License    string            `json:"license,omitempty"`
	Modelfile  string            `json:"modelfile,omitempty"`
	Parameters string            `json:"parameters,omitempty"`
	Template   string            `json:"template,omitempty"`
	System     string            `json:"system,omitempty"`
	Details    Details           `json:"details,omitempty"`
	ModelInfo  map[string]string `json:"model_info,omitempty"`
}

// VersionResponse contains the Ollama version.
type VersionResponse struct {
	Version string `json:"version"`
}
