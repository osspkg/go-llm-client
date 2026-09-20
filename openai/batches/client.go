// Package batches implements OpenAI batch operations.
package batches

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls batch endpoints.
type Client struct{ transport *transport.Client }

// New creates a batch client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a batch.
func (client *Client) Create(ctx context.Context, input Request) (Batch, error) {
	var output Batch
	err := request.JSON(ctx, client.transport, request.Post, "/batches", "batches.create", input, &output)
	return output, err
}

// Get returns a batch.
func (client *Client) Get(ctx context.Context, id string) (Batch, error) {
	var output Batch
	err := request.JSON(ctx, client.transport, request.Get, "/batches/"+url.PathEscape(id), "batches.get", nil, &output)
	return output, err
}

// Cancel cancels a batch.
func (client *Client) Cancel(ctx context.Context, id string) (Batch, error) {
	var output Batch
	err := request.JSON(ctx, client.transport, http.MethodPost, "/batches/"+url.PathEscape(id)+"/cancel", "batches.cancel", nil, &output)
	return output, err
}
