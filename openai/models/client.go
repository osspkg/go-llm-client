/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package models implements OpenAI model operations.
package models

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls model endpoints.
type Client struct{ transport *transport.Client }

// New creates a model client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List returns available models.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/models", "models.list", nil, &output)
	return output, err
}

// Get returns one model.
func (client *Client) Get(ctx context.Context, id string) (Model, error) {
	var output Model
	endpoint := "/models/" + url.PathEscape(id)
	err := request.JSON(ctx, client.transport, request.Get, endpoint, "models.get", nil, &output)
	return output, err
}

// Delete deletes a fine-tuned model.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	endpoint := path.Join("/models", url.PathEscape(id))
	err := request.JSON(ctx, client.transport, http.MethodDelete, endpoint, "models.delete", nil, &output)
	return output, err
}
