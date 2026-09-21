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
func New(c *transport.Client) *Client { return &Client{c} }

// Create performs the Create operation.
func (c *Client) Create(ctx context.Context, in Request) (Response, error) {
	var o Response
	err := request.JSON(ctx, c.transport, http.MethodPost, "/completion", "completion.create", in, &o)
	return o, err
}

// CreateStream performs the CreateStream operation.
func (c *Client) CreateStream(ctx context.Context, in Request) (stream.Iterator[Response], error) {
	in.Stream = true
	b, err := request.Stream(ctx, c.transport, http.MethodPost, "/completion", "completion.stream", in, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(b, request.Decode[Response], 0), nil
}
