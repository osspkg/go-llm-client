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

// Decode decodes a generated native provider model.
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

// JSON executes a native JSON request and decodes its generated response.
func JSON(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, output easyjson.Unmarshaler) error { //nolint:revive // low-level request metadata stays together.
	body, err := codec.Marshal(input)
	if err != nil {
		return err
	}
	data, err := client.Request(ctx, method, endpoint, operation, body, "application/json")
	if err != nil {
		return capabilityError(err, endpoint, operation)
	}
	if output == nil {
		return nil
	}
	if err := codec.Unmarshal(data, output); err != nil {
		return &llmerrors.DecodeError{Operation: operation, Cause: err}
	}
	return nil
}

// Stream opens a native SSE response stream.
func Stream(ctx context.Context, client *transport.Client, method, endpoint, operation string, input easyjson.Marshaler, accept string) (io.ReadCloser, error) { //nolint:revive // low-level request metadata stays together.
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	responseBody, err := client.Stream(ctx, method, endpoint, operation, body, "application/json", accept)
	if err != nil {
		return nil, capabilityError(err, endpoint, operation)
	}
	return responseBody, nil
}

func capabilityError(err error, endpoint, operation string) error {
	var httpErr *llmerrors.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &llmerrors.CapabilityError{Operation: operation, Endpoint: endpoint, Cause: err}
	}
	return err
}
