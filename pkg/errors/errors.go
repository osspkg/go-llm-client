/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package errors defines errors shared by provider clients.
package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig indicates invalid client configuration.
	ErrInvalidConfig = errors.New("llm client: invalid config")
	// ErrInvalidRequest indicates invalid caller input.
	ErrInvalidRequest = errors.New("llm client: invalid request")
	// ErrBodyTooLarge indicates that a configured body limit was exceeded.
	ErrBodyTooLarge = errors.New("llm client: response body too large")
	// ErrProtocol indicates malformed provider protocol data.
	ErrProtocol = errors.New("llm client: protocol error")
	// ErrStreamClosed indicates that a stream was already closed.
	ErrStreamClosed = errors.New("llm client: stream closed")
)

// HTTPError describes a non-successful HTTP response without retaining its body.
type HTTPError struct {
	StatusCode   int
	RequestID    string
	ProviderCode string
	ProviderType string
	Param        string
	Message      string
}

// Error implements error.
func (e *HTTPError) Error() string {
	if e == nil {
		return "llm client: http error"
	}
	if e.Message == "" {
		return fmt.Sprintf("llm client: http status %d", e.StatusCode)
	}
	return fmt.Sprintf("llm client: http status %d: %s", e.StatusCode, e.Message)
}

// Unwrap exposes the shared invalid-request-independent HTTP error category.
func (e *HTTPError) Unwrap() error { return ErrProtocol }

// ValidationError identifies a field that failed request validation.
type ValidationError struct {
	Field string
	Msg   string
}

// Error implements error.
func (e *ValidationError) Error() string {
	if e == nil {
		return "llm client: invalid request"
	}
	if e.Field == "" {
		return "llm client: invalid request: " + e.Msg
	}
	return "llm client: invalid request: " + e.Field + ": " + e.Msg
}

// Unwrap exposes ErrInvalidRequest for errors.Is.
func (e *ValidationError) Unwrap() error { return ErrInvalidRequest }

// DecodeError identifies a response that could not be decoded.
type DecodeError struct {
	Operation string
	Cause     error
}

// Error implements error.
func (e *DecodeError) Error() string {
	if e == nil {
		return "llm client: decode error"
	}
	if e.Operation == "" {
		return fmt.Sprintf("llm client: decode response: %v", e.Cause)
	}
	return fmt.Sprintf("llm client: decode %s response: %v", e.Operation, e.Cause)
}

// Unwrap exposes the underlying decode failure.
func (e *DecodeError) Unwrap() error { return e.Cause }

// CapabilityError identifies an endpoint that the selected provider build does
// not expose. Cause retains the bounded underlying protocol or HTTP error.
type CapabilityError struct {
	// Operation is the provider operation name used by the client.
	Operation string
	// Endpoint is the relative endpoint that was unavailable.
	Endpoint string
	// Cause is the underlying bounded protocol error.
	Cause error
}

// Error implements error.
func (e *CapabilityError) Error() string {
	if e == nil {
		return "llm client: capability error"
	}
	if e.Operation == "" {
		return "llm client: unsupported capability: " + e.Endpoint
	}
	return fmt.Sprintf("llm client: unsupported capability %s: %s", e.Operation, e.Endpoint)
}

// Unwrap exposes both the protocol category and the bounded cause.
func (e *CapabilityError) Unwrap() []error { return []error{ErrProtocol, e.Cause} }
