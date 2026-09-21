// Package llama provides typed clients for native llama.cpp server endpoints.
package llama

import (
	"net/http"
	"time"

	"go.osspkg.com/llm-client/llama/completions"
	"go.osspkg.com/llm-client/llama/embeddings"
	"go.osspkg.com/llm-client/llama/lora"
	"go.osspkg.com/llm-client/llama/metrics"
	"go.osspkg.com/llm-client/llama/models"
	"go.osspkg.com/llm-client/llama/rerank"
	"go.osspkg.com/llm-client/llama/server"
	"go.osspkg.com/llm-client/llama/slots"
	"go.osspkg.com/llm-client/llama/templates"
	"go.osspkg.com/llm-client/llama/tokenization"
	"go.osspkg.com/llm-client/pkg/auth"
	"go.osspkg.com/llm-client/pkg/transport"
)

const defaultBaseURL = "http://localhost:8080"

type (
	// Option configures a native llama.cpp client.
	Option func(*config) error
	config struct{ options []transport.Option }
)

// Client exposes native llama.cpp bounded contexts.
type Client struct {
	Completions  *completions.Client
	Embeddings   *embeddings.Client
	Tokenization *tokenization.Client
	Templates    *templates.Client
	Server       *server.Client
	Lora         *lora.Client
	Metrics      *metrics.Client
	Models       *models.Client
	Slots        *slots.Client
	Rerank       *rerank.Client
}

// New creates a client on shared transport.
func New(opts ...Option) (*Client, error) {
	c := config{}
	for _, o := range opts {
		if o != nil {
			if err := o(&c); err != nil {
				return nil, err
			}
		}
	}
	t, err := transport.New(defaultBaseURL, c.options...)
	if err != nil {
		return nil, err
	}
	return &Client{Completions: completions.New(t), Embeddings: embeddings.New(t), Tokenization: tokenization.New(t), Templates: templates.New(t), Server: server.New(t), Lora: lora.New(t), Metrics: metrics.New(t), Models: models.New(t), Slots: slots.New(t), Rerank: rerank.New(t)}, nil
}

// WithBaseURL performs the WithBaseURL operation.
func WithBaseURL(v string) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithBaseURL(v)); return nil }
}

// WithHTTPClient performs the WithHTTPClient operation.
func WithHTTPClient(v *http.Client) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithHTTPClient(v)); return nil }
}

// WithAuthProvider performs the WithAuthProvider operation.
func WithAuthProvider(v auth.HeaderProvider) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithAuthProvider(v)); return nil }
}

// WithRequestTimeout performs the WithRequestTimeout operation.
func WithRequestTimeout(v time.Duration) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithRequestTimeout(v)); return nil }
}

// WithStreamTimeout performs the WithStreamTimeout operation.
func WithStreamTimeout(v time.Duration) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithStreamTimeout(v)); return nil }
}

// WithMaxResponseBody bounds buffered JSON and metrics responses.
func WithMaxResponseBody(v int64) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithMaxResponseBody(v)); return nil }
}

// WithMaxRequestBody bounds request payloads and reader-backed uploads.
func WithMaxRequestBody(v int64) Option {
	return func(c *config) error { c.options = append(c.options, transport.WithMaxRequestBody(v)); return nil }
}
