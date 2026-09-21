# Справочник API и архитектуры

`go.osspkg.com/llm-client` содержит независимые provider-native клиенты. Каждый
клиент выполняет типизированные HTTP- или WebSocket-вызовы API провайдера; Go-код
не запускает модель и не вызывает Go callback изнутри LLM.

## Провайдеры и границы пакетов

- `openai` — OpenAI-compatible REST, полный закреплённый public/admin REST,
  media, uploads и Realtime WebSocket.
- `anthropic` — стабильные Anthropic Messages, Count Tokens, Models, Files и
  Message Batches. Agents и другие managed-agent beta-домены в этот этап не
  входят.
- `llama` — native REST API llama.cpp server. OpenAI-compatible `/v1/*`
  остаётся доступным через `openai.WithBaseURL` и здесь не дублируется.
- `ollama` — официальный Ollama HTTP API.
- `pkg/*` — независимые от провайдеров transport, authentication, typed errors,
  pagination, stream parsers, codec и WebSocket lifecycle.

Доменные пакеты владеют native-моделями запросов, ответов и событий своего
провайдера: `anthropic/messages`, `anthropic/models`, `anthropic/files`,
`anthropic/batches` и соответствующие пакеты `llama/*`. Это сохраняет wire
semantics и делает неподдерживаемые возможности явными.

## Создание клиента и авторизация

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

Anthropic по умолчанию использует `https://api.anthropic.com/v1` и отправляет
`anthropic-version: 2023-06-01`. `WithAPIKey` добавляет `x-api-key`, а
`WithBearerToken` — `Authorization: Bearer`. Workspace- и beta-заголовки
задаются явными опциями. `WithAuthProvider` позволяет возвращать собственные
заголовки для каждой операции.

Native llama.cpp по умолчанию использует `http://localhost:8080`; его auth
callback также вызывается перед каждым запросом. `RequestMeta.Domain` всегда
содержит hostname назначения, в том числе для custom base URL. Credentials не
попадают в URL и сообщения ошибок.

## Anthropic Messages

```go
response, err := anthropicClient.Messages.Create(ctx, messages.Request{
	Model:     "claude-3-5-sonnet-latest",
	MaxTokens: 256,
	Messages: []messages.Message{{
		Role:    "user",
		Content: messages.TextContent("Объясни SSE одним предложением."),
	}},
})
```

Messages предоставляет типизированные структуры для text, image, document,
thinking, redacted-thinking, tool-use и tool-result блоков. Поля
`Message.Content` и другие настоящие Anthropic string-or-array или
unions представлены явными provider union wrappers; schema-defined JSON,
например tool arguments, по-прежнему использует `json.RawMessage`.

`CreateStream` возвращает `stream.Iterator[messages.StreamEvent]`. Общий SSE
parser обрабатывает Anthropic event frames, `[DONE]`, границы событий по пустой
строке, ограничение размера, malformed JSON, отмену контекста и явный `Close`.
`Batches.Results` возвращает ограниченный NDJSON iterator с результатами
завершённого batch. Скачивание содержимого файла возвращает `io.ReadCloser`,
которым владеет вызывающий код; upload читает `io.Reader` через ограниченный
multipart request.

## Native llama.cpp

Native-клиент предоставляет:

- `Completions` для `/completion`: string, token-array, mixed и multimodal
  prompt, sampling settings, cache/slot/LoRA overrides, typed timings, stop
  metadata и SSE streaming;
- `Embeddings` для `/embeddings` и singular alias `/embedding`;
- `Tokenization` для `/tokenize` и `/detokenize`, включая union token piece в
  форме строки или массива байт;
- `Templates` для `/apply-template`;
- `Server` для `/health`, `/props` и изменения `/props`;
- `Slots` для списка и save/restore/erase prompt cache;
- `Lora` для списка адаптеров и изменения их scale;
- `Metrics` для ограниченного Prometheus text response;
- `Models` для router list, download, load, unload и model lifecycle SSE;
- `Rerank` для native reranking endpoint.

В конкретной сборке llama.cpp native endpoint может быть выключен или отсутствовать.
Клиент вернёт `errors.CapabilityError`, а исходный ограниченный `HTTPError`
останется доступен через `errors.As`; незаметного перехода на OpenAI- или
Anthropic-endpoint не происходит.

## Потоки и владение ресурсами

Все типизированные HTTP-потоки реализуют:

```go
type Iterator[T any] interface {
	Next(context.Context) bool
	Value() T
	Err() error
	Close() error
}
```

Iterator владеет своим response body. Вызывайте `Close` на каждом пути, в том
числе при досрочном завершении. Для callback-стиля есть `stream.ForEach`.
Parser не использует неограниченный scanner: response body и event lines
ограничены настройками transport и stream.

## Сериализация и ошибки

Пакеты моделей содержат `//go:generate easyjson -all types.go`, а сгенерированные
`*_easyjson.go` коммитятся. Запускайте `go generate ./...`; generated-файлы нельзя
редактировать вручную. Неизвестные JSON-поля игнорируются. `RawMessage`
используется только там, где upstream явно разрешает произвольный JSON: schemas,
tool input и документированные union payloads.

Общий transport задаёт лимиты dial, TLS handshake, response headers, request,
stream и body. Автоматические retries разрешены только для безопасных
идемпотентных методов; generation, uploads, streaming и WebSocket не повторяются.
Через `errors.Is` и `errors.As` можно различить validation, cancellation,
protocol, decode и `HTTPError`.

## Версии контрактов и разработка

Владельцы endpoint и upstream snapshots указаны в
[`docs/capability-matrix.md`](docs/capability-matrix.md). Для llama.cpp там
зафиксирован документированный snapshot ветки `master` и дата проверки; между
релизами llama.cpp набор optional native endpoints может измениться.

```text
go generate ./...
make lint
make tests
go test -race ./...
go vet ./...
go mod verify
git diff --check
```

Тесты используют deterministic `httptest` и локальные protocol fixtures. Живые
credentials провайдеров не требуются.
