package ollama_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/ollama"
	"go.osspkg.com/llm-client/ollama/chat"
	"go.osspkg.com/llm-client/ollama/generate"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestChatAndGenerateStreamsUseNDJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch request.URL.Path {
		case "/api/chat":
			var input chat.Request
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
				t.Errorf("decode chat request: %v", err)
			}
			if input.Model != "llama-test" {
				t.Errorf("model = %q", input.Model)
			}
			return response(http.StatusOK, `{"model":"llama-test","message":{"role":"assistant","content":"hello"},"done":true}`), nil
		case "/api/generate":
			return response(http.StatusOK, `{"model":"llama-test","response":"one","done":false}`+"\n"+`{"model":"llama-test","response":"two","done":true}`), nil
		default:
			return response(http.StatusNotFound, "not found"), nil
		}
	})}

	client, err := ollama.New(ollama.WithBaseURL("http://example.test"), ollama.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Chat().Create(context.Background(), chat.Request{Model: "llama-test"})
	if err != nil {
		t.Fatal(err)
	}
	if response.Message.Content != "hello" || !response.Done {
		t.Fatalf("chat response = %#v", response)
	}

	iterator, err := client.Generate().CreateStream(context.Background(), generate.Request{Model: "llama-test"})
	if err != nil {
		t.Fatal(err)
	}
	var values []string
	for iterator.Next(context.Background()) {
		values = append(values, iterator.Value().Response)
	}
	if err := iterator.Err(); err != nil {
		t.Fatal(err)
	}
	if len(values) != 2 || values[0] != "one" || values[1] != "two" {
		t.Fatalf("stream values = %v", values)
	}
}

func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
