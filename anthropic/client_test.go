package anthropic_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/anthropic"
	"go.osspkg.com/llm-client/anthropic/batches"
	"go.osspkg.com/llm-client/anthropic/messages"
	"go.osspkg.com/llm-client/pkg/auth"
	"go.osspkg.com/llm-client/pkg/stream"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestClientAddsAnthropicHeadersAndCallsMessages(t *testing.T) {
	client, err := anthropic.New(
		anthropic.WithAPIKey("secret"),
		anthropic.WithWorkspace("workspace"),
		anthropic.WithBeta("test-beta"),
		anthropic.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/v1/messages" || r.Header.Get("x-api-key") != "secret" || r.Header.Get("anthropic-version") != "2023-06-01" || r.Header.Get("anthropic-workspace-id") != "workspace" || r.Header.Get("anthropic-beta") != "test-beta" {
				t.Errorf("request = %s headers=%v", r.URL, r.Header)
			}
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"msg_1","type":"message","role":"assistant","content":[],"model":"claude","usage":{"input_tokens":1,"output_tokens":2}}`))}, nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	content, err := messages.TextContent("hello")
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Messages().Create(context.Background(), messages.Request{Model: "claude", MaxTokens: 16, Messages: []messages.Message{{Role: "user", Content: content}}})
	if err != nil {
		t.Fatal(err)
	}
	if response.ID != "msg_1" || response.Usage.OutputTokens != 2 {
		t.Fatalf("response = %#v", response)
	}
}

func TestAnthropicAuthReceivesDestinationDomain(t *testing.T) {
	var got auth.RequestMeta
	client, err := anthropic.New(
		anthropic.WithBaseURL("https://anthropic.example.test/v1"),
		anthropic.WithAuthProvider(func(_ context.Context, meta auth.RequestMeta) (http.Header, error) {
			got = meta
			return make(http.Header), nil
		}),
		anthropic.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"input_tokens":3}`))}, nil
		})}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Messages().CountTokens(context.Background(), messages.CountTokensRequest{Model: "claude"}); err != nil {
		t.Fatal(err)
	}
	if got.Domain != "anthropic.example.test" {
		t.Fatalf("domain = %q", got.Domain)
	}
}

func TestAnthropicMessageStreamAndBatchResults(t *testing.T) {
	client, err := anthropic.New(anthropic.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body string
		switch r.URL.Path {
		case "/v1/messages":
			body = "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{}}\n\n" +
				"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"hi\"}}\n\n" +
				"event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"}}\n\n" +
				"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
		case "/v1/messages/batches/batch_1/results":
			body = "{\"custom_id\":\"request_1\",\"result\":{\"type\":\"errored\",\"error\":{\"type\":\"invalid_request\",\"message\":\"bad input\"}}}\n"
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			body = `{}`
		}
		contentType := "text/event-stream"
		if strings.HasSuffix(r.URL.Path, "/results") {
			contentType = "application/jsonl"
		}
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	iterator, err := client.Messages().CreateStream(context.Background(), messages.Request{Model: "claude", MaxTokens: 8})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = iterator.Close() }()
	if !iterator.Next(context.Background()) || iterator.Value().Type != "message_start" {
		t.Fatalf("first stream event = %#v, err=%v", iterator.Value(), iterator.Err())
	}
	if !iterator.Next(context.Background()) || iterator.Value().Index != 0 {
		t.Fatalf("second stream event = %#v, err=%v", iterator.Value(), iterator.Err())
	}
	delta, err := iterator.Value().DecodeContentBlockDelta()
	if err != nil || delta.Text != "hi" {
		t.Fatalf("content delta = %#v, err=%v", delta, err)
	}
	if !iterator.Next(context.Background()) || iterator.Value().Type != "message_delta" {
		t.Fatalf("third stream event = %#v, err=%v", iterator.Value(), iterator.Err())
	}
	messageDelta, err := iterator.Value().DecodeMessageDelta()
	if err != nil || messageDelta.StopReason != "end_turn" {
		t.Fatalf("message delta = %#v, err=%v", messageDelta, err)
	}
	if !iterator.Next(context.Background()) || iterator.Value().Type != "message_stop" {
		t.Fatalf("fourth stream event = %#v, err=%v", iterator.Value(), iterator.Err())
	}
	next := iterator.Next(context.Background())
	if next || iterator.Err() != nil {
		t.Fatalf("stream terminal state: next=%v err=%v", next, iterator.Err())
	}
	results, err := client.Batches().Results(context.Background(), "batch_1")
	if err != nil {
		t.Fatal(err)
	}
	var result batches.Result
	if err := stream.ForEach(context.Background(), results, func(value batches.Result) error {
		result = value
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if result.CustomID != "request_1" || result.Result.Error == nil || result.Result.Error.Message != "bad input" {
		t.Fatalf("batch result = %#v", result)
	}
}
