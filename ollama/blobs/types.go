/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package blobs

//go:generate easyjson -all types.go

// Upload identifies an Ollama blob upload.
type Upload struct {
	// Digest identifies an immutable blob or model content by SHA-256 digest.
	Digest string `json:"digest"`
}
