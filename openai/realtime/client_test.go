package realtime_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"go.osspkg.com/llm-client/openai"
	"go.osspkg.com/llm-client/openai/realtime"
	"go.osspkg.com/llm-client/pkg/auth"
)

func TestClient_ConnectSendReceiveAndClose(t *testing.T) {
	t.Parallel()

	server := newServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v1/realtime" {
			t.Errorf("path = %q, want /v1/realtime", request.URL.Path)
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		if request.URL.Query().Get("model") != "gpt-realtime" {
			t.Errorf("model = %q, want gpt-realtime", request.URL.Query().Get("model"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("authorization header = %q", request.Header.Get("Authorization"))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		connection, err := websocket.Accept(writer, request, nil)
		if err != nil {
			t.Errorf("accept websocket: %v", err)
			return
		}
		defer func() { _ = connection.Close(websocket.StatusNormalClosure, "") }()

		messageType, data, err := connection.Read(request.Context())
		if err != nil {
			t.Errorf("read client event: %v", err)
			return
		}
		if messageType != websocket.MessageText {
			t.Errorf("message type = %v, want text", messageType)
			return
		}
		if !strings.Contains(string(data), `"type":"session.update"`) {
			t.Errorf("client event = %s", data)
			return
		}
		if err := connection.Write(request.Context(), websocket.MessageText, []byte(`{"type":"session.created","session":{"id":"sess_123","model":"gpt-realtime"}}`)); err != nil {
			t.Errorf("write server event: %v", err)
		}
	}))

	client, err := openai.New(
		openai.WithBaseURL(server.URL+"/v1"),
		openai.WithAuthProvider(auth.StaticBearer("test-token")),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	session, err := client.Realtime().Connect(t.Context(), "gpt-realtime")
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			t.Errorf("close: %v", closeErr)
		}
	}()

	if err := session.Send(t.Context(), realtime.Event{Type: "session.update"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	event, err := session.Receive(t.Context())
	if err != nil {
		t.Fatalf("receive: %v", err)
	}
	if event.Type != "session.created" || event.Session == nil || event.Session.ID != "sess_123" {
		t.Fatalf("event = %#v, want session.created for sess_123", event)
	}
}

func TestSession_ReceiveHonorsCancellation(t *testing.T) {
	t.Parallel()

	connected := make(chan struct{})
	server := newServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		connection, err := websocket.Accept(writer, request, nil)
		if err != nil {
			t.Errorf("accept websocket: %v", err)
			return
		}
		defer func() { _ = connection.Close(websocket.StatusNormalClosure, "") }()
		close(connected)
		<-request.Context().Done()
	}))

	client, err := openai.New(openai.WithBaseURL(server.URL + "/v1"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	session, err := client.Realtime().Connect(t.Context(), "gpt-realtime")
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() { _ = session.Close() }()

	select {
	case <-connected:
	case <-time.After(time.Second):
		t.Fatal("server did not accept connection")
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = session.Receive(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("receive error = %v, want context cancellation", err)
	}
}

func newServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("local TCP listener is not permitted: %v", err)
	}
	server := httptest.NewUnstartedServer(handler)
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	return server
}
