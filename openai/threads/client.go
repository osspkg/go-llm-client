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
