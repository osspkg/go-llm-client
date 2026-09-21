/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package templates provides typed templates API operations.
package templates

import "encoding/json"

//go:generate easyjson -all types.go

// Request contains chat messages to format with the configured template.
type Request struct {
	// Messages contains the ordered conversation turns passed to the template.
	Messages []Message `json:"messages"`
	// ChatTemplate selects an explicit template instead of the server default.
	ChatTemplate string `json:"chat_template,omitempty"`
	// Tools contains optional tool definitions used by templates that render tools.
	Tools []Tool `json:"tools,omitempty"`
}

// Message is one role/content pair consumed by the native chat template.
type Message struct {
	// Role identifies the speaker, such as system, user, assistant, or tool.
	Role string `json:"role"`
	// Content contains the text rendered for this conversation turn.
	Content string `json:"content"`
	// Name optionally identifies the participant for templates that use named speakers.
	Name string `json:"name,omitempty"`
}

// Tool describes a function definition that a chat template may render.
type Tool struct {
	// Type identifies the tool kind, normally "function".
	Type string `json:"type"`
	// Function contains the function name and schema rendered by the template.
	Function Function `json:"function,omitempty"`
}

// Function describes a tool function available to the chat template.
type Function struct {
	// Name is the function identifier shown to the model.
	Name string `json:"name"`
	// Description explains when the model should call the function.
	Description string `json:"description,omitempty"`
	// Parameters is the JSON Schema for function arguments.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// Response contains the formatted prompt.
type Response struct {
	// Prompt is the text produced by the server's chat template.
	Prompt string `json:"prompt"`
}
