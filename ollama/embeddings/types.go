/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package embeddings

import "encoding/json"

//go:generate easyjson -all types.go

// Request creates one or more embeddings.
type Request struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Input contains the text, tokens, or input items to process.
	Input json.RawMessage `json:"input"`
	// Truncate controls whether inputs exceeding the context window are truncated.
	Truncate *bool `json:"truncate,omitempty"`
	// Dimensions requests the number of dimensions for generated embeddings.
	Dimensions int `json:"dimensions,omitempty"`
	// KeepAlive sets how long Ollama keeps the model loaded.
	KeepAlive string `json:"keep_alive,omitempty"`
	// Options contains model-specific generation options.
	Options json.RawMessage `json:"options,omitempty"`
}

// Response contains embeddings.
type Response struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Embeddings contains one numeric vector per input item.
	Embeddings [][]float64 `json:"embeddings"`
	// TotalDuration contains total Ollama processing time in nanoseconds.
	TotalDuration int64 `json:"total_duration,omitempty"`
	// LoadDuration contains Ollama model-load time in nanoseconds.
	LoadDuration int64 `json:"load_duration,omitempty"`
	// PromptEvalCount counts input tokens evaluated by Ollama.
	PromptEvalCount int `json:"prompt_eval_count,omitempty"`
}
