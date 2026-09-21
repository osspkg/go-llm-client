// Package request provides typed request API operations.
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
	var o T
	target, ok := any(&o).(easyjson.Unmarshaler)
	if !ok {
		return o, errors.New("model does not implement easyjson unmarshaler")
	}
	if err := codec.Unmarshal(data, target); err != nil {
		return o, err
	}
	return o, nil
}

// JSON performs the JSON operation.
func JSON(ctx context.Context, c *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, output easyjson.Unmarshaler) error { //nolint:revive // low-level request metadata stays together.
	b, err := codec.Marshal(input)
	if err != nil {
		return err
	}
	d, err := c.Request(ctx, method, endpoint, operation, b, "application/json")
	if err != nil {
		return capabilityError(err, endpoint, operation)
	}
	if output == nil {
		return nil
	}
	if err := codec.Unmarshal(d, output); err != nil {
		return &llmerrors.DecodeError{Operation: operation, Cause: err}
	}
	return nil
}

// Stream performs the Stream operation.
func Stream(ctx context.Context, c *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, accept string) (io.ReadCloser, error) { //nolint:revive // low-level request metadata stays together.
	b, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	body, err := c.Stream(ctx, method, endpoint, operation, b, "application/json", accept)
	if err != nil {
		return nil, capabilityError(err, endpoint, operation)
	}
	return body, nil
}

func capabilityError(err error, endpoint, operation string) error {
	var httpErr *llmerrors.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &llmerrors.CapabilityError{Operation: operation, Endpoint: endpoint, Cause: err}
	}
	return err
}
