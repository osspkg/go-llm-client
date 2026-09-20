// Package stream contains bounded typed streaming readers for provider protocols.
package stream

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"

	llmerrors "go.osspkg.com/llm-client/pkg/errors"
)

const defaultMaxEventSize int64 = 4 << 20

// Iterator is a synchronous, owned stream of typed values.
type Iterator[T any] interface {
	Next(ctx context.Context) bool
	Value() T
	Err() error
	Close() error
}

// ForEach consumes an iterator and closes it on return.
func ForEach[T any](ctx context.Context, iterator Iterator[T], fn func(T) error) error {
	if iterator == nil {
		return llmerrors.ErrInvalidRequest
	}
	defer func() { _ = iterator.Close() }()
	for iterator.Next(ctx) {
		if err := fn(iterator.Value()); err != nil {
			return err
		}
	}
	return iterator.Err()
}

// JSONDecoder decodes one provider event from JSON bytes.
type JSONDecoder[T any] func([]byte) (T, error)

// NDJSON reads one JSON value per line.
type NDJSON[T any] struct {
	reader    *bufio.Reader
	body      io.ReadCloser
	decode    JSONDecoder[T]
	maxSize   int64
	value     T
	err       error
	closed    bool
	closeOnce bool
}

// NewNDJSON creates a bounded newline-delimited JSON iterator.
func NewNDJSON[T any](body io.ReadCloser, decode JSONDecoder[T], maxSize int64) *NDJSON[T] {
	if maxSize <= 0 {
		maxSize = defaultMaxEventSize
	}
	reader := io.Reader(strings.NewReader(""))
	if body != nil {
		reader = body
	}
	stream := &NDJSON[T]{
		reader:  bufio.NewReader(reader),
		body:    body,
		decode:  decode,
		maxSize: maxSize,
	}
	if body == nil || decode == nil {
		stream.err = llmerrors.ErrInvalidRequest
	}
	return stream
}

// Next reads the next non-empty JSON line.
func (s *NDJSON[T]) Next(ctx context.Context) bool {
	if s.closed || s.err != nil {
		return false
	}
	if ctx == nil {
		s.err = llmerrors.ErrInvalidRequest
		_ = s.Close()
		return false
	}
	if err := ctx.Err(); err != nil {
		s.err = err
		_ = s.Close()
		return false
	}
	for {
		line, err := s.readLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				_ = s.Close()
				return false
			}
			s.err = err
			_ = s.Close()
			return false
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		value, decodeErr := s.decode(line)
		if decodeErr != nil {
			s.err = &llmerrors.DecodeError{Cause: decodeErr}
			_ = s.Close()
			return false
		}
		s.value = value
		return true
	}
}

// Value returns the most recently decoded value.
func (s *NDJSON[T]) Value() T { return s.value }

// Err returns the terminal stream error, if any.
func (s *NDJSON[T]) Err() error { return s.err }

// Close releases the underlying response body.
func (s *NDJSON[T]) Close() error {
	if s.closeOnce {
		return nil
	}
	s.closeOnce = true
	s.closed = true
	if s.body == nil {
		return nil
	}
	return s.body.Close()
}

func (s *NDJSON[T]) readLine() ([]byte, error) {
	var line []byte
	for {
		part, isPrefix, err := s.reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) && len(line) > 0 {
				return line, nil
			}
			return nil, err
		}
		if int64(len(line)+len(part)) > s.maxSize {
			return nil, llmerrors.ErrBodyTooLarge
		}
		line = append(line, part...)
		if !isPrefix {
			return line, nil
		}
	}
}

// SSE reads Server-Sent Events whose data payload is JSON.
type SSE[T any] struct {
	reader    *bufio.Reader
	body      io.ReadCloser
	decode    JSONDecoder[T]
	maxSize   int64
	value     T
	err       error
	closed    bool
	closeOnce bool
	done      bool
}

// NewSSE creates a bounded Server-Sent Events iterator.
func NewSSE[T any](body io.ReadCloser, decode JSONDecoder[T], maxSize int64) *SSE[T] {
	if maxSize <= 0 {
		maxSize = defaultMaxEventSize
	}
	reader := io.Reader(strings.NewReader(""))
	if body != nil {
		reader = body
	}
	stream := &SSE[T]{
		reader:  bufio.NewReader(reader),
		body:    body,
		decode:  decode,
		maxSize: maxSize,
	}
	if body == nil || decode == nil {
		stream.err = llmerrors.ErrInvalidRequest
	}
	return stream
}

// Next reads the next JSON data event. The terminal [DONE] event ends the stream.
func (s *SSE[T]) Next(ctx context.Context) bool {
	if s.closed || s.err != nil || s.done {
		return false
	}
	if ctx == nil {
		s.err = llmerrors.ErrInvalidRequest
		_ = s.Close()
		return false
	}
	if err := ctx.Err(); err != nil {
		s.err = err
		_ = s.Close()
		return false
	}
	var data []byte
	for {
		line, err := s.readLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(data) == 0 {
					_ = s.Close()
					return false
				}
				break
			}
			s.err = err
			_ = s.Close()
			return false
		}
		if len(line) == 0 {
			if len(data) == 0 {
				continue
			}
			break
		}
		if line[0] == ':' {
			continue
		}
		if !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}
		chunk := bytes.TrimSpace(line[len("data:"):])
		if int64(len(data)+len(chunk)) > s.maxSize {
			s.err = llmerrors.ErrBodyTooLarge
			_ = s.Close()
			return false
		}
		if len(data) > 0 {
			data = append(data, '\n')
		}
		data = append(data, chunk...)
	}
	if strings.TrimSpace(string(data)) == "[DONE]" {
		s.done = true
		_ = s.Close()
		return false
	}
	value, err := s.decode(data)
	if err != nil {
		s.err = &llmerrors.DecodeError{Cause: err}
		_ = s.Close()
		return false
	}
	s.value = value
	return true
}

// Value returns the most recently decoded value.
func (s *SSE[T]) Value() T { return s.value }

// Err returns the terminal stream error, if any.
func (s *SSE[T]) Err() error { return s.err }

// Close releases the underlying response body.
func (s *SSE[T]) Close() error {
	if s.closeOnce {
		return nil
	}
	s.closeOnce = true
	s.closed = true
	if s.body == nil {
		return nil
	}
	return s.body.Close()
}

func (s *SSE[T]) readLine() ([]byte, error) {
	line, err := s.reader.ReadBytes('\n')
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	if errors.Is(err, io.EOF) && len(line) > 0 {
		return line, nil
	}
	if int64(len(line)) > s.maxSize {
		return nil, llmerrors.ErrBodyTooLarge
	}
	return line, err
}
