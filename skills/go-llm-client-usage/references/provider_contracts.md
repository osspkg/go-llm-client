# Provider contracts and authoritative references

Use the snapshot table and endpoint matrix in this file as the local contract
baseline. When an endpoint or field is unclear, inspect the pinned upstream
source and then the current package model; do not infer wire names from another
provider.

## Contract snapshots

These revisions and check dates define the wire-contract baseline used by this
skill:

| Provider | Contract | Revision or check date |
| --- | --- | --- |
| OpenAI | OpenAPI snapshot | `ddface9bd361f5fe37943291d23ee2ca72cbcc2b` |
| OpenAI | API reference | checked 2026-09-20 |
| Ollama | OpenAPI schema | `6383a0fa9cbf97494b847226e189f6e36b401a08` |
| Ollama | API guide | checked 2026-09-20 |
| Anthropic | Stable API `2023-06-01` and documentation snapshot | checked 2026-09-21 |
| llama.cpp | Native server API, `master` snapshot | checked 2026-09-21 |

## Endpoint ownership

| Domain | HTTP surface | Streaming | Client |
| --- | --- | --- | --- |
| OpenAI Responses | `/v1/responses` | SSE | `openai.Responses` |
| OpenAI Chat | `/v1/chat/completions` | SSE | `openai.Chat` |
| OpenAI Completions | `/v1/completions` | SSE | `openai.Completions` |
| OpenAI Embeddings | `/v1/embeddings` | none | `openai.Embeddings` |
| OpenAI Models | `/v1/models` | none | `openai.Models` |
| OpenAI Files | `/v1/files` | bounded binary download | `openai.Files` |
| OpenAI Batches | `/v1/batches` | none | `openai.Batches` |
| OpenAI Audio | `/v1/audio/*` | bounded binary speech | `openai.Audio` |
| OpenAI Images | `/v1/images/*` | endpoint-dependent | `openai.Images` |
| OpenAI Moderations | `/v1/moderations` | none | `openai.Moderations` |
| OpenAI stateful resources | assistants, threads, runs, vector stores, containers, evals | endpoint-dependent | `openai.Stateful` |
| OpenAI organization | `/v1/organization/*` | none | `openai.Organization` |
| OpenAI Realtime | `/v1/realtime` | WebSocket | `openai.Realtime` |
| Ollama generation | `/api/generate`, `/api/chat` | NDJSON | `ollama.Generate`, `ollama.Chat` |
| Ollama embeddings | `/api/embed`, `/api/embeddings` | none | `ollama.Embeddings` |
| Ollama models | `/api/tags`, `/api/ps`, `/api/show`, `/api/create`, `/api/pull`, `/api/push`, `/api/copy`, `/api/delete`, `/api/version` | NDJSON lifecycle | `ollama.Models` |
| Anthropic Messages | `/v1/messages`, `/v1/messages/count_tokens` | SSE for messages | `anthropic.Messages` |
| Anthropic Models | `/v1/models`, `/v1/models/{id}` | none | `anthropic.Models` |
| Anthropic Files | `/v1/files/*` | bounded binary download | `anthropic.Files` |
| Anthropic Batches | `/v1/messages/batches/*` | JSONL results | `anthropic.Batches` |
| llama.cpp completion | `/completion` | SSE | `llama.Completions` |
| llama.cpp embeddings | `/embeddings`, `/embedding` | none | `llama.Embeddings` |
| llama.cpp tokenization | `/tokenize`, `/detokenize` | none | `llama.Tokenization` |
| llama.cpp templates and server | `/apply-template`, `/health`, `/props` | none | `llama.Templates`, `llama.Server` |
| llama.cpp operations | `/slots`, `/lora-adapters`, `/metrics`, `/models`, `/models/sse`, `/rerank` | model lifecycle SSE | `llama.Slots`, `llama.LoRA`, `llama.Metrics`, `llama.Models`, `llama.Rerank` |

OpenAI-compatible routes served by llama.cpp remain owned by `openai` and are
not duplicated under `llama`.

## Authoritative upstream sources

These links are the authoritative references for wire contracts:

- OpenAI OpenAPI snapshot: [openai-openapi](https://github.com/openai/openai-openapi/tree/ddface9bd361f5fe37943291d23ee2ca72cbcc2b)
- OpenAI API reference: [API reference](https://platform.openai.com/docs/api-reference/introduction)
- Ollama OpenAPI schema: [openapi.yaml](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/openapi.yaml)
- Ollama API guide: [API docs](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/api.md)
- Anthropic API overview: [API overview](https://docs.anthropic.com/en/api/overview)
- Anthropic Messages: [Messages API](https://docs.anthropic.com/en/api/messages)
- llama.cpp native server: [server README](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md)
- llama.app API: [llama.app API](https://llama.app/docs/api)

## Compatibility rules

- OpenAI-compatible services may use `openai.WithBaseURL`; their capability
  support must be explicit and no hidden fallback is allowed.
- Native llama.cpp endpoints are exposed only under `llama`; its OpenAI-shaped
  `/v1/*` endpoints remain in `openai`.
- Anthropic managed-agent beta domains are outside the stable client scope.
- Ollama streams are NDJSON; OpenAI, Anthropic, and native llama.cpp streams
  use SSE where the endpoint advertises streaming.
- New upstream endpoints require an intentional revision and capability-matrix
  update before adding public methods.
