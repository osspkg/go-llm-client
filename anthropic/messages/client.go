/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package messages implements the Anthropic Messages bounded context.
package messages

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/anthropic/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create performs the Create operation.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/messages", "messages.create", input, &output)
	return output, err
}

// CreateStream performs the CreateStream operation.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[StreamEvent], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/messages", "messages.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamEvent], 0), nil
}

// CountTokens performs the CountTokens operation.
func (client *Client) CountTokens(ctx context.Context, input CountTokensRequest) (CountTokensResponse, error) {
	var output CountTokensResponse
	err := request.JSON(ctx, client.transport, request.Post, "/messages/count_tokens", "messages.count_tokens", input, &output)
	return output, err
}
