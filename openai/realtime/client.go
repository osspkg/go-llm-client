// Package realtime implements the OpenAI Realtime WebSocket bounded context.
package realtime

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/transport"
	llmwebsocket "go.osspkg.com/llm-client/pkg/websocket"
)

// Client opens Realtime sessions.
type Client struct{ transport *transport.Client }

// New creates a Realtime client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Connect opens a Realtime session for a model.
func (client *Client) Connect(ctx context.Context, model string) (*Session, error) {
	if model == "" {
		return nil, errors.New("model is required")
	}
	endpoint := path.Join("/realtime")
	query := url.Values{"model": []string{model}}
	endpoint += "?" + query.Encode()
	target, err := client.transport.WebSocketURL(endpoint)
	if err != nil {
		return nil, err
	}
	headers, err := client.transport.WebSocketHeaders(ctx, "realtime.connect")
	if err != nil {
		return nil, err
	}
	connection, err := llmwebsocket.Dial(ctx, target, client.transport.HTTPClient(), headers)
	if err != nil {
		return nil, fmt.Errorf("connect realtime: %w", err)
	}
	return &Session{connection: connection}, nil
}

// Session is an explicitly owned Realtime WebSocket session.
type Session struct {
	connection *llmwebsocket.Session
}

// Send sends one typed client event.
func (session *Session) Send(ctx context.Context, event Event) error {
	if session == nil || session.connection == nil {
		return errors.New("send realtime event: session is nil")
	}
	data, err := codec.Marshal(event)
	if err != nil {
		return err
	}
	if err := session.connection.SendText(ctx, data); err != nil {
		return fmt.Errorf("send realtime event: %w", err)
	}
	return nil
}

// Receive receives and decodes one server event.
func (session *Session) Receive(ctx context.Context) (Event, error) {
	var output Event
	if session == nil || session.connection == nil {
		return output, errors.New("receive realtime event: session is nil")
	}
	data, err := session.connection.ReceiveText(ctx)
	if err != nil {
		return output, fmt.Errorf("receive realtime event: %w", err)
	}
	return request.Decode[Event](data)
}

// Close closes the WebSocket session.
func (session *Session) Close() error {
	if session == nil || session.connection == nil {
		return nil
	}
	return session.connection.Close()
}
