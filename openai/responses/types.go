/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package responses

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a response.
type Request struct {
	// Model selects the response model.
	Model string `json:"model"`
	// Input contains upstream-defined input items.
	Input json.RawMessage `json:"input,omitempty"`
	// Instructions provide system-level guidance.
	Instructions string `json:"instructions,omitempty"`
	// Stream requests incremental SSE events.
	Stream bool `json:"stream,omitempty"`
	// MaxOutputTokens limits generated tokens.
	MaxOutputTokens int `json:"max_output_tokens,omitempty"`
	// Metadata attaches application-defined labels.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Tools declares tools available to the model.
	Tools []Tool `json:"tools,omitempty"`
}

// Tool describes a Responses tool declaration.
type Tool struct {
	// Type selects the tool kind.
	Type string `json:"type"`
	// Name identifies a function tool.
	Name string `json:"name,omitempty"`
	// Description tells the model when to invoke the tool.
	Description string `json:"description,omitempty"`
	// Parameters is the function's JSON Schema.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// Response is the typed response object.
type Response struct {
	// ID identifies the response.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// CreatedAt is the Unix creation time.
	CreatedAt int64 `json:"created_at"`
	// Model identifies the model that generated the response.
	Model string `json:"model"`
	// Status reports response lifecycle state.
	Status string `json:"status"`
	// Output contains structured output items.
	Output []OutputItem `json:"output,omitempty"`
	// OutputText is provider-assembled text output.
	OutputText string `json:"output_text,omitempty"`
	// Usage reports token accounting.
	Usage *Usage `json:"usage,omitempty"`
	// Error contains a provider failure, if any.
	Error *Error `json:"error,omitempty"`
}

// OutputItem is a response output item.
type OutputItem struct {
	// ID identifies this output item.
	ID string `json:"id"`
	// Type identifies the output item variant.
	Type string `json:"type"`
	// Role identifies the producing conversation role.
	Role string `json:"role,omitempty"`
	// Status reports item lifecycle state.
	Status string `json:"status,omitempty"`
	// Content contains typed content parts.
	Content []OutputContent `json:"content,omitempty"`
	// Text contains direct item text when present.
	Text string `json:"text,omitempty"`
}

// OutputContent is a typed output content item.
type OutputContent struct {
	// Type identifies the content variant.
	Type string `json:"type"`
	// Text contains generated text.
	Text string `json:"text,omitempty"`
	// Refusal contains a safety refusal.
	Refusal string `json:"refusal,omitempty"`
	// Annotations contains citations or related metadata.
	Annotations []Annotation `json:"annotations,omitempty"`
}

// Annotation is an output annotation.
type Annotation struct {
	// Type identifies the annotation variant.
	Type string `json:"type"`
	// URL is the referenced resource location.
	URL string `json:"url,omitempty"`
}

// Usage contains token accounting.
type Usage struct {
	// InputTokens counts input tokens.
	InputTokens int `json:"input_tokens"`
	// OutputTokens counts generated tokens.
	OutputTokens int `json:"output_tokens"`
	// TotalTokens is the total token count.
	TotalTokens int `json:"total_tokens"`
}

// Error is a provider response error.
type Error struct {
	// Code is the provider error code.
	Code string `json:"code,omitempty"`
	// Message is a safe provider error description.
	Message string `json:"message,omitempty"`
}

// StreamEvent is one Responses SSE event.
type StreamEvent struct {
	// Type identifies the SSE event variant.
	Type string `json:"type"`
	// SequenceNumber orders events from the provider.
	SequenceNumber int `json:"sequence_number,omitempty"`
	// ResponseID identifies the response being streamed.
	ResponseID string `json:"response_id,omitempty"`
	// Delta contains incremental text.
	Delta string `json:"delta,omitempty"`
	// Text contains complete text supplied by an event.
	Text string `json:"text,omitempty"`
	// Response contains a complete response when supplied.
	Response *Response `json:"response,omitempty"`
}

// DeleteResponse confirms deletion of a stored response.
type DeleteResponse struct {
	// ID identifies the deleted response.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Deleted confirms deletion.
	Deleted bool `json:"deleted"`
}

// InputItem is an input item returned by the Responses API.
type InputItem struct {
	// ID identifies the input item when assigned.
	ID string `json:"id,omitempty"`
	// Type identifies the input item variant.
	Type string `json:"type"`
	// Role identifies the conversation role.
	Role string `json:"role,omitempty"`
	// Status reports item lifecycle state.
	Status string `json:"status,omitempty"`
	// Content preserves the spec-defined arbitrary input payload.
	Content json.RawMessage `json:"content,omitempty"`
}

// InputItemsResponse is a cursor page of response input items.
type InputItemsResponse struct {
	// Object identifies the provider list kind.
	Object string `json:"object"`
	// Data contains input items.
	Data []InputItem `json:"data"`
	// FirstID is the first cursor item.
	FirstID string `json:"first_id,omitempty"`
	// LastID is the last cursor item.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page exists.
	HasMore bool `json:"has_more"`
}
