/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package models

//go:generate easyjson -all types.go

// Model is an installed or running model.
type Model struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// Model selects the model used for the operation.
	Model string `json:"model,omitempty"`
	// ModifiedAt contains the model modification timestamp.
	ModifiedAt string `json:"modified_at,omitempty"`
	// Size contains the model or file size in bytes.
	Size int64 `json:"size,omitempty"`
	// Digest identifies an immutable blob or model content by SHA-256 digest.
	Digest string `json:"digest,omitempty"`
	// Details contains model format, family, and quantization metadata.
	Details Details `json:"details,omitempty"`
	// ExpiresAt contains the time a running model will be unloaded.
	ExpiresAt string `json:"expires_at,omitempty"`
	// SizeVRAM contains the VRAM occupied by the running model in bytes.
	SizeVRAM int64 `json:"size_vram,omitempty"`
}

// Details contains model metadata.
type Details struct {
	// ParentModel identifies the model from which this model was derived.
	ParentModel string `json:"parent_model,omitempty"`
	// Format selects the model or response format.
	Format string `json:"format,omitempty"`
	// Family identifies the model family.
	Family string `json:"family,omitempty"`
	// Families lists model families supported by the model.
	Families []string `json:"families,omitempty"`
	// ParameterSize describes the model parameter count.
	ParameterSize string `json:"parameter_size,omitempty"`
	// QuantizationLevel describes the model quantization level.
	QuantizationLevel string `json:"quantization_level,omitempty"`
}

// ListResponse lists local models.
type ListResponse struct {
	// Models contains the available or running models.
	Models []Model `json:"models"`
}

// Request is shared by model lifecycle operations.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// From identifies the source model used to create a model.
	From string `json:"from,omitempty"`
	// Files maps model file names to uploaded blob SHA-256 digests.
	Files map[string]string `json:"files,omitempty"`
	// Template contains the prompt template used by a model.
	Template string `json:"template,omitempty"`
	// System contains the default system prompt for a model.
	System string `json:"system,omitempty"`
	// License contains one or more model license texts.
	License []string `json:"license,omitempty"`
	// Stream requests incremental events instead of one buffered response.
	Stream bool `json:"stream,omitempty"`
	// Insecure allows insecure model creation when the provider supports it.
	Insecure bool `json:"insecure,omitempty"`
	// Quantize selects the quantization format for model creation.
	Quantize string `json:"quantize,omitempty"`
}

// Progress is a model lifecycle event.
type Progress struct {
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status"`
	// Digest identifies an immutable blob or model content by SHA-256 digest.
	Digest string `json:"digest,omitempty"`
	// Total contains the total number of bytes expected.
	Total int64 `json:"total,omitempty"`
	// Completed contains the number of bytes processed so far.
	Completed int64 `json:"completed,omitempty"`
}

// ShowResponse contains model information.
type ShowResponse struct {
	// License contains one or more model license texts.
	License string `json:"license,omitempty"`
	// Modelfile contains the provider value used for the `Modelfile` field.
	Modelfile string `json:"modelfile,omitempty"`
	// Parameters contains the JSON Schema or arguments accepted by a function.
	Parameters string `json:"parameters,omitempty"`
	// Template contains the prompt template used by a model.
	Template string `json:"template,omitempty"`
	// System contains the default system prompt for a model.
	System string `json:"system,omitempty"`
	// Details contains model format, family, and quantization metadata.
	Details Details `json:"details,omitempty"`
	// ModelInfo contains provider-specific model metadata.
	ModelInfo map[string]string `json:"model_info,omitempty"`
}

// VersionResponse contains the Ollama version.
type VersionResponse struct {
	// Version contains the Ollama server version.
	Version string `json:"version"`
}
