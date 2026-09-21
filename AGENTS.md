# Project instructions

## Architecture

- This repository is a Go library (`go.osspkg.com/llm-client`) targeting Go 1.26.
- `openai`, `anthropic`, `llama`, and `ollama` are independent provider bounded contexts.
- Provider domain packages own provider-native request, response, event, and error models.
- Provider-independent infrastructure belongs under `pkg` and must not import provider packages.
- Native llama.cpp `/v1/*` compatibility routes are owned by `openai`; `llama` contains only native routes.
- Anthropic managed-agent beta domains are not part of the stable client surface.
- Root provider clients compose domain clients but do not hide provider-specific wire semantics.

## Dependencies and serialization

- Runtime dependencies are the standard library, `github.com/mailru/easyjson`, and
  `github.com/coder/websocket` only for OpenAI Realtime.
- Do not add test frameworks, goleak, OpenAPI generators, or convenience HTTP libraries.
- Models live in importable packages, never `package main`.
- Generated `*_easyjson.go` files are committed and must be regenerated with `go generate ./...`.
- Do not hand-edit generated files. Use typed unions instead of `interface{}` where possible.
- Unknown JSON fields are ignored. `json.RawMessage` is allowed only where the upstream
  contract explicitly permits arbitrary JSON (for example a JSON schema or tool arguments).

## API and resource rules

- Context is the first argument for every blocking operation.
- Constructors validate options and return errors instead of panicking.
- Exported APIs use documented typed errors and preserve causes with `%w`.
- HTTP response bodies and streams have explicit ownership and are always closed.
- No credentials may appear in URLs, logs, or returned diagnostic strings.
- Response bodies, stream events, multipart inputs, and retry counts are bounded.
- The default request-body limit is 64 MiB and the default buffered response limit is 32 MiB;
  callers can lower or raise them explicitly with transport/client options.
- No unowned background goroutines. Realtime sessions have explicit `Close` lifecycle.
- Clients are immutable after construction and safe for concurrent use unless documented otherwise.

## Development and validation

Run, in order when applicable:

```text
go generate ./...
make lint
make tests
go test -race ./...
go vet ./...
go mod verify
git diff --check
```

`make lint` is the repository lint gate. `make tests` and the race detector are required
for behavioral validation. Never claim a live provider integration test passed unless it
was actually run with explicit credentials and network access.

## Documentation

Public Go APIs and executable examples are written in English. `DOC.ru.md` contains the
Russian usage and architecture guide. Update README, docs, examples, capability matrix,
and tests when public behavior changes.
