// Package chat implements the OpenAI Chat Completions bounded context.
package chat

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls chat completion endpoints.
type Client struct{ transport *transport.Client }

// New creates a chat client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a non-streaming chat completion.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/chat/completions", "chat.create", input, &output)
	return output, err
}

// CreateStream creates a typed chat SSE stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[StreamChunk], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/chat/completions", "chat.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamChunk], 0), nil
}
