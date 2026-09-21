// Package tokenization provides typed tokenization API operations.
package tokenization

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

// Tokenize performs the Tokenize operation.
func (c *Client) Tokenize(ctx context.Context, in TokenizeRequest) (TokenizeResponse, error) {
	var o TokenizeResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/tokenize", "tokenize", in, &o)
	return o, err
}

// Detokenize performs the Detokenize operation.
func (c *Client) Detokenize(ctx context.Context, in DetokenizeRequest) (DetokenizeResponse, error) {
	var o DetokenizeResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/detokenize", "detokenize", in, &o)
	return o, err
}
