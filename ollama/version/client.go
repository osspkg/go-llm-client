/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package version implements the Ollama version endpoint.
package version

import (
	"context"

	"go.osspkg.com/llm-client/ollama/internal/request"
	"go.osspkg.com/llm-client/ollama/models"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Response describes an Ollama server version.
type Response = models.VersionResponse

// Client calls the version endpoint.
type Client struct{ transport *transport.Client }

// New creates a version client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Get returns the Ollama server version.
func (client *Client) Get(ctx context.Context) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Get, "/api/version", "ollama.version.get", nil, &output)
	return output, err
}
