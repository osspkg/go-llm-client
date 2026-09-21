/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package files

//go:generate easyjson -all types.go

// File describes an uploaded file.
type File struct {
	// ID identifies the file.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Bytes is the uploaded file size.
	Bytes int64 `json:"bytes"`
	// CreatedAt is the Unix upload time.
	CreatedAt int64 `json:"created_at"`
	// Filename is the original file name.
	Filename string `json:"filename"`
	// Purpose records the API feature that may use the file.
	Purpose string `json:"purpose"`
	// Status reports provider-side processing state.
	Status string `json:"status,omitempty"`
}

// ListResponse lists files.
type ListResponse struct {
	// Object identifies the provider list kind.
	Object string `json:"object"`
	// Data contains the returned files.
	Data []File `json:"data"`
	// HasMore reports whether another page exists.
	HasMore bool `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	// ID identifies the deleted file.
	ID string `json:"id"`
	// Object identifies the provider object kind.
	Object string `json:"object"`
	// Deleted confirms that deletion succeeded.
	Deleted bool `json:"deleted"`
}
