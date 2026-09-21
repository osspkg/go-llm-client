/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package moderations implements OpenAI moderation operations.
package moderations

import (
	"context"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls moderation endpoints.
type Client struct{ transport *transport.Client }

// New creates a moderation client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create classifies input.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/moderations", "moderations.create", input, &output)
	return output, err
}
