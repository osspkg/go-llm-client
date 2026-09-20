// Package embeddings implements Ollama /api/embed.
package embeddings

import (
	"context"

	"go.osspkg.com/llm-client/ollama/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls embedding endpoints.
type Client struct{ transport *transport.Client }

// New creates an embeddings client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates embeddings.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/api/embed", "ollama.embed", input, &output)
	return output, err
}

// Legacy creates an embedding through the deprecated /api/embeddings endpoint.
func (client *Client) Legacy(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/api/embeddings", "ollama.embeddings_legacy", input, &output)
	return output, err
}
