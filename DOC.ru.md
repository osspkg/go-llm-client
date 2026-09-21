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

Авторизация подписки Codex доступна через `auth.NewCodexDeviceAuth`. Передайте
в `Login` callback, который покажет `CodexDeviceCode.VerificationURL` и
`UserCode`, затем передайте `HeaderProvider()` полученной сессии в
`openai.WithAuthProvider`. Сессия обновляет access token, если срок его
действия известен. Device code и выданные tokens являются секретами: не
пишите их в логи и сохраняйте только в защищённом хранилище.

### Интеграция Codex с хранилищем и UI

`pkg/auth` намеренно не сохраняет credentials. Граница хранения принадлежит
приложению: используйте зашифрованные поля в базе данных, KMS-backed secret
store или системный keyring. Типичный интерфейс на уровне приложения:

```go
type CodexTokenStore interface {
	Load(context.Context, string) (auth.CodexTokens, error)
	Save(context.Context, string, auth.CodexTokens) error
	Delete(context.Context, string) error
}
```

`string` идентифицирует пользователя или другой subject приложения, а не
токен Codex. Хранилище должно отличать отсутствие записи от ошибки хранения и
никогда не писать поля токена в логи.

Для web UI обмен device code следует выполнять только на backend:

1. `POST /auth/codex/start` вызывает `RequestDeviceCode` и сохраняет полный
   `auth.CodexDeviceCode` на backend с коротким TTL и привязкой к текущему
   пользователю.
2. В браузер возвращаются только `VerificationURL`, `UserCode` и идентификатор
   состояния или срок действия. Нельзя сериализовать и затем собирать
   `CodexDeviceCode` заново: `Complete` использует его приватные
   device-auth identifier и polling interval.
3. UI открывает verification URL и показывает одноразовый код.
4. `POST /auth/codex/complete` загружает ожидающий device code, вызывает
   `Complete`, сохраняет `session.Tokens()` в защищённое хранилище и привязывает
   сессию к backend-сессии пользователя.
5. После перезапуска процесса токены загружаются из хранилища, затем перед
   созданием OpenAI-клиента вызывается `NewSession`:

```go
tokens, err := tokenStore.Load(ctx, subjectID)
if err != nil {
	return err
}

session, err := codexAuth.NewSession(tokens)
if err != nil {
	return err
}

client, err := openai.New(
	openai.WithAuthProvider(session.HeaderProvider()),
)
if err != nil {
	return err
}
```

Браузер получает только непрозрачную cookie с флагами `HttpOnly`, `Secure` и
подходящим `SameSite`. Access token, refresh token, ID token и ожидающий
device code остаются на backend. Endpoints запуска и завершения должны быть
защищены собственной authentication/authorization схемой приложения, CSRF
защитой, rate limit и проверкой принадлежности состояния пользователю.

`HeaderProvider` обновляет истекающий access token в памяти. После операции
клиента снова сохраняйте `session.Tokens()`, чтобы не потерять новый
refresh token при завершении процесса:

```go
if _, err := client.Chat().Create(ctx, request); err != nil {
	return err
}
if err := tokenStore.Save(ctx, subjectID, session.Tokens()); err != nil {
	return err
}
```

Если нужна немедленная персистентность, оберните `HeaderProvider`, сравнивайте
текущий snapshot токенов с последним сохранённым и записывайте данные только
после refresh. Обёртка должна передавать исходный context и не включать
credentials в ошибки или логи.

Для desktop UI используется та же backend-facing схема, но token set следует
хранить в системном keyring. UI-процесс получает статус авторизации и данные
для показа, но не access/refresh token. При logout удалите token set из
хранилища, отбросьте in-memory session и инвалидируйте сессию приложения.

Native llama.cpp по умолчанию использует `http://localhost:8080`; его auth
callback также вызывается перед каждым запросом. `RequestMeta.Domain` всегда
содержит hostname назначения, в том числе для custom base URL. Credentials не
попадают в URL и сообщения ошибок.

## Anthropic Messages

```go
content, err := messages.TextContent("Объясни SSE одним предложением.")
if err != nil {
	return err
}
response, err := anthropicClient.Messages().Create(ctx, messages.Request{
	Model:     "claude-3-5-sonnet-latest",
	MaxTokens: 256,
	Messages: []messages.Message{{
		Role:    "user",
		Content: content,
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

Native-клиент после создания неизменяем и предоставляет доменные клиенты
через accessors:

- `Completions()` для `/completion`: string, token-array, mixed и multimodal
  prompt, sampling settings, cache/slot/LoRA overrides, typed timings, stop
  metadata и SSE streaming;
- `Embeddings()` для `/embeddings` и singular alias `/embedding`;
- `Tokenization()` для `/tokenize` и `/detokenize`, включая union token piece в
  форме строки или массива байт;
- `Templates()` для `/apply-template`;
- `Server()` для `/health`, `/props` и изменения `/props`;
- `Slots()` для списка и save/restore/erase prompt cache;
- `LoRA()` для списка адаптеров и изменения их scale;
- `Metrics()` для ограниченного Prometheus text response;
- `Models()` для router list, download, load, unload и model lifecycle SSE;
- `Rerank()` для native reranking endpoint.

В конкретной сборке llama.cpp native endpoint может быть выключен или отсутствовать.
Клиент вернёт `errors.CapabilityError`, а исходный ограниченный `HTTPError`
останется доступен через `errors.As`; незаметного перехода на OpenAI- или
Anthropic-endpoint не происходит.

## Схемы переменных окружения сервера

Метаданные переменных окружения не привязаны к экземпляру клиента. Root-пакеты
предоставляют `ollama.EnvironmentScheme()` и `llama.EnvironmentScheme()`. Каждая
функция возвращает `pkg/environment.Scheme` со списком переменных сервера:

```go
scheme := llama.EnvironmentScheme()
for _, variable := range scheme.Variables {
	fmt.Printf("%s=%q: %s\n", variable.Name, variable.Default, variable.Description)
}
```

Каждый `environment.Variable` содержит точное имя переменной, конечный список
`AllowedValues`, если провайдер задаёт перечисление, значение `Default` и
понятное описание назначения и формата значения. Пустой `AllowedValues` означает
скалярный или зависящий от провайдера формат: путь, duration, число или список,
а не разрешение произвольного JSON. Функции не читают и не изменяют
`os.Environ` и не требуют создания root-клиента.

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
