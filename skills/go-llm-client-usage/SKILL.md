---
name: go-llm-client-usage
description: Use go.osspkg.com/llm-client from Go code by selecting the correct provider-native client, building typed requests, configuring auth, consuming bounded streams, and handling typed errors. Applies to this repository's OpenAI-compatible, Anthropic-compatible, native llama.cpp, and Ollama integrations.
---

# Go LLM Client Usage

Use the current repository source and documentation as the source of truth. Do
not invent a provider-neutral API: this library intentionally keeps provider
models and bounded contexts separate.

## Select the provider

- Use `openai` for OpenAI REST, OpenAI-compatible `/v1` servers, and Realtime.
  This also covers llama.cpp's OpenAI-compatible routes when they are enabled.
- Use `anthropic` for the stable Anthropic Messages API and its Models, Files,
  and Message Batches domains.
- Use `llama` only for native llama.cpp routes such as `/completion`,
  `/tokenize`, `/slots`, and `/props`.
- Use `ollama` for Ollama `/api/*` operations and NDJSON generation/chat.

Root clients are immutable after construction. Access domain clients through
methods such as `client.Chat()`, `client.Responses()`, `client.Messages()`,
`client.Completions()`, and `client.Generate()`; do not expect exported mutable
domain fields.

## Requests and authentication

Prefer provider options for standard credentials (`WithAPIKey`,
`WithBearerToken`, or `auth.StaticBearer`) and `WithBaseURL` for compatible
servers. For dynamic credentials use `WithAuthProvider`. The callback runs for
each HTTP request and WebSocket handshake, and its `auth.RequestMeta.Domain`
always contains the destination hostname, including custom base URLs. Never
put credentials in query parameters, URLs, logs, or returned errors.

Use typed provider request models. Union constructors that can fail return an
error, for example `messages.TextContent` and
`completions.StringPrompt`; handle that error before sending the request.

## Streams and resources

Streaming methods return `pkg/stream.Iterator[T]` with `Next(ctx)`, `Value()`,
`Err()`, and `Close()`. Always close an iterator, including early-return paths;
use `stream.ForEach` when callback consumption is clearer. SSE is used by
OpenAI, Anthropic, and native llama.cpp streams; Ollama uses NDJSON. Pass
contexts through every blocking operation and inspect `Err()` after iteration.

Responses and event lines are bounded. Do not replace the shared stream parser
with an unbounded `bufio.Scanner`, and do not add automatic retries to
generation, upload, streaming, or WebSocket operations.

## Errors and tests

Use `errors.Is` and `errors.As` with `pkg/errors` to distinguish validation,
cancellation, body-limit, decode, protocol, HTTP, and capability errors. A
missing native llama.cpp endpoint is represented by `CapabilityError`; it does
not silently fall back to OpenAI routes.

For library changes, use deterministic `httptest` or local WebSocket fixtures;
never require provider credentials or live APIs. Run `go generate ./...`,
`make lint`, `make tests`, and the relevant race/vet checks. Generated
easyjson files are committed and must be regenerated rather than edited.

## Supporting references

Read only the reference needed for the task:

- [api_reference.md](references/api_reference.md) — package map, root accessors,
  options, stream and error contracts.
- [provider_contracts.md](references/provider_contracts.md) — local contract
  documents, pinned revisions, and authoritative upstream links.
- [usage_examples.md](references/usage_examples.md) — provider, auth, stream,
  error, upload, and Realtime examples.
