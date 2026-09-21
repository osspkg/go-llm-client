// Package templates provides typed templates API operations.
package templates

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

// Apply formats chat messages with the server's native prompt template.
func (client *Client) Apply(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, http.MethodPost, "/apply-template", "templates.apply", input, &output)
	return output, err
}
