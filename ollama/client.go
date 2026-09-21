/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

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
	generate   *generate.Client
	chat       *chat.Client
	embeddings *embeddings.Client
	models     *models.Client
	blobs      *blobs.Client
	version    *version.Client
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
		generate: generate.New(client), chat: chat.New(client),
		embeddings: embeddings.New(client), models: models.New(client), blobs: blobs.New(client), version: version.New(client),
	}, nil
}

// Generate returns the immutable Ollama generate domain client.
func (client *Client) Generate() *generate.Client { return client.generate }

// Chat returns the immutable Ollama chat domain client.
func (client *Client) Chat() *chat.Client { return client.chat }

// Embeddings returns the immutable Ollama embeddings domain client.
func (client *Client) Embeddings() *embeddings.Client { return client.embeddings }

// Models returns the immutable Ollama model-management domain client.
func (client *Client) Models() *models.Client { return client.models }

// Blobs returns the immutable Ollama blob domain client.
func (client *Client) Blobs() *blobs.Client { return client.blobs }

// Version returns the immutable Ollama version domain client.
func (client *Client) Version() *version.Client { return client.version }

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
