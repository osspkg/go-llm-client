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
