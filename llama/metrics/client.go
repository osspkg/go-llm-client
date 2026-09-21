/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

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

// Get returns the bounded Prometheus exposition text. Pass an empty model for
// single-model servers; router mode can use the model query parameter.
func (client *Client) Get(ctx context.Context, model string) (string, error) {
	endpoint := "/metrics"
	if model != "" {
		endpoint += "?" + url.Values{"model": []string{model}}.Encode()
	}
	data, err := client.transport.Request(ctx, http.MethodGet, endpoint, "metrics.get", nil, "")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
