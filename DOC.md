# API and architecture reference

`go.osspkg.com/llm-client` contains two provider-native clients:

```go
client, err := openai.New(
	openai.WithBaseURL("https://api.openai.com/v1"),
	openai.WithAuthProvider(auth.StaticBearer("token")),
)
if err != nil {
	return err
}

completion, err := client.Chat.Create(ctx, chat.Request{Model: "model"})
```

The provider exposes an HTTP or WebSocket API. The Go library performs typed
transport calls to that provider; the Go client does not execute a model or
invoke a Go callback from inside the LLM.

## Package boundaries

`openai/*` and `ollama/*` are bounded contexts containing provider-native
requests, responses, and operations. `pkg/*` contains provider-independent
transport, auth, errors, pagination, stream, codec, and WebSocket lifecycle
components. Provider types never cross into `pkg`.

## Streams and resource ownership

HTTP streams implement `stream.Iterator[T]`:

```go
type Iterator[T any] interface {
	Next(context.Context) bool
	Value() T
	Err() error
	Close() error
}
```

SSE recognizes `data:` frames and `[DONE]`; Ollama NDJSON handles bounded
newline-delimited values. Every iterator owns its response body and callers
must call `Close`, including on early termination. Realtime sessions use
`Connect`, `Send`, `Receive`, and `Close`; writes and reads are serialized by
the common WebSocket lifecycle wrapper.

## Serialization

Importable provider model packages contain `//go:generate easyjson -all
types.go` and committed `*_easyjson.go` output. Unknown JSON fields are
ignored. `json.RawMessage` is used only for upstream fields that intentionally
carry arbitrary JSON, such as schemas and tool arguments.

## Timeouts, limits, retries, and errors

The shared transport has bounded response and request bodies, dial/TLS/header
timeouts, context propagation, and optional retries. Retries are disabled for
generation, uploads, streams, and WebSockets; only safe idempotent methods may
be retried. HTTP errors retain status and request ID without retaining the
provider body or credentials. Use `errors.Is` and `errors.As` to inspect typed
validation, protocol, decode, cancellation, and HTTP errors.

## Contract revisions

See [docs/capability-matrix.md](docs/capability-matrix.md) for pinned upstream
revisions and endpoint ownership.
