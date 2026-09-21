/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package server provides typed server API operations.
package server

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

// Health reports whether the native server is ready to accept requests.
func (client *Client) Health(ctx context.Context) (HealthResponse, error) {
	var output HealthResponse
	err := request.JSON(ctx, client.transport, http.MethodGet, "/health", "server.health", nil, &output)
	return output, err
}

// Props returns native server properties and default generation settings.
func (client *Client) Props(ctx context.Context) (Props, error) {
	var output Props
	err := request.JSON(ctx, client.transport, http.MethodGet, "/props", "server.props", nil, &output)
	return output, err
}

// UpdateProps changes mutable native server properties.
func (client *Client) UpdateProps(ctx context.Context, input Props) (Props, error) {
	var output Props
	err := request.JSON(ctx, client.transport, http.MethodPost, "/props", "server.props.update", input, &output)
	return output, err
}
