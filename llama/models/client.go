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
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List returns models currently known by the native router.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, http.MethodGet, "/models", "models.list", nil, &output)
	return output, err
}

// Load loads a model and returns the router operation result.
func (client *Client) Load(ctx context.Context, input ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/models/load", "models.load", input, &output)
	return output, err
}

// Unload unloads a model and returns the router operation result.
func (client *Client) Unload(ctx context.Context, input ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/models/unload", "models.unload", input, &output)
	return output, err
}

// Download starts a background model download through the router.
func (client *Client) Download(ctx context.Context, input ModelRequest) (OperationResponse, error) {
	var output OperationResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/models", "models.download", input, &output)
	return output, err
}

// Events opens the bounded SSE stream of router model lifecycle events.
func (client *Client) Events(ctx context.Context) (stream.Iterator[Event], error) {
	body, err := client.transport.Stream(ctx, http.MethodGet, "/models/sse", "models.events", nil, "", "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[Event], 0), nil
}
