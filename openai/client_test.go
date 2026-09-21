package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"go.osspkg.com/llm-client/openai"
	"go.osspkg.com/llm-client/openai/capability"
	"go.osspkg.com/llm-client/openai/chat"
	"go.osspkg.com/llm-client/openai/responses"
	"go.osspkg.com/llm-client/pkg/auth"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestChatAndResponsesStreamsUseTypedTransport(t *testing.T) {
	var authCalls atomic.Int32
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("authorization = %q", request.Header.Get("Authorization"))
		}
		switch request.URL.Path {
		case "/v1/chat/completions":
			var input chat.Request
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
				t.Errorf("decode chat request: %v", err)
			}
			if input.Model != "test-model" {
				t.Errorf("model = %q", input.Model)
			}
			return response(http.StatusOK, "application/json", `{"id":"chat-1","object":"chat.completion","model":"test-model","choices":[{"index":0,"message":{"role":"assistant","content":"hello"}}]}`), nil
		case "/v1/responses":
			return response(http.StatusOK, "text/event-stream", "data: {\"type\":\"response.output_text.delta\",\"delta\":\"hi\"}\n\ndata: [DONE]\n\n"), nil
		default:
			return response(http.StatusNotFound, "text/plain", "not found"), nil
		}
	})}

	client, err := openai.New(
		openai.WithBaseURL("https://example.test/v1"),
		openai.WithHTTPClient(httpClient),
		openai.WithAuthProvider(func(context.Context, auth.RequestMeta) (http.Header, error) {
			authCalls.Add(1)
			return http.Header{"Authorization": []string{"Bearer test-token"}}, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Chat().Create(context.Background(), chat.Request{Model: "test-model"})
	if err != nil {
		t.Fatal(err)
	}
	if response.ID != "chat-1" || len(response.Choices) != 1 {
		t.Fatalf("response = %#v", response)
	}

	iterator, err := client.Responses().CreateStream(context.Background(), responses.Request{Model: "test-model"})
	if err != nil {
		t.Fatal(err)
	}
	if !iterator.Next(context.Background()) || iterator.Value().Delta != "hi" {
		t.Fatalf("stream value = %#v", iterator.Value())
	}
	next := iterator.Next(context.Background())
	if next || iterator.Err() != nil {
		t.Fatalf("stream terminal state: next=%v err=%v", next, iterator.Err())
	}
	if authCalls.Load() != 2 {
		t.Fatalf("auth calls = %d", authCalls.Load())
	}
}

func TestDisabledCapabilityStopsBeforeTransport(t *testing.T) {
	client, err := openai.New(
		openai.WithCapability(capability.Chat, false),
		openai.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("transport should not be called")
			return nil, errors.New("unexpected transport call")
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Chat().Create(context.Background(), chat.Request{Model: "test-model"})
	if !errors.Is(err, llmerrors.ErrInvalidRequest) {
		t.Fatalf("error = %v", err)
	}
}

func TestStoredResourceOperationsUseExpectedRoutes(t *testing.T) {
	t.Parallel()

	expected := []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/v1/chat/completions/chat_123"},
		{method: http.MethodPatch, path: "/v1/chat/completions/chat_123"},
		{method: http.MethodDelete, path: "/v1/chat/completions/chat_123"},
		{method: http.MethodGet, path: "/v1/chat/completions/chat_123/messages"},
		{method: http.MethodGet, path: "/v1/responses/resp_123"},
		{method: http.MethodDelete, path: "/v1/responses/resp_123"},
		{method: http.MethodPost, path: "/v1/responses/resp_123/cancel"},
		{method: http.MethodGet, path: "/v1/responses/resp_123/input_items"},
		{method: http.MethodGet, path: "/v1/batches"},
	}
	index := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if index >= len(expected) {
			t.Fatalf("unexpected request %s %s", request.Method, request.URL.Path)
		}
		want := expected[index]
		index++
		if request.Method != want.method || request.URL.Path != want.path {
			t.Errorf("request = %s %s, want %s %s", request.Method, request.URL.Path, want.method, want.path)
		}
		switch request.URL.Path {
		case "/v1/chat/completions/chat_123":
			if request.Method == http.MethodDelete {
				return response(http.StatusOK, "application/json", `{"id":"chat_123","object":"chat.completion.deleted","deleted":true}`), nil
			}
			return response(http.StatusOK, "application/json", `{"id":"chat_123","object":"chat.completion","model":"test","choices":[]}`), nil
		case "/v1/chat/completions/chat_123/messages":
			return response(http.StatusOK, "application/json", `{"object":"list","data":[],"has_more":false}`), nil
		case "/v1/responses/resp_123":
			if request.Method == http.MethodDelete {
				return response(http.StatusOK, "application/json", `{"id":"resp_123","object":"response","deleted":true}`), nil
			}
			return response(http.StatusOK, "application/json", `{"id":"resp_123","object":"response","model":"test","status":"completed"}`), nil
		case "/v1/responses/resp_123/cancel":
			return response(http.StatusOK, "application/json", `{"id":"resp_123","object":"response","model":"test","status":"cancelled"}`), nil
		case "/v1/responses/resp_123/input_items":
			return response(http.StatusOK, "application/json", `{"object":"list","data":[],"has_more":false}`), nil
		case "/v1/batches":
			return response(http.StatusOK, "application/json", `{"object":"list","data":[],"has_more":false}`), nil
		default:
			return response(http.StatusNotFound, "application/json", `{}`), nil
		}
	})}
	client, err := openai.New(openai.WithBaseURL("https://example.test/v1"), openai.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := client.Chat().Get(ctx, "chat_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Chat().Update(ctx, "chat_123", chat.UpdateRequest{Metadata: map[string]string{"key": "value"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Chat().Delete(ctx, "chat_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Chat().Messages(ctx, "chat_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Responses().Get(ctx, "resp_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Responses().Delete(ctx, "resp_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Responses().Cancel(ctx, "resp_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Responses().InputItems(ctx, "resp_123"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Batches().List(ctx); err != nil {
		t.Fatal(err)
	}
	if index != len(expected) {
		t.Fatalf("requests = %d, want %d", index, len(expected))
	}
}

func response(status int, contentType, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(strings.NewReader(body))}
}
