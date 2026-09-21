/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package tokenization provides typed tokenization API operations.
package tokenization

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/llama/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Tokenize converts text into model token IDs and optional token pieces.
func (client *Client) Tokenize(ctx context.Context, input TokenizeRequest) (TokenizeResponse, error) {
	var output TokenizeResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/tokenize", "tokenize", input, &output)
	return output, err
}

// Detokenize converts model token IDs back into text and optional pieces.
func (client *Client) Detokenize(ctx context.Context, input DetokenizeRequest) (DetokenizeResponse, error) {
	var output DetokenizeResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/detokenize", "detokenize", input, &output)
	return output, err
}
