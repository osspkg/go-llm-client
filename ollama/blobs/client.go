// Package blobs implements Ollama blob existence and upload operations.
package blobs

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls blob endpoints.
type Client struct{ transport *transport.Client }

// New creates a blobs client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Exists reports whether a digest is available.
func (client *Client) Exists(ctx context.Context, digest string) (bool, error) {
	if digest == "" {
		return false, errors.New("digest is required")
	}
	_, err := client.transport.Request(ctx, http.MethodHead, "/api/blobs/"+url.PathEscape(digest), "ollama.blobs.exists", nil, "")
	if err != nil {
		var httpErr *llmerrors.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Upload streams a blob to its digest endpoint.
func (client *Client) Upload(ctx context.Context, digest string, body io.Reader) error {
	if digest == "" || body == nil {
		return errors.New("digest and body are required")
	}
	_, err := client.transport.RequestReader(ctx, http.MethodPost, "/api/blobs/"+url.PathEscape(digest), "ollama.blobs.upload", body, "application/octet-stream")
	return err
}
