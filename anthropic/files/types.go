/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package files provides typed Anthropic Files API operations.
package files

//go:generate easyjson -all types.go

// File describes an uploaded provider file.
type File struct {
	// ID uniquely identifies the uploaded file.
	ID string `json:"id"`
	// Type identifies the provider resource type.
	Type string `json:"type"`
	// Name is the original filename supplied during upload.
	Name string `json:"filename,omitempty"`
	// MimeType is the detected or declared media type.
	MimeType string `json:"mime_type,omitempty"`
	// SizeBytes is the file size reported by the provider.
	SizeBytes int64 `json:"size_bytes,omitempty"`
	// CreatedAt is the provider creation timestamp in ISO-8601 form.
	CreatedAt string `json:"created_at,omitempty"`
}

// ListResponse is a cursor page of uploaded files.
type ListResponse struct {
	// Data contains files returned for the current page.
	Data []File `json:"data"`
	// HasMore reports whether another page can be requested.
	HasMore bool `json:"has_more"`
	// FirstID is the first file ID in the page.
	FirstID string `json:"first_id,omitempty"`
	// LastID is the last file ID in the page.
	LastID string `json:"last_id,omitempty"`
}

// DeleteResponse confirms file deletion.
type DeleteResponse struct {
	// ID identifies the deleted file.
	ID string `json:"id"`
	// Type identifies the provider resource type.
	Type string `json:"type"`
	// Deleted reports whether deletion succeeded.
	Deleted bool `json:"deleted"`
}
