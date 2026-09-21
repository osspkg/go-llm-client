# Implementation plan

Status markers: `[ ]` pending, `[x]` complete, `[!]` blocked or needs review.

## Decisions

- OpenAI-compatible HTTP is the primary OpenAI provider contract.
- The full pinned OpenAI public/admin REST surface and Realtime WebSocket are in scope.
- Ollama HTTP follows its pinned OpenAPI schema and includes streaming/model management.
- Provider-native typed APIs are primary; there is no mandatory normalized facade.
- Runtime dependencies are stdlib, easyjson, and coder/websocket for Realtime only.
- Unknown JSON fields are ignored; arbitrary JSON is retained only in spec-defined fields.
- Streams use `Next/Value/Err/Close`; Realtime uses an explicit Session lifecycle.
- Tests are deterministic mocks only; no live provider credentials are required.
- Anthropic stable API uses the documented `2023-06-01` version; managed-agent
  beta domains remain out of scope.
- Native llama.cpp uses its native routes only; OpenAI-compatible `/v1/*` routes
  remain owned by `openai`.
- Missing native llama.cpp routes are surfaced as `pkg/errors.CapabilityError`
  while preserving the underlying bounded `HTTPError` in the error chain.
- SSE line accumulation uses bounded `bufio.Reader.ReadLine` fragments; it never
  buffers an unbounded physical line before applying the event limit.

## Tasks

- [x] 01. Create project rules and this implementation plan. **CP-01:** documentation exists.
- [!] 02. Pin OpenAI/Ollama upstream revisions and capability matrix. **CP-02:** revisions and implemented domain families are documented; a full per-operation audit against the pinned OpenAI snapshot remains required.
- [x] 03. Add and verify approved dependencies. **CP-03:** dependency graph is constrained.
- [x] 04. Implement provider-independent auth, errors, transport, pagination, stream, codec, and WebSocket packages. **CP-04:** no provider imports under `pkg`.
- [x] 05. Implement bounded HTTP transport, auth callback, retry policy, capability gate, and body limits. **CP-05:** transport tests pass.
- [x] 06. Implement typed error/protocol model and cleanup paths. **CP-06:** errors are inspectable and bodies close.
- [x] 07. Add easyjson model generation and committed generated output. **CP-07:** `go generate ./...` is reproducible.
- [x] 08. Implement the OpenAI root client and capability matrix. **CP-08:** custom base URL/auth/capability gating works.
- [!] 09. Implement OpenAI generation domains. **CP-09:** core generation endpoints and streams are typed; the complete snapshot operation audit remains required.
- [x] 10. Implement OpenAI files, uploads, and batches. **CP-10:** multipart/binary/batch lifecycle APIs exist with bounded request bodies.
- [!] 11. Implement OpenAI fine-tuning, assistants, threads, runs, vector stores, containers, and evals. **CP-11:** domain packages and core lifecycle methods exist; full pinned REST operation coverage remains required.
- [!] 12. Implement OpenAI audio and image domains. **CP-12:** JSON and streaming multipart media paths exist; all snapshot variants still need audit.
- [!] 13. Implement OpenAI organization/admin domains. **CP-13:** organization operations and per-operation auth metadata exist; the full admin surface remains required.
- [x] 14. Implement OpenAI Realtime Session over coder/websocket. **CP-14:** explicit session, handshake auth, typed send/receive, cancellation and close lifecycle are covered by a deterministic local WebSocket contract test. The test skips only where the execution sandbox prohibits binding a loopback listener.
- [x] 15. Implement the Ollama root client and pinned core domains. **CP-15:** generation, embeddings, model lifecycle, blobs, and version operations are typed.
- [x] 16. Add stream parser fuzzing and allocation benchmarks. **CP-16:** bounded parser fuzz smoke runs and `b.Loop` benchmarks pass.
- [x] 17. Add deterministic HTTP/WebSocket contract tests. **CP-17:** HTTP, SSE/NDJSON and Realtime WebSocket contract tests require no credentials or external provider network.
- [x] 18. Add README, DOC.md, DOC.ru.md, examples, and capability matrix docs. **CP-18:** example tests compile.
- [x] 19. Add Makefile generation/verification targets and run quality gates. **CP-19:** `make verify` passes generation, lint, tests, race, vet, module verification and diff checks. `govulncheck` is best-effort and cannot fetch its database from this sandbox.
- [!] 20. Perform final API, security, concurrency, and worktree review. **CP-20:** transport and Realtime lifecycle have been reviewed; a full typed-operation audit against the current pinned OpenAI snapshot remains open.
- [x] 21. Add the Anthropic stable provider root and bounded contexts. **CP-21:** Messages, Count Tokens, Models, Files, Batches, auth headers, cursor pages, multipart upload/download, SSE events, and JSONL batch results are typed and covered by deterministic tests.
- [x] 22. Add the native llama.cpp provider root and bounded contexts. **CP-22:** completion, embeddings, tokenization, templates, server properties, slots, LoRA, metrics, router models, and reranking use native routes without OpenAI fallback.
- [x] 23. Add provider-specific unions, comments, fuzz targets, benchmarks, and documentation. **CP-23:** JSON fields have developer-facing GoDoc, generated easyjson files are reproducible, provider examples and capability documentation describe the supported surfaces.
- [x] 24. Run the final quality gates for the Anthropic and llama.cpp additions. **CP-24:** generation, lint, unit tests, race tests, vet, module verification, and diff checks pass; vulnerability scanning remains best-effort when its database is unavailable.

## Quality gate log

| Gate | Result |
|---|---|
| `go generate ./...` | passed |
| `make lint` | passed (`0 issues`; `govulncheck` warning: vulnerability DB unavailable) |
| `make tests` | passed (`go test ./...`) |
| `go test -race ./...` | passed |
| `go vet ./...` | passed |
| `go mod verify` | passed |
| `git diff --check` | passed |
| Stream/provider fuzz smoke | passed (`pkg/stream:FuzzSSEDoesNotPanic`, `anthropic/messages:FuzzStreamEventJSON`, `llama:FuzzNativePromptJSON`, 1s each) |
| Stream/native benchmarks | passed (`go test ./pkg/stream ./llama -run=^$ -bench=. -benchtime=1x`) |
