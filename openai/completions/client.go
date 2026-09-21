/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package completions implements legacy OpenAI completions.
package completions

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls completion endpoints.
type Client struct{ transport *transport.Client }

// New creates a completion client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a completion.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/completions", "completions.create", input, &output)
	return output, err
}

// CreateStream creates a typed completion stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[Response], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/completions", "completions.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[Response], 0), nil
}
