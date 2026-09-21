/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package tokenization provides typed native tokenization operations.
package tokenization

import (
	"encoding/json"
	"fmt"
)

//go:generate easyjson -all types.go

// TokenizeRequest controls conversion of text into model tokens.
type TokenizeRequest struct {
	// Content is the text to tokenize.
	Content string `json:"content"`
	// AddSpecial requests insertion of special tokens such as BOS.
	AddSpecial bool `json:"add_special,omitempty"`
	// ParseSpecial controls whether special-token spellings are interpreted as tokens.
	ParseSpecial bool `json:"parse_special,omitempty"`
	// WithPieces requests token IDs together with their text/byte pieces.
	WithPieces bool `json:"with_pieces,omitempty"`
}

// TokenPiece is one token ID optionally paired with its text or raw byte representation.
//
//easyjson:skip
type TokenPiece struct {
	// ID is the model vocabulary token ID.
	ID int `json:"id"`
	// Piece is a UTF-8 string or a byte array when the token is not valid UTF-8.
	Piece    PieceValue `json:"piece,omitempty"`
	hasPiece bool
}

// UnmarshalJSON accepts either a bare token ID or a detailed token object.
func (piece *TokenPiece) UnmarshalJSON(data []byte) error {
	var id int
	if err := json.Unmarshal(data, &id); err == nil {
		piece.ID = id
		piece.Piece = PieceValue{}
		piece.hasPiece = false
		return nil
	}
	type tokenPieceAlias TokenPiece
	var value tokenPieceAlias
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("decode token: %w", err)
	}
	*piece = TokenPiece(value)
	piece.hasPiece = true
	return nil
}

// MarshalJSON emits a bare token ID when no piece details are present.
func (piece TokenPiece) MarshalJSON() ([]byte, error) {
	if !piece.hasPiece && piece.Piece.Text == "" && piece.Piece.Bytes == nil {
		return json.Marshal(piece.ID)
	}
	type tokenPieceAlias TokenPiece
	return json.Marshal(tokenPieceAlias(piece))
}

// PieceValue is the typed string-or-bytes token piece union.
//
//easyjson:skip
type PieceValue struct {
	// Text contains the UTF-8 representation when the server returned a string.
	Text string
	// Bytes contains raw bytes when the server returned a byte array.
	Bytes []byte
}

// UnmarshalJSON decodes either a JSON string or an array of byte values.
func (piece *PieceValue) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		piece.Text = text
		piece.Bytes = nil
		return nil
	}
	var raw []int
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode token piece: %w", err)
	}
	piece.Text = ""
	piece.Bytes = make([]byte, len(raw))
	for index, value := range raw {
		if value < 0 || value > 255 {
			return fmt.Errorf("decode token piece: byte %d out of range", value)
		}
		piece.Bytes[index] = byte(value)
	}
	return nil
}

// MarshalJSON encodes the text or raw byte representation.
func (piece PieceValue) MarshalJSON() ([]byte, error) {
	if piece.Text != "" || piece.Bytes == nil {
		return json.Marshal(piece.Text)
	}
	values := make([]int, len(piece.Bytes))
	for index, value := range piece.Bytes {
		values[index] = int(value)
	}
	return json.Marshal(values)
}

// TokenizeResponse contains token IDs or detailed token pieces.
type TokenizeResponse struct {
	// Tokens contains IDs when with_pieces is false, or detailed pieces otherwise.
	Tokens []TokenPiece `json:"tokens,omitempty"`
}

// DetokenizeRequest contains token IDs to convert back into text.
type DetokenizeRequest struct {
	// Tokens contains vocabulary token IDs in model order.
	Tokens []int `json:"tokens"`
}

// DetokenizeResponse contains the reconstructed text.
type DetokenizeResponse struct {
	// Content is the text reconstructed from the supplied tokens.
	Content string `json:"content"`
}
