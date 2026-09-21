/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package request contains OpenAI domain request plumbing.
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

// Decode decodes a generated model value.
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

// JSON executes a JSON request and decodes a generated response model.
func JSON(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, output easyjson.Unmarshaler) error { //nolint:revive // provider helper mirrors transport metadata.
	body, err := codec.Marshal(input)
	if err != nil {
		return err
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

// Bytes executes a request with a caller-provided body and returns bounded bytes.
func Bytes(ctx context.Context, client *transport.Client, method, endpoint, operation string, body []byte, contentType string) ([]byte, error) { //nolint:revive // provider helper mirrors transport metadata.
	return client.Request(ctx, method, endpoint, operation, body, contentType)
}

// Stream opens a provider stream.
func Stream(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, accept string) (io.ReadCloser, error) { //nolint:revive // provider helper mirrors transport metadata.
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	return client.Stream(ctx, method, endpoint, operation, body, "application/json", accept)
}

// Method aliases keep provider methods readable.
// HTTP method aliases keep provider methods readable.
const (
	Get    = http.MethodGet
	Post   = http.MethodPost
	Delete = http.MethodDelete
	Patch  = http.MethodPatch
)
