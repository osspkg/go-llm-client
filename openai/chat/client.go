/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package chat implements the OpenAI Chat Completions bounded context.
package chat

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls chat completion endpoints.
type Client struct{ transport *transport.Client }

// New creates a chat client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a non-streaming chat completion.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/chat/completions", "chat.create", input, &output)
	return output, err
}

// CreateStream creates a typed chat SSE stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[StreamChunk], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/chat/completions", "chat.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamChunk], 0), nil
}

// List returns stored chat completions.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/chat/completions", "chat.list", nil, &output)
	return output, err
}

// Get returns a stored chat completion.
func (client *Client) Get(ctx context.Context, id string) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Get, "/chat/completions/"+url.PathEscape(id), "chat.get", nil, &output)
	return output, err
}

// Update changes metadata associated with a stored chat completion.
func (client *Client) Update(ctx context.Context, id string, input UpdateRequest) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Patch, "/chat/completions/"+url.PathEscape(id), "chat.update", input, &output)
	return output, err
}

// Delete removes a stored chat completion.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, request.Delete, "/chat/completions/"+url.PathEscape(id), "chat.delete", nil, &output)
	return output, err
}

// Messages returns the messages associated with a stored chat completion.
func (client *Client) Messages(ctx context.Context, id string) (MessagesResponse, error) {
	var output MessagesResponse
	err := request.JSON(ctx, client.transport, request.Get, "/chat/completions/"+url.PathEscape(id)+"/messages", "chat.messages.list", nil, &output)
	return output, err
}
