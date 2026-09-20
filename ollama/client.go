// Package ollama provides a typed client for the Ollama HTTP API.
package ollama

import (
	"net/http"
	"time"

	"go.osspkg.com/llm-client/ollama/blobs"
	"go.osspkg.com/llm-client/ollama/chat"
	"go.osspkg.com/llm-client/ollama/embeddings"
	"go.osspkg.com/llm-client/ollama/generate"
	"go.osspkg.com/llm-client/ollama/models"
	"go.osspkg.com/llm-client/ollama/version"
	"go.osspkg.com/llm-client/pkg/auth"
	"go.osspkg.com/llm-client/pkg/transport"
)

const defaultBaseURL = "http://localhost:11434"

// Option configures an Ollama client.
type Option func(*config) error

type config struct{ transportOptions []transport.Option }

// Client is a concurrent-safe Ollama facade.
type Client struct {
	Generate   *generate.Client
	Chat       *chat.Client
	Embeddings *embeddings.Client
	Models     *models.Client
	Blobs      *blobs.Client
	Version    *version.Client
}

// New creates an Ollama client.
func New(options ...Option) (*Client, error) {
	configuration := config{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}
	client, err := transport.New(defaultBaseURL, configuration.transportOptions...)
	if err != nil {
		return nil, err
	}
	return &Client{
		Generate: generate.New(client), Chat: chat.New(client),
		Embeddings: embeddings.New(client), Models: models.New(client), Blobs: blobs.New(client), Version: version.New(client),
	}, nil
}

// WithBaseURL changes the Ollama endpoint.
func WithBaseURL(baseURL string) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithBaseURL(baseURL))
		return nil
	}
}

// WithAuthProvider configures optional per-request headers.
func WithAuthProvider(provider auth.HeaderProvider) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithAuthProvider(provider))
		return nil
	}
}

// WithHTTPClient supplies a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithHTTPClient(client))
		return nil
	}
}

// WithMaxRequestBody limits streamed and buffered request bodies.
func WithMaxRequestBody(limit int64) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithMaxRequestBody(limit))
		return nil
	}
}

// WithRequestTimeout sets the default request timeout.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithRequestTimeout(timeout))
		return nil
	}
}
