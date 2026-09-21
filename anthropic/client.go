/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

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

// Client is a concurrent-safe facade over Anthropic bounded contexts.
type Client struct {
	messages *messages.Client
	models   *models.Client
	files    *files.Client
	batches  *batches.Client
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
	return &Client{messages: messages.New(t), models: models.New(t), files: files.New(t), batches: batches.New(t)}, nil
}

// Messages returns the immutable Anthropic Messages domain client.
func (client *Client) Messages() *messages.Client { return client.messages }

// Models returns the immutable Anthropic Models domain client.
func (client *Client) Models() *models.Client { return client.models }

// Files returns the immutable Anthropic Files domain client.
func (client *Client) Files() *files.Client { return client.files }

// Batches returns the immutable Anthropic Message Batches domain client.
func (client *Client) Batches() *batches.Client { return client.batches }

// WithBaseURL changes the Anthropic endpoint, including for compatible servers.
func WithBaseURL(baseURL string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithBaseURL(baseURL))
		return nil
	}
}

// WithHTTPClient supplies a custom HTTP client, normally for tests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHTTPClient(client))
		return nil
	}
}

// WithAuthProvider configures per-request authentication headers.
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

// WithVersion sets the Anthropic API version header.
func WithVersion(version string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-version", version))
		return nil
	}
}

// WithWorkspace sets the Anthropic workspace header.
func WithWorkspace(id string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-workspace-id", id))
		return nil
	}
}

// WithBeta enables an Anthropic beta feature header.
func WithBeta(value string) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithHeader("anthropic-beta", value))
		return nil
	}
}

// WithRequestTimeout sets the default non-stream request timeout.
func WithRequestTimeout(value time.Duration) Option {
	return func(c *config) error {
		c.transportOptions = append(c.transportOptions, transport.WithRequestTimeout(value))
		return nil
	}
}

// WithStreamTimeout sets an optional total stream timeout.
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
