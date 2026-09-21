// Package completions provides typed completions API operations.
package completions

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

// Create sends a native completion request and waits for the complete result.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, http.MethodPost, "/completion", "completion.create", input, &output)
	return output, err
}

// CreateStream sends a native completion request and returns incremental SSE events.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[Response], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/completion", "completion.stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[Response], 0), nil
}
