# Local API reference

This is a routing reference, not a replacement for Go doc. Verify signatures
against the current package source before proposing code.

## Root clients and domain accessors

All provider roots are constructed with `New(options ...Option)` and expose
immutable domain clients through methods:

| Root | Domain accessors |
| --- | --- |
| `openai` | `Responses()`, `Chat()`, `Completions()`, `Assistants()`, `Threads()`, `Runs()`, `VectorStores()`, `Containers()`, `Evals()`, `Embeddings()`, `Models()`, `Files()`, `Uploads()`, `Batches()`, `FineTuning()`, `Audio()`, `Images()`, `Moderations()`, `Stateful()`, `Organization()`, `Realtime()` |
| `anthropic` | `Messages()`, `Models()`, `Files()`, `Batches()` |
| `llama` | `Completions()`, `Embeddings()`, `Tokenization()`, `Templates()`, `Server()`, `LoRA()`, `Metrics()`, `Models()`, `Slots()`, `Rerank()` |
| `ollama` | `Generate()`, `Chat()`, `Embeddings()`, `Models()`, `Blobs()`, `Version()` |

The root accessor returns a provider-specific bounded context. Use the domain
package's request and response types; do not mix models across providers.

## Important constructors and options

- `openai.New`: use `WithBaseURL` for compatible servers,
  `WithAuthProvider` for dynamic headers, and `WithCapability` to explicitly
  disable an operation family.
- `anthropic.New`: use `WithAPIKey`, `WithBearerToken`, `WithVersion`,
  `WithWorkspace`, and `WithBeta` as needed. The default API version is
  `2023-06-01`.
- `llama.New`: native endpoint default is `http://localhost:8080`.
- `ollama.New`: default endpoint is `http://localhost:11434`.
- Every root supports a custom HTTP client and request/body limits through its
  documented options.
- `llama.EnvironmentScheme()` and `ollama.EnvironmentScheme()` return
  client-independent `pkg/environment.Scheme` metadata for provider server
  environment variables.

`auth.HeaderProvider` receives `auth.RequestMeta`. `RequestMeta.Domain` is
always the destination hostname without port, path, query, or credentials.

## Domain selection

| Need | Package and operation |
| --- | --- |
| OpenAI-compatible chat | `openai.Chat().Create` or `CreateStream` |
| OpenAI Responses | `openai.Responses().Create` or `CreateStream` |
| Anthropic messages | `anthropic.Messages().Create` or `CreateStream` |
| Native llama.cpp completion | `llama.Completions().Create` or `CreateStream` |
| Ollama chat/generation | `ollama.Chat().CreateStream`, `ollama.Generate().CreateStream` |
| Typed embeddings | provider-specific `Embeddings().Create` |
| Native tokenization | `llama.Tokenization().Tokenize` / `Detokenize` |
| Native server state | `llama.Server()`, `Slots()`, `Models()`, or `Metrics()` |

Do not use `llama` for `/v1/chat/completions` or `/v1/embeddings`; those
OpenAI-compatible routes belong to `openai`.

## Environment schemes

Environment schemes are metadata functions, not client methods. Each variable
contains its exact name, finite `AllowedValues` when applicable, documented
`Default`, and a description of its purpose and value format. An empty
`AllowedValues` slice means the provider accepts a scalar or provider-dependent
format such as a path, duration, integer, or comma-separated list. Calling a
scheme does not read or modify the process environment.

## Stream and error contracts

```go
type Iterator[T any] interface {
	Next(context.Context) bool
	Value() T
	Err() error
	Close() error
}
```

Call `Close` on every path. After `Next` returns false, inspect `Err`. Use
`stream.ForEach` for callback-style consumption. Use `errors.As` for
`*errors.HTTPError`, `*errors.DecodeError`, `*errors.ValidationError`, and
`*errors.CapabilityError`; use `errors.Is` for shared sentinel categories.

## Reference ownership

This skill keeps its routing, contract, and example guidance in its own
reference files. Use [provider_contracts.md](provider_contracts.md) for the
versioned wire-contract baseline and [usage_examples.md](usage_examples.md)
for complete examples.
