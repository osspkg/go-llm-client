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
func New(c *transport.Client) *Client { return &Client{c} }

// Apply performs the Apply operation.
func (c *Client) Apply(ctx context.Context, in Request) (Response, error) {
	var o Response
	err := request.JSON(ctx, c.transport, http.MethodPost, "/apply-template", "templates.apply", in, &o)
	return o, err
}
