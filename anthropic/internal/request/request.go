/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package request contains Anthropic request plumbing.
package request

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/mailru/easyjson"

	"go.osspkg.com/llm-client/pkg/codec"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

const (
	// Get is the HTTP GET method.
	Get = http.MethodGet
	// Post is the HTTP POST method.
	Post = http.MethodPost
	// Delete is the HTTP DELETE method.
	Delete = http.MethodDelete
)

// Decode performs the Decode operation.
func Decode[T any](data []byte) (T, error) {
	var output T
	target, ok := any(&output).(easyjson.Unmarshaler)
	if !ok {
		return output, errors.New("model does not implement easyjson unmarshaler")
	}
	if err := codec.Unmarshal(data, target); err != nil {
		return output, err
	}
	return output, nil
}

// JSON performs the JSON operation.
func JSON(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, output easyjson.Unmarshaler) error { //nolint:revive
	var body []byte
	var err error
	if input != nil {
		body, err = codec.Marshal(input)
		if err != nil {
			return err
		}
	}
	data, err := client.Request(ctx, method, endpoint, operation, body, "application/json")
	if err != nil {
		return err
	}
	if output == nil || len(data) == 0 {
		return nil
	}
	if err := codec.Unmarshal(data, output); err != nil {
		return &llmerrors.DecodeError{Operation: operation, Cause: err}
	}
	return nil
}

// Stream performs the Stream operation.
func Stream(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, accept string) (io.ReadCloser, error) { //nolint:revive
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	return client.Stream(ctx, method, endpoint, operation, body, "application/json", accept)
}
