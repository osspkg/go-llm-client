/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package rerank provides typed rerank API operations.
package rerank

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

// Create ranks the submitted documents against a query.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, http.MethodPost, "/rerank", "rerank.create", input, &output)
	return output, err
}
