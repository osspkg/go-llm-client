/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package threads implements the OpenAI Threads bounded context.
package threads

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is a thread request payload.
type Request = stateful.Request

// Resource is a thread resource.
type Resource = stateful.Resource

// DeleteResponse confirms thread deletion.
type DeleteResponse = stateful.DeleteResponse

// ListResponse is a cursor page of thread resources.
type ListResponse = stateful.ListResponse

// Client calls thread endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a threads client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a thread.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateThread(ctx, input)
}

// Get returns a thread.
func (client *Client) Get(ctx context.Context, id string) (Resource, error) {
	return client.stateful.GetThread(ctx, id)
}

// Delete deletes a thread.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	return client.stateful.DeleteThread(ctx, id)
}

// Update updates a thread.
func (client *Client) Update(ctx context.Context, id string, input Request) (Resource, error) {
	return client.stateful.UpdateThread(ctx, id, input)
}

// CreateAndRun creates a thread and starts its initial run.
func (client *Client) CreateAndRun(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateThreadAndRun(ctx, input)
}

// ListMessages lists messages in a thread.
func (client *Client) ListMessages(ctx context.Context, id string) (ListResponse, error) {
	return client.stateful.ListMessages(ctx, id)
}

// CreateMessage adds a message to a thread.
func (client *Client) CreateMessage(ctx context.Context, id string, input Request) (Resource, error) {
	return client.stateful.CreateMessage(ctx, id, input)
}

// GetMessage returns a message in a thread.
func (client *Client) GetMessage(ctx context.Context, threadID, messageID string) (Resource, error) {
	return client.stateful.GetMessage(ctx, threadID, messageID)
}

// UpdateMessage updates a message in a thread.
func (client *Client) UpdateMessage(ctx context.Context, threadID, messageID string, input Request) (Resource, error) {
	return client.stateful.UpdateMessage(ctx, threadID, messageID, input)
}

// DeleteMessage deletes a message in a thread.
func (client *Client) DeleteMessage(ctx context.Context, threadID, messageID string) (DeleteResponse, error) {
	return client.stateful.DeleteMessage(ctx, threadID, messageID)
}
