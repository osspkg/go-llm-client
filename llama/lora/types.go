/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package lora provides typed LoRA adapter API operations.
package lora

import (
	"encoding/json"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

//go:generate easyjson -all types.go

// Adapter identifies a loaded LoRA adapter and its global scale.
type Adapter struct {
	// ID is the server-assigned adapter identifier used in completion requests.
	ID int `json:"id"`
	// Path is the filesystem path of the loaded adapter.
	Path string `json:"path"`
	// Scale is the global multiplier applied to this adapter; zero disables it.
	Scale float64 `json:"scale"`
}

// ListResponse is the native endpoint's array of LoRA adapters.
type ListResponse []Adapter

// MarshalJSON encodes the native adapter array.
func (response ListResponse) MarshalJSON() ([]byte, error) {
	type responseAlias []Adapter
	return json.Marshal(responseAlias(response))
}

// UnmarshalJSON decodes the native adapter array.
func (response *ListResponse) UnmarshalJSON(data []byte) error {
	type responseAlias []Adapter
	var value responseAlias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*response = ListResponse(value)
	return nil
}

// MarshalEasyJSON implements easyjson.Marshaler for the adapter array.
func (response ListResponse) MarshalEasyJSON(writer *jwriter.Writer) {
	data, err := response.MarshalJSON()
	writer.Raw(data, err)
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler for the adapter array.
func (response *ListResponse) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	data := lexer.Raw()
	if lexer.Ok() {
		lexer.AddError(response.UnmarshalJSON(data))
	}
}
