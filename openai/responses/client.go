// Package responses implements the OpenAI Responses bounded context.
package responses

import (
	"context"
	"net/http"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls Responses endpoints.
type Client struct{ transport *transport.Client }

// New creates a Responses client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a non-streaming response.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/responses", "responses.create", input, &output)
	return output, err
}

// CreateStream creates a typed SSE response stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[StreamEvent], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/responses", "responses.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamEvent], 0), nil
}
