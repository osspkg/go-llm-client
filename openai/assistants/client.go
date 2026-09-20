// Package assistants implements the OpenAI Assistants bounded context.
package assistants

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is an assistant request payload.
type Request = stateful.Request

// Resource is an assistant resource.
type Resource = stateful.Resource

// ListResponse is an assistant list response.
type ListResponse = stateful.ListResponse

// DeleteResponse confirms assistant deletion.
type DeleteResponse = stateful.DeleteResponse

// Client calls assistant endpoints.
type Client struct{ stateful *stateful.Client }

// New creates an assistants client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates an assistant.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateAssistant(ctx, input)
}

// Get returns an assistant.
func (client *Client) Get(ctx context.Context, id string) (Resource, error) {
	return client.stateful.GetAssistant(ctx, id)
}

// List lists assistants.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	return client.stateful.ListAssistants(ctx)
}

// Delete deletes an assistant.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	return client.stateful.DeleteAssistant(ctx, id)
}
