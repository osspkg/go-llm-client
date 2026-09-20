// Package auth contains request-scoped authentication helpers.
package auth

import (
	"context"
	"net/http"
)

// RequestMeta describes the operation for which headers are requested.
type RequestMeta struct {
	Method string
	URL    string
	// Domain is the destination hostname. It never includes credentials, paths,
	// query parameters, or a port.
	Domain    string
	Operation string
	Streaming bool
	WebSocket bool
}

// HeaderProvider supplies headers for one outbound request.
// The returned header map is copied before it is applied to a request.
type HeaderProvider func(context.Context, RequestMeta) (http.Header, error)

// StaticBearer returns a provider that adds a bearer token to every request.
func StaticBearer(token string) HeaderProvider {
	return func(context.Context, RequestMeta) (http.Header, error) {
		headers := make(http.Header)
		if token != "" {
			headers.Set("Authorization", "Bearer "+token)
		}
		return headers, nil
	}
}

// Static returns a provider that returns a defensive copy of headers.
func Static(headers http.Header) HeaderProvider {
	copyHeaders := headers.Clone()
	return func(context.Context, RequestMeta) (http.Header, error) {
		return copyHeaders.Clone(), nil
	}
}
