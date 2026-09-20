package request

import (
	"context"
	"io"

	"go.osspkg.com/llm-client/pkg/transport"
)

// Reader executes a request from a bounded streaming body.
func Reader(ctx context.Context, client *transport.Client, method, endpoint, operation string, body io.Reader, contentType string) ([]byte, error) { //nolint:revive // provider helper mirrors transport metadata.
	return client.RequestReader(ctx, method, endpoint, operation, body, contentType)
}

// StreamReader opens a provider stream backed by a one-shot reader body.
func StreamReader(ctx context.Context, client *transport.Client, method, endpoint, operation string, body io.Reader, contentType, accept string) (io.ReadCloser, error) { //nolint:revive // provider helper mirrors transport metadata.
	return client.StreamReader(ctx, method, endpoint, operation, body, contentType, accept)
}
