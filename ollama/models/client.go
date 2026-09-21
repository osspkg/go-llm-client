/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package models implements Ollama model lifecycle operations.
package models

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/ollama/internal/request"
	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls model lifecycle endpoints.
type Client struct{ transport *transport.Client }

// New creates a model client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List returns local models.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/api/tags", "ollama.models.list", nil, &output)
	return output, err
}

// Running returns models currently loaded in memory.
func (client *Client) Running(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/api/ps", "ollama.models.running", nil, &output)
	return output, err
}

// Show returns model information.
func (client *Client) Show(ctx context.Context, model string) (ShowResponse, error) {
	var output ShowResponse
	err := request.JSON(ctx, client.transport, request.Post, "/api/show", "ollama.models.show", Request{Model: model}, &output)
	return output, err
}

// Create creates a model.
func (client *Client) Create(ctx context.Context, input Request) (stream.Iterator[Progress], error) {
	return client.lifecycleStream(ctx, "/api/create", "ollama.models.create", input)
}

// Pull downloads a model.
func (client *Client) Pull(ctx context.Context, input Request) (stream.Iterator[Progress], error) {
	return client.lifecycleStream(ctx, "/api/pull", "ollama.models.pull", input)
}

// Push uploads a model.
func (client *Client) Push(ctx context.Context, input Request) (stream.Iterator[Progress], error) {
	return client.lifecycleStream(ctx, "/api/push", "ollama.models.push", input)
}

// Copy copies a model.
func (client *Client) Copy(ctx context.Context, source, destination string) error {
	input := Request{Model: destination, From: source}
	body, err := codec.Marshal(input)
	if err != nil {
		return err
	}
	_, err = client.transport.Request(ctx, http.MethodPost, "/api/copy", "ollama.models.copy", body, "application/json")
	return err
}

// Delete deletes a model.
func (client *Client) Delete(ctx context.Context, model string) error {
	input := Request{Model: model}
	body, err := codec.Marshal(input)
	if err != nil {
		return err
	}
	_, err = client.transport.Request(ctx, http.MethodDelete, "/api/delete", "ollama.models.delete", body, "application/json")
	return err
}

// Version returns the server version.
func (client *Client) Version(ctx context.Context) (VersionResponse, error) {
	var output VersionResponse
	err := request.JSON(ctx, client.transport, request.Get, "/api/version", "ollama.version", nil, &output)
	return output, err
}

func (client *Client) lifecycleStream(ctx context.Context, endpoint, operation string, input Request) (stream.Iterator[Progress], error) {
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	responseBody, err := client.transport.Stream(ctx, http.MethodPost, endpoint, operation, body, "application/json", "application/x-ndjson")
	if err != nil {
		return nil, err
	}
	return stream.NewNDJSON(responseBody, request.Decode[Progress], 0), nil
}
