package images

import "io"

//go:generate easyjson -all types.go

// Request generates or edits an image.
type Request struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Image          string `json:"image,omitempty"`
}

// EditInput is a multipart image edit request.
//
//easyjson:skip
type EditInput struct {
	Model    string
	Prompt   string
	Filename string
	Image    io.Reader
	Mask     io.Reader
	N        int
	Size     string
}

// VariationInput is a multipart image variation request.
//
//easyjson:skip
type VariationInput struct {
	Model    string
	Filename string
	Image    io.Reader
	N        int
	Size     string
}

// Response contains generated images.
type Response struct {
	Created int64   `json:"created"`
	Data    []Image `json:"data"`
}

// Image is an image URL or base64 payload.
type Image struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}
