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
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List returns all configured LoRA adapters.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, http.MethodGet, "/lora-adapters", "lora.list", nil, &output)
	return output, err
}

// Set replaces the server's active LoRA adapter configuration.
func (client *Client) Set(ctx context.Context, input ListResponse) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/lora-adapters", "lora.set", input, &output)
	return output, err
}
