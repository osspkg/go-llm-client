// Package runs implements the OpenAI Runs bounded context.
package runs

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is a run request payload.
type Request = stateful.Request

// Resource is a run resource.
type Resource = stateful.Resource

// Client calls run endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a runs client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a run in a thread.
func (client *Client) Create(ctx context.Context, threadID string, input Request) (Resource, error) {
	return client.stateful.CreateRun(ctx, threadID, input)
}

// Get returns a run.
func (client *Client) Get(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.stateful.GetRun(ctx, threadID, runID)
}

// Cancel cancels a run.
func (client *Client) Cancel(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.stateful.CancelRun(ctx, threadID, runID)
}
