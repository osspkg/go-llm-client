// Package anthropic provides a typed client for the stable Anthropic API.
package anthropic

import (
	"context"
	"net/http"
	"time"

	"go.osspkg.com/llm-client/anthropic/batches"
	"go.osspkg.com/llm-client/anthropic/files"
	"go.osspkg.com/llm-client/anthropic/messages"
	"go.osspkg.com/llm-client/anthropic/models"
	"go.osspkg.com/llm-client/pkg/auth"
	"go.osspkg.com/llm-client/pkg/transport"
)

const defaultBaseURL = "https://api.anthropic.com/v1"

type (
	// Option configures an Anthropic client.
	Option func(*config) error
	config struct{ transportOptions []transport.Option }
)

// Client describes the Client API value.
type Client struct {
	Messages *messages.Client
	Models   *models.Client
	Files    *files.Client
	Batches  *batches.Client
}

// New creates a client on shared transport.
func New(options ...Option) (*Client, error) {
	c := config{}
	c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-version", "2023-06-01"))
	for _, option := range options {
		if option != nil {
			if err := option(&c); err != nil {
				return nil, err
			}
		}
	}
	t, err := transport.New(defaultBaseURL, c.transportOptions...)
	if err != nil {
		return nil, err
	}
	return &Client{Messages: messages.New(t), Models: models.New(t), Files: files.New(t), Batches: batches.New(t)}, nil
}

// WithBaseURL performs the WithBaseURL operation.
func WithBaseURL(baseURL string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithBaseURL(baseURL))
		return nil
	}
}

// WithHTTPClient performs the WithHTTPClient operation.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHTTPClient(client))
		return nil
	}
}

// WithAuthProvider performs the WithAuthProvider operation.
func WithAuthProvider(provider auth.HeaderProvider) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithAuthProvider(provider))
		return nil
	}
}

// WithAPIKey performs the WithAPIKey operation.
func WithAPIKey(key string) Option {
	return WithAuthProvider(func(_ context.Context, _ auth.RequestMeta) (http.Header, error) {
		h := make(http.Header)
		if key != "" {
			h.Set("x-api-key", key)
		}
		return h, nil
	})
}

// WithBearerToken configures Authorization: Bearer authentication.
func WithBearerToken(token string) Option {
	return WithAuthProvider(func(_ context.Context, _ auth.RequestMeta) (http.Header, error) {
		h := make(http.Header)
		if token != "" {
			h.Set("Authorization", "Bearer "+token)
		}
		return h, nil
	})
}

// WithVersion performs the WithVersion operation.
func WithVersion(version string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-version", version))
		return nil
	}
}

// WithWorkspace performs the WithWorkspace operation.
func WithWorkspace(id string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-workspace-id", id))
		return nil
	}
}

// WithBeta performs the WithBeta operation.
func WithBeta(value string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-beta", value))
		return nil
	}
}

// WithRequestTimeout performs the WithRequestTimeout operation.
func WithRequestTimeout(value time.Duration) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithRequestTimeout(value))
		return nil
	}
}

// WithStreamTimeout performs the WithStreamTimeout operation.
func WithStreamTimeout(value time.Duration) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithStreamTimeout(value))
		return nil
	}
}

// WithMaxResponseBody bounds buffered JSON responses and error payloads.
func WithMaxResponseBody(value int64) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithMaxResponseBody(value))
		return nil
	}
}

// WithMaxRequestBody bounds request payloads and file uploads.
func WithMaxRequestBody(value int64) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithMaxRequestBody(value))
		return nil
	}
}
