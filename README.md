# go-llm-client

[![Go Reference](https://pkg.go.dev/badge/go.osspkg.com/llm-client.svg)](https://pkg.go.dev/go.osspkg.com/llm-client)
[![Go Report Card](https://goreportcard.com/badge/go.osspkg.com/llm-client)](https://goreportcard.com/report/go.osspkg.com/llm-client)
![Go version](https://img.shields.io/badge/Go-1.26.8-00ADD8?logo=go&logoColor=white)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)](LICENSE)

Typed Go clients for OpenAI-compatible providers and Ollama. The library calls
provider HTTP and WebSocket APIs directly and exposes provider-native models
organized by bounded context.

## Contents

- [Features](#features)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Streaming](#streaming)
- [Provider coverage](#provider-coverage)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Features

- OpenAI-compatible REST client with custom base URLs and per-request auth.
- OpenAI domain packages for responses, chat, completions, embeddings, models,
  files, uploads, batches, media, stateful resources, organization APIs, and
  Realtime WebSocket sessions.
- Ollama client with generation, chat, embeddings, model management, blobs, and
  version APIs.
- Typed SSE and NDJSON iterators with explicit `Next`, `Value`, `Err`, and
  `Close` lifecycle methods.
- Bounded response and event sizes, context cancellation, typed protocol errors,
  and deterministic `httptest`-based tests.
- `easyjson`-generated provider models; `github.com/coder/websocket` is used
  only by the OpenAI Realtime package.

## Installation

```bash
go get go.osspkg.com/llm-client
```

The module requires Go 1.26.8 or newer. See the
[capability matrix](docs/capability-matrix.md) for endpoint-level coverage and
the [API reference](DOC.md) for package details.

## Quick start

The following example creates an OpenAI-compatible client and sends a chat
request. Set `OPENAI_API_KEY` before running it.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.osspkg.com/llm-client/openai"
	"go.osspkg.com/llm-client/openai/chat"
	"go.osspkg.com/llm-client/pkg/auth"
)

func main() {
	client, err := openai.New(
		openai.WithAuthProvider(auth.StaticBearer(os.Getenv("OPENAI_API_KEY"))),
	)
	if err != nil {
		panic(err)
	}

	response, err := client.Chat.Create(context.Background(), chat.Request{
		Model: "model",
		Messages: []chat.Message{{
			Role:    "user",
			Content: json.RawMessage([]byte(`"Hello"`)),
		}},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(response.Choices[0].Message.Content))
}
```

For an OpenAI-compatible service, provide its endpoint explicitly:

```go
client, err := openai.New(
	openai.WithBaseURL("http://localhost:8080/v1"),
	openai.WithAuthProvider(auth.StaticBearer("token")),
)
```

Ollama uses `http://localhost:11434` by default:

```go
package main

import (
	"context"
	"fmt"

	"go.osspkg.com/llm-client/ollama"
	"go.osspkg.com/llm-client/ollama/chat"
)

func main() {
	client, err := ollama.New()
	if err != nil {
		panic(err)
	}

	response, err := client.Chat.Create(context.Background(), chat.Request{
		Model:    "llama3.2",
		Messages: []chat.Message{{Role: "user", Content: "Hello"}},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Message.Content)
}
```

## Streaming

Streaming calls return typed iterators. The caller owns the iterator and must
close it when processing is complete.

```go
events, err := client.Chat.CreateStream(ctx, chat.Request{
	Model: "model",
	Messages: []chat.Message{{
		Role:    "user",
		Content: json.RawMessage([]byte(`"Tell me a short joke."`)),
	}},
})
if err != nil {
	return err
}
defer events.Close()

for events.Next(ctx) {
	chunk := events.Value()
	// Process the typed chunk.
	_ = chunk
}
if err := events.Err(); err != nil {
	return err
}
```

The same iterator contract is used for OpenAI SSE and Ollama NDJSON streams.
See [`pkg/stream`](pkg/stream) and the streaming sections in [`DOC.md`](DOC.md)
for parser limits and error handling.

## Provider coverage

| Provider | Package | Transport | Default endpoint |
| --- | --- | --- | --- |
| OpenAI-compatible | [`openai`](openai) | HTTP; Realtime WebSocket | `https://api.openai.com/v1` |
| Ollama | [`ollama`](ollama) | HTTP; NDJSON streaming | `http://localhost:11434` |

Provider APIs remain separate. There is no provider-neutral facade that hides
provider-specific capabilities or request models.

For the current endpoint and capability status, see:

- [`docs/capability-matrix.md`](docs/capability-matrix.md)
- [`DOC.md`](DOC.md) — English API and architecture reference
- [`DOC.ru.md`](DOC.ru.md) — Russian guide
- [`PLAN.md`](PLAN.md) — implementation checkpoints and compatibility notes

## Development

Generated `easyjson` files are committed to the repository and must not be
edited manually. Run the complete local quality workflow:

```bash
go generate ./...
make lint
make tests
go test -race ./...
go vet ./...
go mod verify
git diff --check
```

Tests use deterministic mocks and do not require provider credentials or live
network access. See [`AGENTS.md`](AGENTS.md) for architecture, dependency,
serialization, security, and resource-lifecycle rules.

## Contributing

Contributions are welcome. Before opening a pull request:

1. Keep provider-specific types and operations in their provider/domain package.
2. Keep shared transport, auth, error, pagination, stream, and WebSocket
   lifecycle code in `pkg` without provider imports.
3. Regenerate easyjson output with `go generate ./...`.
4. Run the development checks listed above.
5. Add deterministic tests for new protocol behavior and resource ownership.

## License

This project is available under the [BSD 3-Clause License](LICENSE).
