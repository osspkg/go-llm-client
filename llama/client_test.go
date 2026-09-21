package llama_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/llama"
	"go.osspkg.com/llm-client/llama/completions"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestNativeCompletionUsesNativeRoute(t *testing.T) {
	client, err := llama.New(llama.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/completion" {
			t.Errorf("path = %q", r.URL.Path)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"content":"hello","stop":true}`))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Completions.Create(context.Background(), completions.Request{Prompt: []byte(`"hello"`)})
	if err != nil {
		t.Fatal(err)
	}
	if response.Content != "hello" || !response.Stop {
		t.Fatalf("response = %#v", response)
	}
}
