/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package generate implements Ollama /api/generate.
package generate

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/ollama/internal/request"
	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls generate endpoints.
type Client struct{ transport *transport.Client }

// New creates a generate client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create generates a response without streaming.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/api/generate", "ollama.generate", input, &output)
	return output, err
}

// CreateStream generates a typed NDJSON stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[Response], error) {
	input.Stream = true
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	responseBody, err := client.transport.Stream(ctx, http.MethodPost, "/api/generate", "ollama.generate_stream", body, "application/json", "application/x-ndjson")
	if err != nil {
		return nil, err
	}
	return stream.NewNDJSON(responseBody, request.Decode[Response], 0), nil
}
