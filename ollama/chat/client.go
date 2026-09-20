// Package chat implements Ollama /api/chat.
package chat

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/ollama/internal/request"
	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls chat endpoints.
type Client struct{ transport *transport.Client }

// New creates a chat client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a non-streaming chat response.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/api/chat", "ollama.chat", input, &output)
	return output, err
}

// CreateStream creates an NDJSON chat stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[Response], error) {
	input.Stream = true
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	responseBody, err := client.transport.Stream(ctx, http.MethodPost, "/api/chat", "ollama.chat_stream", body, "application/json", "application/x-ndjson")
	if err != nil {
		return nil, err
	}
	return stream.NewNDJSON(responseBody, request.Decode[Response], 0), nil
}
