/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package generate

import "encoding/json"

//go:generate easyjson -all types.go

// Request generates a completion.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Prompt contains the text the model should complete.
	Prompt string `json:"prompt,omitempty"`
	// Suffix contains the provider value used for the `Suffix` field.
	Suffix string `json:"suffix,omitempty"`
	// Stream requests incremental events instead of one buffered response.
	Stream bool `json:"stream,omitempty"`
	// Format selects the model or response format.
	Format json.RawMessage `json:"format,omitempty"`
	// Options contains model-specific generation options.
	Options json.RawMessage `json:"options,omitempty"`
	// KeepAlive sets how long Ollama keeps the model loaded.
	KeepAlive string `json:"keep_alive,omitempty"`
}

// Response is one generate response or a complete non-stream response.
type Response struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt string `json:"created_at"`
	// Response contains generated text or a generated response object.
	Response string `json:"response,omitempty"`
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
	// Thinking contains the model reasoning text when returned.
	Thinking string `json:"thinking,omitempty"`
}
