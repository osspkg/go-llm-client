/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package request contains Ollama domain request plumbing.
package request

import (
	"context"
	"errors"

	"github.com/mailru/easyjson"

	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/transport"
)

// HTTP method aliases keep provider methods readable.
const (
	Get    = "GET"
	Post   = "POST"
	Delete = "DELETE"
)

// JSON executes a JSON request and decodes a generated response.
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
	return codec.Unmarshal(data, output)
}

// Decode decodes a generated model.
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
