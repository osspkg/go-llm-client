// Package metrics provides the bounded native Prometheus metrics endpoint.
package metrics

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls native metrics operations.
type Client struct{ transport *transport.Client }

// New creates a metrics client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Get returns the bounded Prometheus exposition text. A model is optional for
// single-model servers and required by router mode.
func (client *Client) Get(ctx context.Context, model ...string) (string, error) {
	endpoint := "/metrics"
	if len(model) > 0 && model[0] != "" {
		endpoint += "?" + url.Values{"model": []string{model[0]}}.Encode()
	}
	data, err := client.transport.Request(ctx, http.MethodGet, endpoint, "metrics.get", nil, "")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
