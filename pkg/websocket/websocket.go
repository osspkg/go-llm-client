// Package websocket provides a provider-independent typed-session transport
// boundary for WebSocket clients.
package websocket

import (
	"context"
	"errors"
	"net/http"
	"sync"

	coderwebsocket "github.com/coder/websocket"
)

// Session is an explicitly owned WebSocket connection. Writes are serialized;
// callers must coordinate reads so that at most one goroutine reads at a time.
type Session struct {
	connection *coderwebsocket.Conn
	writeMu    sync.Mutex
	readMu     sync.Mutex
	closeOnce  sync.Once
	closeErr   error
}

// Dial opens a text-capable WebSocket session.
func Dial(ctx context.Context, target string, client *http.Client, headers http.Header) (*Session, error) {
	if ctx == nil || target == "" {
		return nil, errors.New("websocket: context and target are required")
	}
	connection, response, err := coderwebsocket.Dial(ctx, target, &coderwebsocket.DialOptions{
		HTTPClient: client,
		HTTPHeader: headers,
	})
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return &Session{connection: connection}, nil
}

// SendText writes one text message. Concurrent writes are serialized.
func (session *Session) SendText(ctx context.Context, data []byte) error {
	if session == nil || session.connection == nil {
		return errors.New("websocket: session is nil")
	}
	session.writeMu.Lock()
	defer session.writeMu.Unlock()
	return session.connection.Write(ctx, coderwebsocket.MessageText, data)
}

// ReceiveText reads one text message. Concurrent reads are serialized.
func (session *Session) ReceiveText(ctx context.Context) ([]byte, error) {
	if session == nil || session.connection == nil {
		return nil, errors.New("websocket: session is nil")
	}
	session.readMu.Lock()
	defer session.readMu.Unlock()
	messageType, data, err := session.connection.Read(ctx)
	if err != nil {
		return nil, err
	}
	if messageType != coderwebsocket.MessageText {
		return nil, errors.New("websocket: unexpected binary message")
	}
	return data, nil
}

// Close performs an idempotent normal WebSocket close.
func (session *Session) Close() error {
	if session == nil || session.connection == nil {
		return nil
	}
	session.closeOnce.Do(func() {
		session.closeErr = session.connection.Close(coderwebsocket.StatusNormalClosure, "")
	})
	return session.closeErr
}
