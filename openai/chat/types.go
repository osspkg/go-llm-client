/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package chat

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a chat completion.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Messages contains the conversation history sent to the model.
	Messages []Message `json:"messages"`
	// Stream requests incremental events instead of one buffered response.
	Stream bool `json:"stream,omitempty"`
	// MaxTokens contains the provider value used for the `MaxTokens` field.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Temperature contains the provider value used for the `Temperature` field.
	Temperature float64 `json:"temperature,omitempty"`
	// TopP contains the provider value used for the `TopP` field.
	TopP float64 `json:"top_p,omitempty"`
	// Tools declares tools the model may call.
	Tools []Tool `json:"tools,omitempty"`
	// ToolChoice contains the provider value used for the `ToolChoice` field.
	ToolChoice json.RawMessage `json:"tool_choice,omitempty"`
	// ResponseFormat selects the image response representation, such as URL or base64.
	ResponseFormat json.RawMessage `json:"response_format,omitempty"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateRequest changes metadata associated with a stored chat completion.
type UpdateRequest struct {
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata"`
}

// Message is a chat message.
type Message struct {
	// Role identifies the organization or message role.
	Role string `json:"role"`
	// Content contains the message or content-part text.
	Content json.RawMessage `json:"content,omitempty"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// ToolCalls contains tool calls produced by the model.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	// ToolCallID contains the provider value used for the `ToolCallID` field.
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// Tool describes a function tool.
type Tool struct {
	// Type identifies the object or event variant.
	Type string `json:"type"`
	// Function contains the function declaration or invocation.
	Function FunctionTool `json:"function"`
}

// FunctionTool describes a callable function.
type FunctionTool struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// Description explains a tool or resource so the model can use it correctly.
	Description string `json:"description,omitempty"`
	// Parameters contains the JSON Schema or arguments accepted by a function.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall is a model-generated tool call.
type ToolCall struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Type identifies the object or event variant.
	Type string `json:"type"`
	// Index identifies the zero-based position of an item or choice.
	Index int `json:"index,omitempty"`
	// Function contains the function declaration or invocation.
	Function FunctionCall `json:"function"`
}

// FunctionCall contains generated function arguments.
type FunctionCall struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Arguments contains the JSON arguments supplied to a function.
	Arguments string `json:"arguments,omitempty"`
}

// Response is a chat completion.
type Response struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Created contains the Unix timestamp when the completion was created.
	Created int64 `json:"created"`
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Choices contains generated alternatives or streamed choices.
	Choices []Choice `json:"choices"`
	// Usage contains token accounting for the operation.
	Usage *Usage `json:"usage,omitempty"`
	// ServiceTier contains the provider value used for the `ServiceTier` field.
	ServiceTier string `json:"service_tier,omitempty"`
}

// Choice is a completion choice.
type Choice struct {
	// Index identifies the zero-based position of an item or choice.
	Index int `json:"index"`
	// Message contains the provider error message.
	Message Message `json:"message"`
	// Delta contains incremental text or argument content.
	Delta Message `json:"delta,omitempty"`
	// FinishReason explains why model generation stopped.
	FinishReason string `json:"finish_reason,omitempty"`
}

// Usage contains token accounting.
type Usage struct {
	// PromptTokens counts tokens in the input prompt.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens counts tokens generated for the completion.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens contains total input and output tokens.
	TotalTokens int `json:"total_tokens"`
}

// StreamChunk is one streamed chat event.
type StreamChunk struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Created contains the Unix timestamp when the completion was created.
	Created int64 `json:"created"`
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Choices contains generated alternatives or streamed choices.
	Choices []Choice `json:"choices"`
	// Usage contains token accounting for the operation.
	Usage *Usage `json:"usage,omitempty"`
}

// ListResponse lists stored chat completions.
type ListResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Response `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// DeleteResponse confirms deletion of a stored chat completion.
type DeleteResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
}

// MessagesResponse lists the messages associated with a stored completion.
type MessagesResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Message `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}
