// Package server provides typed server API operations.
package server

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

// Health performs the Health operation.
func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	var o HealthResponse
	err := request.JSON(ctx, c.transport, http.MethodGet, "/health", "server.health", nil, &o)
	return o, err
}

// Props performs the Props operation.
func (c *Client) Props(ctx context.Context) (Props, error) {
	var o Props
	err := request.JSON(ctx, c.transport, http.MethodGet, "/props", "server.props", nil, &o)
	return o, err
}

// UpdateProps performs the UpdateProps operation.
func (c *Client) UpdateProps(ctx context.Context, in Props) (Props, error) {
	var o Props
	err := request.JSON(ctx, c.transport, http.MethodPost, "/props", "server.props.update", in, &o)
	return o, err
}
