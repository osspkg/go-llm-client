// Package lora provides typed lora API operations.
package lora

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

// List performs the List operation.
func (c *Client) List(ctx context.Context) (ListResponse, error) {
	var o ListResponse
	err := request.JSON(ctx, c.transport, http.MethodGet, "/lora-adapters", "lora.list", nil, &o)
	return o, err
}

// Set performs the Set operation.
func (c *Client) Set(ctx context.Context, in ListResponse) (ListResponse, error) {
	var o ListResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/lora-adapters", "lora.set", in, &o)
	return o, err
}
