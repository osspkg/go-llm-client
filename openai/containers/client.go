// Package containers implements the OpenAI Containers bounded context.
package containers

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is a container request payload.
type Request = stateful.Request

// Resource is a container resource.
type Resource = stateful.Resource

// ListResponse is a cursor page of container resources.
type ListResponse = stateful.ListResponse

// DeleteResponse confirms container deletion.
type DeleteResponse = stateful.DeleteResponse

// Client calls container endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a containers client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a container.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateContainer(ctx, input)
}

// List lists containers.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	return client.stateful.ListContainers(ctx)
}

// Get returns a container.
func (client *Client) Get(ctx context.Context, id string) (Resource, error) {
	return client.stateful.GetContainer(ctx, id)
}

// Delete deletes a container.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	return client.stateful.DeleteContainer(ctx, id)
}
