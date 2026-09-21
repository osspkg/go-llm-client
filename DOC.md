# API and architecture reference

`go.osspkg.com/llm-client` contains independent provider-native clients. Each
client performs typed HTTP or WebSocket calls to an LLM provider; the Go code
does not execute a model or invoke a Go callback inside the model.

## Providers and package boundaries

- `openai` covers OpenAI-compatible REST, the pinned OpenAI public/admin REST
  domains, media, uploads, and Realtime WebSocket.
- `anthropic` covers the stable Anthropic Messages, Count Tokens, Models,
  Files, and Message Batches APIs. Agents and other managed-agent beta domains
  are intentionally outside this release.
- `llama` covers native llama.cpp server endpoints. Its OpenAI-compatible
  `/v1/*` endpoints remain available through `openai.WithBaseURL` and are not
  duplicated here.
- `ollama` covers the official Ollama HTTP API.
- `pkg/*` contains provider-independent transport, authentication, typed
  errors, pagination, stream parsers, codec, and WebSocket lifecycle helpers.

Domain packages own their provider-native request, response, and event models:
`anthropic/messages`, `anthropic/models`, `anthropic/files`,
`anthropic/batches`, and the corresponding `llama/*` packages. This preserves
wire semantics and makes unsupported provider capabilities explicit.

## Construction and authentication

```go
anthropicClient, err := anthropic.New(
	 anthropic.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
)
if err != nil {
	return err
}

llamaClient, err := llama.New(
	llama.WithBaseURL("http://localhost:8080"),
)
```

Anthropic defaults to `https://api.anthropic.com/v1` and sends
`anthropic-version: 2023-06-01`. `WithAPIKey` sends `x-api-key`; use
`WithBearerToken` for `Authorization: Bearer`. Workspace and beta headers are
explicit options. `WithAuthProvider` can return custom headers per operation.

Native llama.cpp defaults to `http://localhost:8080`; its auth callback is also
called for every request. `RequestMeta.Domain` always contains the destination
hostname, including custom base URLs. Credentials are never put in URLs or
error messages.

## Anthropic Messages

```go
content, err := messages.TextContent("Explain SSE in one sentence.")
if err != nil {
	return err
}
response, err := anthropicClient.Messages().Create(ctx, messages.Request{
	Model:     "claude-3-5-sonnet-latest",
	MaxTokens: 256,
	Messages: []messages.Message{{
		Role:    "user",
		Content: content,
	}},
})
```

Messages expose typed structures for text, image, document, thinking,
redacted-thinking, tool-use, and tool-result blocks. `Message.Content` and
other fields that are genuine Anthropic string-or-array or schema-defined JSON
unions remain represented by explicit provider union wrappers; schema-defined
JSON such as tool arguments continues to use `json.RawMessage`.

`CreateStream` returns `stream.Iterator[messages.StreamEvent]`. The shared SSE
parser recognizes Anthropic event frames, `[DONE]`, blank-line boundaries,
bounded event size, malformed JSON, cancellation, and explicit `Close`.
`Batches.Results` returns a bounded NDJSON iterator for completed batch results.
File content downloads return an `io.ReadCloser` owned by the caller; uploads
consume an `io.Reader` through a bounded multipart request.

## Native llama.cpp

The native client is immutable after construction and exposes domain clients
through accessors:

- `Completions()` for `/completion`, including string, token-array, mixed, and
  multimodal prompts, sampling settings, cache/slot/LoRA overrides, typed
  timings, stop metadata, and SSE streaming;
- `Embeddings()` for `/embeddings` and the singular `/embedding` alias;
- `Tokenization()` for `/tokenize` and `/detokenize`, including token-piece
  string-or-byte unions;
- `Templates()` for `/apply-template`;
- `Server()` for `/health`, `/props`, and `/props` updates;
- `Slots()` for listing and save/restore/erase cache operations;
- `LoRA()` for listing and updating adapter scales;
- `Metrics()` for bounded Prometheus text;
- `Models()` for router listing, download, load, unload, and model lifecycle SSE;
- `Rerank()` for the native reranking endpoint.

Native endpoints can be disabled or absent in a particular llama.cpp build.
The client returns `errors.CapabilityError` for a missing native route; its
cause remains available through `errors.As` as the bounded `HTTPError`. It does
not silently fall back to an OpenAI or Anthropic endpoint.

## Streams and ownership

All typed HTTP streams implement:

```go
type Iterator[T any] interface {
	Next(context.Context) bool
	Value() T
	Err() error
	Close() error
}
```

The iterator owns its response body. Call `Close` on every path, including
early termination. `stream.ForEach` is available for callback-style consumers.
The parser never uses an unbounded scanner: response bodies and event lines are
bounded by transport and stream options.

## Serialization and errors

Provider model packages contain `//go:generate easyjson -all types.go` and the
generated `*_easyjson.go` files are committed. Run `go generate ./...`; never
edit generated files manually. Unknown JSON fields are ignored. `RawMessage`
is used only for upstream-defined arbitrary JSON such as schemas, tool input,
and documented union payloads.

The shared transport has dial, TLS handshake, response-header, request, stream,
and body limits. Automatic retries are restricted to safe idempotent methods;
generation, uploads, streaming, and WebSocket calls are not retried. Inspect
failures with `errors.Is` and `errors.As` for validation, cancellation,
protocol, decode, and `HTTPError` categories.

## Contract revisions and development

The endpoint ownership and upstream snapshots are recorded in
[`docs/capability-matrix.md`](docs/capability-matrix.md). The current
llama.cpp server contract is the documented `master` snapshot checked on the
date recorded there; llama.cpp can add or remove optional native endpoints
between releases.

```text
go generate ./...
make lint
make tests
go test -race ./...
go vet ./...
go mod verify
git diff --check
```

Tests use deterministic `httptest` and local protocol fixtures. No live
provider credentials are required.
