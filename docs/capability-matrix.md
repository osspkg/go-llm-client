# Provider capability matrix

The matrix is intentionally versioned. Endpoint behavior is checked against the
following upstream revisions:

| Provider | Contract | Revision | Reference |
| --- | --- | --- | --- |
| OpenAI | OpenAPI snapshot | `ddface9bd361f5fe37943291d23ee2ca72cbcc2b` | [openai-openapi](https://github.com/openai/openai-openapi/tree/ddface9bd361f5fe37943291d23ee2ca72cbcc2b) |
| OpenAI | API reference | checked 2026-09-20 | [API reference](https://platform.openai.com/docs/api-reference/introduction) |
| Ollama | OpenAPI schema | `6383a0fa9cbf97494b847226e189f6e36b401a08` | [openapi.yaml](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/openapi.yaml) |
| Ollama | API guide | checked 2026-09-20 | [API docs](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/api.md) |
| Anthropic | Stable API version `2023-06-01`, docs snapshot | checked 2026-09-21 | [API overview](https://docs.anthropic.com/en/api/overview) |
| llama.cpp | Native server API, `master` snapshot | checked 2026-09-21 | [server API](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md) |

| Domain | HTTP surface | Streaming | Auth metadata | Client |
| --- | --- | --- | --- | --- |
| OpenAI Responses | `/v1/responses` | SSE | per request | `openai.Responses` |
| OpenAI Chat | `/v1/chat/completions` | SSE | per request | `openai.Chat` |
| OpenAI Completions | `/v1/completions` | SSE | per request | `openai.Completions` |
| OpenAI Embeddings | `/v1/embeddings` | none | per request | `openai.Embeddings` |
| OpenAI Models | `/v1/models` | none | per request | `openai.Models` |
| OpenAI Files | `/v1/files` | binary download | per request | `openai.Files` |
| OpenAI Batches | `/v1/batches` | none | per request | `openai.Batches` |
| OpenAI Audio | `/v1/audio/*` | binary speech | per request | `openai.Audio` |
| OpenAI Images | `/v1/images/*` | SSE for generation and edit | per request | `openai.Images` |
| OpenAI Moderations | `/v1/moderations` | none | per request | `openai.Moderations` |
| OpenAI stateful resources | assistants, threads, runs, vector stores, containers, evals | endpoint-dependent | per request | `openai.Stateful` |
| OpenAI organization | `/v1/organization/*` | none | per request | `openai.Organization` |
| OpenAI Realtime | `/v1/realtime` | WebSocket | handshake callback | `openai.Realtime` |
| Ollama generation | `/api/generate`, `/api/chat` | NDJSON | per request | `ollama.Generate`, `ollama.Chat` |
| Ollama embeddings | `/api/embed`, `/api/embeddings` | none | per request | `ollama.Embeddings` |
| Ollama models | `/api/tags`, `/api/ps`, `/api/show`, `/api/create`, `/api/pull`, `/api/push`, `/api/copy`, `/api/delete`, `/api/version` | NDJSON lifecycle | per request | `ollama.Models` |
| Anthropic Messages | `/v1/messages`, `/v1/messages/count_tokens` | SSE for messages | per request | `anthropic.Messages` |
| Anthropic Models | `/v1/models`, `/v1/models/{id}` | none | per request | `anthropic.Models` |
| Anthropic Files | `/v1/files/*` | bounded binary download | per request | `anthropic.Files` |
| Anthropic Batches | `/v1/messages/batches/*` | JSONL results | per request | `anthropic.Batches` |
| llama.cpp completion | `/completion` | SSE | per request | `llama.Completions` |
| llama.cpp embeddings | `/embeddings`, `/embedding` | none | per request | `llama.Embeddings` |
| llama.cpp tokenization | `/tokenize`, `/detokenize` | none | per request | `llama.Tokenization` |
| llama.cpp templates/server | `/apply-template`, `/health`, `/props` | none | per request | `llama.Templates`, `llama.Server` |
| llama.cpp operations | `/slots`, `/lora-adapters`, `/metrics`, `/models`, `/models/sse`, `/rerank` | model lifecycle SSE | per request | `llama.Slots`, `llama.Lora`, `llama.Metrics`, `llama.Models`, `llama.Rerank` |

The root clients expose provider-native types. A capability matrix is a local
configuration snapshot; it does not silently emulate or fall back to another
endpoint when a provider does not support an operation.

The implementation covers the listed Anthropic and native llama.cpp endpoints.
Anthropic beta managed-agent domains and llama.cpp internal `/tools` endpoints
are explicitly out of scope. Endpoints added after the recorded snapshots
require a deliberate revision update and matrix review before being exposed.
