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

// Client calls evaluation endpoints.
type Client struct{ stateful *stateful.Client }

// New creates an evals client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates an evaluation.
func (client *Client) Create(ctx context.Context, input Request) (Resource, error) {
	return client.stateful.CreateEval(ctx, input)
}
