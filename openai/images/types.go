package images

import "io"

//go:generate easyjson -all types.go

// Request generates or edits an image.
type Request struct {
	// Model optionally selects the image model.
	Model string `json:"model,omitempty"`
	// Prompt describes the requested image.
	Prompt string `json:"prompt"`
	// N requests this many images.
	N int `json:"n,omitempty"`
	// Size selects the generated image dimensions.
	Size string `json:"size,omitempty"`
	// Quality selects the supported rendering quality.
	Quality string `json:"quality,omitempty"`
	// ResponseFormat selects URL or base64 image output.
	ResponseFormat string `json:"response_format,omitempty"`
	// Image supplies an optional source image for compatible requests.
	Image string `json:"image,omitempty"`
	// Stream requests incremental SSE events.
	Stream bool `json:"stream,omitempty"`
}

// EditInput is a multipart image edit request.
//
//easyjson:skip
type EditInput struct {
	// Model optionally selects the image model.
	Model string
	// Prompt describes the desired edit.
	Prompt string
	// Filename names the uploaded image part.
	Filename string
	// Image provides the source image bytes.
	Image io.Reader
	// Mask optionally constrains the editable area.
	Mask io.Reader
	// N requests this many edited images.
	N int
	// Size selects the generated image dimensions.
	Size string
	// Stream requests incremental SSE events.
	Stream bool
}

// VariationInput is a multipart image variation request.
//
//easyjson:skip
type VariationInput struct {
	// Model optionally selects the image model.
	Model string
	// Filename names the uploaded image part.
	Filename string
	// Image provides the source image bytes.
	Image io.Reader
	// N requests this many variations.
	N int
	// Size selects the generated image dimensions.
	Size string
}

// Response contains generated images.
type Response struct {
	// Created is the Unix response creation time.
	Created int64 `json:"created"`
	// Data contains generated image results.
	Data []Image `json:"data"`
}

// Image is an image URL or base64 payload.
type Image struct {
	// URL is a provider-hosted image location when requested.
	URL string `json:"url,omitempty"`
	// B64JSON contains an image encoded as base64 JSON data.
	B64JSON string `json:"b64_json,omitempty"`
	// RevisedPrompt contains the provider-normalized prompt.
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// StreamEvent is an SSE event emitted while generating or editing an image.
type StreamEvent struct {
	// Type identifies the image streaming event.
	Type string `json:"type"`
	// B64JSON contains the partial or final image data.
	B64JSON string `json:"b64_json,omitempty"`
	// PartialImageIndex identifies a progressive image update.
	PartialImageIndex int `json:"partial_image_index,omitempty"`
}
