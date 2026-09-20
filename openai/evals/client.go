// Package evals implements the OpenAI Evals bounded context.
package evals

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is an evaluation request payload.
type Request = stateful.Request

// Resource is an evaluation resource.
type Resource = stateful.Resource

// ListResponse is a cursor page of evaluation resources.
type ListResponse = stateful.ListResponse

// DeleteResponse confirms evaluation deletion.
type DeleteResponse = stateful.DeleteResponse

// Client calls evaluation endpoints.
type Client struct{ stateful *stateful.Client }

// New creates an evals client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates an evaluation.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateEval(ctx, input)
}

// List lists evaluations.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	return client.stateful.ListEvals(ctx)
}

// Get returns an evaluation.
func (client *Client) Get(ctx context.Context, id string) (Resource, error) {
	return client.stateful.GetEval(ctx, id)
}

// Update updates an evaluation.
func (client *Client) Update(ctx context.Context, id string, input Request) (Resource, error) {
	return client.stateful.UpdateEval(ctx, id, input)
}

// Delete deletes an evaluation.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	return client.stateful.DeleteEval(ctx, id)
}
