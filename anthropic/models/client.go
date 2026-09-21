/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package models provides typed models API operations.
package models

import (
	"context"
	"net/url"
	"strconv"

	"go.osspkg.com/llm-client/anthropic/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List performs the List operation.
func (client *Client) List(ctx context.Context, limit int, after, before string) (ListResponse, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		query.Set("after_id", after)
	}
	if before != "" {
		query.Set("before_id", before)
	}
	endpoint := "/models"
	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, endpoint, "models.list", nil, &output)
	return output, err
}

// Get performs the Get operation.
func (client *Client) Get(ctx context.Context, id string) (Model, error) {
	var output Model
	err := request.JSON(ctx, client.transport, request.Get, "/models/"+url.PathEscape(id), "models.get", nil, &output)
	return output, err
}
