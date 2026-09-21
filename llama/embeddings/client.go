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
func New(c *transport.Client) *Client { return &Client{c} }

// Create performs the Create operation.
func (c *Client) Create(ctx context.Context, in Request) (Response, error) {
	var o Response
	err := request.JSON(ctx, c.transport, http.MethodPost, "/embeddings", "embeddings.create", in, &o)
	return o, err
}

// CreateSingle creates an embedding through the singular native alias.
func (c *Client) CreateSingle(ctx context.Context, in Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, c.transport, http.MethodPost, "/embedding", "embeddings.create_single", in, &output)
	return output, err
}
