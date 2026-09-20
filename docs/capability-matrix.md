# Provider capability matrix

The matrix is intentionally versioned. Endpoint behavior is checked against the
following upstream revisions:

| Provider | Contract | Revision | Reference |
| --- | --- | --- | --- |
| OpenAI | OpenAPI snapshot | `ddface9bd361f5fe37943291d23ee2ca72cbcc2b` | [openai-openapi](https://github.com/openai/openai-openapi/tree/ddface9bd361f5fe37943291d23ee2ca72cbcc2b) |
| OpenAI | API reference | checked 2026-09-20 | [API reference](https://platform.openai.com/docs/api-reference/introduction) |
| Ollama | OpenAPI schema | `6383a0fa9cbf97494b847226e189f6e36b401a08` | [openapi.yaml](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/openapi.yaml) |
| Ollama | API guide | checked 2026-09-20 | [API docs](https://github.com/ollama/ollama/blob/6383a0fa9cbf97494b847226e189f6e36b401a08/docs/api.md) |

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
| OpenAI Images | `/v1/images/*` | none | per request | `openai.Images` |
| OpenAI Moderations | `/v1/moderations` | none | per request | `openai.Moderations` |
| OpenAI stateful resources | assistants, threads, runs, vector stores, containers, evals | endpoint-dependent | per request | `openai.Stateful` |
| OpenAI organization | `/v1/organization/*` | none | per request | `openai.Organization` |
| OpenAI Realtime | `/v1/realtime` | WebSocket | handshake callback | `openai.Realtime` |
| Ollama generation | `/api/generate`, `/api/chat` | NDJSON | per request | `ollama.Generate`, `ollama.Chat` |
| Ollama embeddings | `/api/embed`, `/api/embeddings` | none | per request | `ollama.Embeddings` |
| Ollama models | `/api/tags`, `/api/ps`, `/api/show`, `/api/create`, `/api/pull`, `/api/push`, `/api/copy`, `/api/delete`, `/api/version` | NDJSON lifecycle | per request | `ollama.Models` |

The root clients expose provider-native types. A capability matrix is a local
configuration snapshot; it does not silently emulate or fall back to another
endpoint when a provider does not support an operation.

The current implementation tracks the listed core endpoints and records the
remaining full-snapshot audit in `PLAN.md`. Endpoints added after the pinned
revisions require a deliberate revision update and a matrix review before being
exposed.
