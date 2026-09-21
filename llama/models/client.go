// Package models provides typed models API operations.
package models

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/llama/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(c *transport.Client) *Client { return &Client{c} }

// List performs the List operation.
func (c *Client) List(ctx context.Context) (ListResponse, error) {
	var o ListResponse
	err := request.JSON(ctx, c.transport, http.MethodGet, "/models", "models.list", nil, &o)
	return o, err
}

// Load loads a model and returns the router operation result.
func (c *Client) Load(ctx context.Context, in ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/models/load", "models.load", in, &output)
	return output, err
}

// Unload unloads a model and returns the router operation result.
func (c *Client) Unload(ctx context.Context, in ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/models/unload", "models.unload", in, &output)
	return output, err
}

// Download starts a background model download through the router.
func (c *Client) Download(ctx context.Context, in ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/models", "models.download", in, &output)
	return output, err
}

// Events opens the bounded SSE stream of router model lifecycle events.
func (c *Client) Events(ctx context.Context) (stream.Iterator[Event], error) {
	body, err := c.transport.Stream(ctx, http.MethodGet, "/models/sse", "models.events", nil, "", "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[Event], 0), nil
}
