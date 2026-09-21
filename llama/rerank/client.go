// Package rerank provides typed rerank API operations.
package rerank

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
	err := request.JSON(ctx, c.transport, http.MethodPost, "/rerank", "rerank.create", in, &o)
	return o, err
}
