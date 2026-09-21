// Package embeddings provides typed embeddings API operations.
package embeddings

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/llama/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create requests pooled or token-level embeddings.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, http.MethodPost, "/embeddings", "embeddings.create", input, &output)
	return output, err
}

// CreateSingle creates an embedding through the singular native alias.
func (client *Client) CreateSingle(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, http.MethodPost, "/embedding", "embeddings.create_single", input, &output)
	return output, err
}
