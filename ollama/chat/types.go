/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package chat

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates a chat response.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Messages contains the conversation history sent to the model.
	Messages []Message `json:"messages"`
	// Stream requests incremental events instead of one buffered response.
	Stream bool `json:"stream,omitempty"`
	// Tools declares tools the model may call.
	Tools []Tool `json:"tools,omitempty"`
	// Think controls reasoning output for models that support thinking.
	Think json.RawMessage `json:"think,omitempty"`
	// Format selects the model or response format.
	Format json.RawMessage `json:"format,omitempty"`
	// Options contains model-specific generation options.
	Options json.RawMessage `json:"options,omitempty"`
	// KeepAlive sets how long Ollama keeps the model loaded.
	KeepAlive string `json:"keep_alive,omitempty"`
}

// Message is an Ollama chat message.
type Message struct {
	// Role identifies the organization or message role.
	Role string `json:"role"`
	// Content contains the message or content-part text.
	Content string `json:"content,omitempty"`
	// Thinking contains the model reasoning text when returned.
	Thinking string `json:"thinking,omitempty"`
	// Images contains base64-encoded images attached to a multimodal message.
	Images []string `json:"images,omitempty"`
	// ToolCalls contains tool calls produced by the model.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Tool is an Ollama tool declaration.
type Tool struct {
	// Type identifies the object or event variant.
	Type string `json:"type"`
	// Function contains the function declaration or invocation.
	Function FunctionTool `json:"function"`
}

// FunctionTool describes a tool.
type FunctionTool struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// Description explains a tool or resource so the model can use it correctly.
	Description string `json:"description,omitempty"`
	// Parameters contains the JSON Schema or arguments accepted by a function.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall is a model tool call.
type ToolCall struct {
	// Function contains the function declaration or invocation.
	Function FunctionCall `json:"function"`
}

// FunctionCall contains tool arguments.
type FunctionCall struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// Arguments contains the JSON arguments supplied to a function.
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// Response is one chat response event.
type Response struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt string `json:"created_at"`
	// Message contains the provider error message.
	Message Message `json:"message"`
	// Done reports whether Ollama finished the streamed operation.
	Done bool `json:"done"`
	// TotalDuration contains total Ollama processing time in nanoseconds.
	TotalDuration int64 `json:"total_duration,omitempty"`
	// LoadDuration contains Ollama model-load time in nanoseconds.
	LoadDuration int64 `json:"load_duration,omitempty"`
	// PromptEvalCount counts input tokens evaluated by Ollama.
	PromptEvalCount int `json:"prompt_eval_count,omitempty"`
	// EvalCount counts output tokens evaluated by Ollama.
	EvalCount int `json:"eval_count,omitempty"`
	// DoneReason explains why generation stopped.
	DoneReason string `json:"done_reason,omitempty"`
}
