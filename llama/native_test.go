package llama_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/llama"
	"go.osspkg.com/llm-client/llama/completions"
	"go.osspkg.com/llm-client/llama/embeddings"
	"go.osspkg.com/llm-client/llama/models"
	"go.osspkg.com/llm-client/llama/rerank"
	"go.osspkg.com/llm-client/llama/server"
	"go.osspkg.com/llm-client/llama/templates"
	"go.osspkg.com/llm-client/llama/tokenization"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
)

func TestNativeDomainRoutesAndWireShapes(t *testing.T) { //nolint:gocyclo,revive // one deterministic contract test covers the native domain matrix.
	client, err := llama.New(llama.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body string
		switch r.URL.Path {
		case "/embeddings":
			body = `[{"index":0,"embedding":[[0.1,0.2]]}]`
		case "/tokenize":
			body = `{"tokens":[{"id":1,"piece":"hi"},{"id":2,"piece":[195,161]}]}`
		case "/detokenize":
			body = `{"content":"hi"}`
		case "/props":
			if r.Method == http.MethodPost {
				body = `{"model_path":"updated.gguf"}`
			} else {
				body = `{"model_path":"model.gguf","total_slots":2,"modalities":{"vision":true},"default_generation_settings":{"temperature":0.7}}`
			}
		case "/slots":
			body = `[{"id_slot":1,"n_ctx":4096}]`
		case "/slots/1":
			body = `{"id_slot":1,"filename":"slot.bin"}`
		case "/lora-adapters":
			body = `[{"id":1,"path":"adapter.gguf","scale":0.5}]`
		case "/models":
			if r.Method == http.MethodGet {
				body = `{"data":[{"id":"model","status":{"value":"loaded"},"architecture":{"input_modalities":["text"]}}]}`
			} else {
				body = `{"success":true}`
			}
		case "/models/load", "/models/unload":
			body = `{"success":true}`
		case "/apply-template":
			body = `{"prompt":"formatted"}`
		case "/health":
			body = `{"status":"ok"}`
		case "/metrics":
			body = "llama_prompt_tokens_total 3\n"
		case "/rerank":
			body = `{"results":[{"index":0,"relevance_score":0.9}]}`
		default:
			t.Errorf("unexpected native path %s", r.URL.Path)
			body = `{}`
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}

	embedding, err := client.Embeddings.Create(context.Background(), embeddings.Request{Content: []byte(`"hello"`)})
	if err != nil || len(embedding) != 1 || len(embedding[0].Embedding) != 1 {
		t.Fatalf("embedding = %#v, err=%v", embedding, err)
	}
	tokens, err := client.Tokenization.Tokenize(context.Background(), tokenization.TokenizeRequest{Content: "hi", WithPieces: true})
	if err != nil || len(tokens.Tokens) != 2 || string(tokens.Tokens[1].Piece.Bytes) != string([]byte{195, 161}) {
		t.Fatalf("tokens = %#v, err=%v", tokens, err)
	}
	props, err := client.Server.Props(context.Background())
	if err != nil || props.ModelPath != "model.gguf" || !props.Modalities.Vision {
		t.Fatalf("props = %#v, err=%v", props, err)
	}
	listedSlots, err := client.Slots.List(context.Background())
	if err != nil || len(listedSlots) != 1 || listedSlots[0].ID != 1 {
		t.Fatalf("slots = %#v, err=%v", listedSlots, err)
	}
	adapters, err := client.Lora.List(context.Background())
	if err != nil || len(adapters) != 1 || adapters[0].Path != "adapter.gguf" {
		t.Fatalf("adapters = %#v, err=%v", adapters, err)
	}
	modelsPage, err := client.Models.List(context.Background())
	if err != nil || len(modelsPage.Data) != 1 {
		t.Fatalf("models = %#v, err=%v", modelsPage, err)
	}
	loaded, err := client.Models.Load(context.Background(), models.ModelRequest{Model: "model"})
	if err != nil || !loaded.Success {
		t.Fatalf("load = %#v, err=%v", loaded, err)
	}
	if _, err := client.Server.Health(context.Background()); err != nil {
		t.Fatal(err)
	}
	updated, err := client.Server.UpdateProps(context.Background(), server.Props{ModelPath: "updated.gguf"})
	if err != nil || updated.ModelPath != "updated.gguf" {
		t.Fatalf("updated props = %#v, err=%v", updated, err)
	}
	formatted, err := client.Templates.Apply(context.Background(), templates.Request{
		Messages: []templates.Message{{Role: "user", Content: "hello"}},
	})
	if err != nil || formatted.Prompt != "formatted" {
		t.Fatalf("template = %#v, err=%v", formatted, err)
	}
	if _, err := client.Slots.Save(context.Background(), 1, "slot.bin"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Slots.Restore(context.Background(), 1, "slot.bin"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Slots.Erase(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	metricsText, err := client.Metrics.Get(context.Background())
	if err != nil || metricsText != "llama_prompt_tokens_total 3\n" {
		t.Fatalf("metrics = %q, err=%v", metricsText, err)
	}
	ranked, err := client.Rerank.Create(context.Background(), rerank.Request{Query: "q", Documents: []string{"d"}})
	if err != nil || len(ranked.Results) != 1 {
		t.Fatalf("rerank = %#v, err=%v", ranked, err)
	}
}

func TestNativeCompletionSSE(t *testing.T) {
	client, err := llama.New(llama.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/completion" {
			t.Errorf("path = %q", r.URL.Path)
		}
		var request completions.Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if string(request.Prompt) != `"hello"` || !request.Stream {
			t.Errorf("request = %#v", request)
		}
		body := "data: {\"content\":\"hi\"}\n\n" + "data: {\"content\":\"!\",\"stop\":true}\n\n"
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"text/event-stream"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	iterator, err := client.Completions.CreateStream(context.Background(), completions.Request{Prompt: completions.StringPrompt("hello")})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = iterator.Close() }()
	if !iterator.Next(context.Background()) || iterator.Value().Content != "hi" {
		t.Fatalf("first event = %#v err=%v", iterator.Value(), iterator.Err())
	}
	if !iterator.Next(context.Background()) || !iterator.Value().Stop {
		t.Fatalf("second event = %#v err=%v", iterator.Value(), iterator.Err())
	}
	if iterator.Next(context.Background()) || iterator.Err() != nil {
		t.Fatalf("terminal event: err=%v", iterator.Err())
	}
}

func TestMissingNativeEndpointReturnsCapabilityError(t *testing.T) {
	client, err := llama.New(llama.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"not found"}}`)),
		}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Templates.Apply(context.Background(), templates.Request{})
	var capabilityErr *llmerrors.CapabilityError
	if !errors.As(err, &capabilityErr) || !errors.Is(err, llmerrors.ErrProtocol) {
		t.Fatalf("error = %v", err)
	}
	var httpErr *llmerrors.HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusNotFound {
		t.Fatalf("cause = %v", err)
	}
}
