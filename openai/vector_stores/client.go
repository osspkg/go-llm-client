/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package vector_stores implements the OpenAI Vector Stores bounded context.
package vector_stores //nolint:revive // package path follows the provider domain name.

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is a vector store request payload.
type Request = stateful.Request

// Resource is a vector store resource.
type Resource = stateful.Resource

// ListResponse is a vector store list response.
type ListResponse = stateful.ListResponse

// Client calls vector store endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a vector stores client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a vector store.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateVectorStore(ctx, input)
}

// List lists vector stores.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	return client.stateful.ListVectorStores(ctx)
}

// Get returns a vector store.
func (client *Client) Get(ctx context.Context, id string) (Resource, error) {
	return client.stateful.GetVectorStore(ctx, id)
}

// Update updates a vector store.
func (client *Client) Update(ctx context.Context, id string, input Request) (Resource, error) {
	return client.stateful.UpdateVectorStore(ctx, id, input)
}

// Delete deletes a vector store.
func (client *Client) Delete(ctx context.Context, id string) (stateful.DeleteResponse, error) {
	return client.stateful.DeleteVectorStore(ctx, id)
}
