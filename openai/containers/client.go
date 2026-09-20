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

// Client calls container endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a containers client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a container.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateContainer(ctx, input)
}
