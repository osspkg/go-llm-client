// Package embeddings provides typed embeddings API operations.
package embeddings

import (
	"encoding/json"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

//go:generate easyjson -all types.go

// Request contains the text or multimodal content to embed and normalization options.
type Request struct {
	// Model selects a model in router mode.
	Model string `json:"model,omitempty"`
	// Content is the text, token input, or multimodal prompt accepted by the native endpoint.
	Content json.RawMessage `json:"content"`
	// Input is the OpenAI-compatible input shape accepted by native /embeddings servers.
	Input json.RawMessage `json:"input,omitempty"`
	// EncodingFormat selects the representation requested by compatible servers.
	EncodingFormat string `json:"encoding_format,omitempty"`
	// Normalization selects the embedding normalization mode; -1 disables normalization.
	Normalization int `json:"embd_normalize,omitempty"`
}

// Item is one native embedding result.
type Item struct {
	// Index is the zero-based input item index.
	Index int `json:"index"`
	// Embedding contains pooled values or one vector per input token when pooling is disabled.
	Embedding [][]float64 `json:"embedding"`
}

// Response is the native endpoint's array of embedding results.
type Response []Item

// MarshalJSON encodes the native response array.
func (response Response) MarshalJSON() ([]byte, error) {
	type responseAlias []Item
	return json.Marshal(responseAlias(response))
}

// UnmarshalJSON decodes the native response array.
func (response *Response) UnmarshalJSON(data []byte) error {
	type responseAlias []Item
	var value responseAlias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*response = Response(value)
	return nil
}

// MarshalEasyJSON implements easyjson.Marshaler for the response array.
func (response Response) MarshalEasyJSON(writer *jwriter.Writer) {
	data, err := response.MarshalJSON()
	writer.Raw(data, err)
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler for the response array.
func (response *Response) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	data := lexer.Raw()
	if lexer.Ok() {
		lexer.AddError(response.UnmarshalJSON(data))
	}
}
