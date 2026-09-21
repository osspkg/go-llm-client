/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package uploads

import "io"

//go:generate easyjson -all types.go

// CreateRequest starts a large file upload.
type CreateRequest struct {
	// Filename contains the original or uploaded file name.
	Filename string `json:"filename"`
	// Purpose identifies the API feature allowed to use the file.
	Purpose string `json:"purpose"`
	// Bytes contains the file or upload size in bytes.
	Bytes int64 `json:"bytes"`
	// MimeType contains the uploaded file MIME type.
	MimeType string `json:"mime_type"`
}

// Upload describes an in-progress or completed upload.
type Upload struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Bytes contains the file or upload size in bytes.
	Bytes int64 `json:"bytes"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// Filename contains the original or uploaded file name.
	Filename string `json:"filename"`
	// Purpose identifies the API feature allowed to use the file.
	Purpose string `json:"purpose"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status"`
}

// Part describes a completed upload part.
type Part struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// UploadID identifies the parent upload.
	UploadID string `json:"upload_id"`
	// Bytes contains the file or upload size in bytes.
	Bytes int64 `json:"bytes"`
}

// CompleteRequest completes an upload with ordered part IDs.
type CompleteRequest struct {
	// PartIDs lists upload parts in the order they should be assembled.
	PartIDs []string `json:"part_ids"`
}

// PartInput contains one bounded upload part reader.
//
//easyjson:skip
type PartInput struct {
	Filename string
	File     io.Reader
}

// DeleteResponse confirms an upload cancellation.
type DeleteResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
}
