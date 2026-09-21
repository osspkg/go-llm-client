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
	completions  *completions.Client
	embeddings   *embeddings.Client
	tokenization *tokenization.Client
	templates    *templates.Client
	server       *server.Client
	loRA         *lora.Client
	metrics      *metrics.Client
	models       *models.Client
	slots        *slots.Client
	rerank       *rerank.Client
}

// New creates a client on shared transport.
func New(opts ...Option) (*Client, error) {
	configuration := config{}
	for _, option := range opts {
		if option != nil {
			if err := option(&configuration); err != nil {
				return nil, err
			}
		}
	}
	client, err := transport.New(defaultBaseURL, configuration.options...)
	if err != nil {
		return nil, err
	}
	return &Client{completions: completions.New(client), embeddings: embeddings.New(client), tokenization: tokenization.New(client), templates: templates.New(client), server: server.New(client), loRA: lora.New(client), metrics: metrics.New(client), models: models.New(client), slots: slots.New(client), rerank: rerank.New(client)}, nil
}

// Completions returns the immutable native completion domain client.
func (client *Client) Completions() *completions.Client { return client.completions }

// Embeddings returns the immutable native embeddings domain client.
func (client *Client) Embeddings() *embeddings.Client { return client.embeddings }

// Tokenization returns the immutable native tokenization domain client.
func (client *Client) Tokenization() *tokenization.Client { return client.tokenization }

// Templates returns the immutable native template domain client.
func (client *Client) Templates() *templates.Client { return client.templates }

// Server returns the immutable native server-management domain client.
func (client *Client) Server() *server.Client { return client.server }

// LoRA returns the immutable native LoRA domain client.
func (client *Client) LoRA() *lora.Client { return client.loRA }

// Metrics returns the immutable native metrics domain client.
func (client *Client) Metrics() *metrics.Client { return client.metrics }

// Models returns the immutable native model-router domain client.
func (client *Client) Models() *models.Client { return client.models }

// Slots returns the immutable native slots domain client.
func (client *Client) Slots() *slots.Client { return client.slots }

// Rerank returns the immutable native reranking domain client.
func (client *Client) Rerank() *rerank.Client { return client.rerank }

// WithBaseURL changes the native llama.cpp server endpoint.
func WithBaseURL(baseURL string) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithBaseURL(baseURL))
		return nil
	}
}

// WithHTTPClient supplies a custom HTTP client, normally for tests.
func WithHTTPClient(client *http.Client) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithHTTPClient(client))
		return nil
	}
}

// WithAuthProvider configures optional per-request headers.
func WithAuthProvider(provider auth.HeaderProvider) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithAuthProvider(provider))
		return nil
	}
}

// WithRequestTimeout sets the default non-stream request timeout.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithRequestTimeout(timeout))
		return nil
	}
}

// WithStreamTimeout sets an optional total stream timeout.
func WithStreamTimeout(timeout time.Duration) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithStreamTimeout(timeout))
		return nil
	}
}

// WithMaxResponseBody bounds buffered JSON and metrics responses.
func WithMaxResponseBody(limit int64) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithMaxResponseBody(limit))
		return nil
	}
}

// WithMaxRequestBody bounds request payloads and reader-backed uploads.
func WithMaxRequestBody(limit int64) Option {
	return func(configuration *config) error {
		configuration.options = append(configuration.options, transport.WithMaxRequestBody(limit))
		return nil
	}
}
