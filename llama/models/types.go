/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package models provides typed native router model operations.
package models

import "encoding/json"

//go:generate easyjson -all types.go

// Model describes a model known to the llama.cpp router.
type Model struct {
	// ID is the router model name used in load and inference requests.
	ID string `json:"id"`
	// Status contains the router's lifecycle state and process details.
	Status ModelStatus `json:"status,omitempty"`
	// Path is the local model file path when the router exposes it.
	Path string `json:"path,omitempty"`
	// Modalities lists supported input or output modalities for compatible servers.
	Modalities []string `json:"modalities,omitempty"`
	// Architecture describes input and output modalities advertised by the router.
	Architecture Architecture `json:"architecture,omitempty"`
}

// ModelStatus describes the router lifecycle state of a model.
type ModelStatus struct {
	// Value is unloaded, loading, loaded, sleeping, or downloading.
	Value string `json:"value,omitempty"`
	// Args contains the server arguments for a running model instance.
	Args []string `json:"args,omitempty"`
	// Failed reports that loading or downloading failed.
	Failed bool `json:"failed,omitempty"`
	// ExitCode is the process exit code after a failed load.
	ExitCode int `json:"exit_code,omitempty"`
	// Progress contains provider-defined download progress details.
	Progress json.RawMessage `json:"progress,omitempty"`
}

// Architecture describes model input and output modalities.
type Architecture struct {
	// InputModalities lists modalities accepted by the model.
	InputModalities []string `json:"input_modalities,omitempty"`
	// OutputModalities lists modalities produced by the model.
	OutputModalities []string `json:"output_modalities,omitempty"`
}

// ListResponse contains models returned by a native router.
type ListResponse struct {
	// Data contains the models available in the router cache.
	Data []Model `json:"data,omitempty"`
	// Models accepts compatible servers that use a models property.
	Models []Model `json:"models,omitempty"`
}

// ModelRequest selects a router model for loading, unloading, or downloading.
type ModelRequest struct {
	// Model is the model alias, repository reference, or preset name.
	Model string `json:"model"`
}

// OperationResponse reports whether a router model operation was accepted.
type OperationResponse struct {
	// Success reports whether the server accepted the operation.
	Success bool `json:"success"`
}

// Event is one router model lifecycle event.
type Event struct {
	// Model identifies the model affected by the event, or * for a full reload.
	Model string `json:"model"`
	// Name identifies model_status, download_progress, model_remove, or models_reload.
	Name string `json:"event"`
	// Data contains provider-defined event details such as status or progress.
	Data json.RawMessage `json:"data,omitempty"`
}
