# Справочник API и архитектуры

Библиотека `go.osspkg.com/llm-client` содержит два provider-native клиента:
OpenAI-compatible и Ollama.

```go
client, err := openai.New(
	openai.WithBaseURL("https://api.openai.com/v1"),
	openai.WithAuthProvider(auth.StaticBearer("token")),
)
if err != nil {
	return err
}

completion, err := client.Chat.Create(ctx, chat.Request{Model: "model"})
```

Провайдер предоставляет HTTP- или WebSocket-API. Go-библиотека выполняет
типизированные transport-вызовы этого API; Go-клиент не запускает модель и не
вызывает Go callback изнутри LLM.

## Границы пакетов

Пакеты `openai/*` и `ollama/*` являются bounded contexts: они владеют
provider-native моделями запросов, ответов и операциями. В `pkg/*` находятся
независимые от провайдеров компоненты transport, auth, errors, pagination,
stream, codec и WebSocket lifecycle. Типы провайдеров не импортируются в `pkg`.

Такое разделение сохраняет wire semantics каждого API и не скрывает
provider-specific возможности за общей нормализованной facade.

## Потоки и владение ресурсами

HTTP-потоки реализуют `stream.Iterator[T]`:

```go
type Iterator[T any] interface {
	Next(context.Context) bool
	Value() T
	Err() error
	Close() error
}
```

SSE распознаёт кадры `data:` и маркер `[DONE]`. Ollama использует ограниченные
значения NDJSON, разделённые переводом строки. Каждый iterator владеет своим
`response body`; вызывающий код обязан вызвать `Close`, в том числе при
досрочном прекращении обработки.

Realtime-сессии используют lifecycle `Connect`, `Send`, `Receive` и `Close`.
Чтение и запись сериализуются общим WebSocket lifecycle wrapper, поэтому
операции имеют явного владельца и контролируемое завершение.

## Сериализация

Importable-пакеты моделей содержат директиву
`//go:generate easyjson -all types.go`, а сгенерированный код `*_easyjson.go`
хранится в репозитории. Generated-файлы нельзя редактировать вручную.

Неизвестные JSON-поля игнорируются. `json.RawMessage` используется только для
полей, которым upstream-спецификация намеренно разрешает произвольный JSON,
например для схем и аргументов инструментов.

## Таймауты, лимиты, повторы и ошибки

Общий transport задаёт ограничение размера response и request body, таймауты
dial, TLS handshake и response headers, передаёт context и поддерживает
настраиваемые retry. Автоматические повторы отключены для generation,
upload, stream и WebSocket operations; повторяться могут только безопасные
идемпотентные методы.

HTTP-ошибки сохраняют HTTP status и request ID, но не сохраняют полный body
провайдера и credentials. Для обработки ошибок используйте `errors.Is` и
`errors.As`: они позволяют различать validation, protocol, decode,
cancellation и HTTP errors.

## Версии контрактов

Закреплённые upstream revisions, владельцы endpoint’ов и текущие возможности
указаны в [матрице capability](docs/capability-matrix.md).
