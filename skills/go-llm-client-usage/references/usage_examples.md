# Usage examples

The examples use provider-native types and the current accessor-based root API.
They are patterns for application code; replace model names and endpoints with
values supported by the selected server.

## OpenAI-compatible chat

```go
ctx := context.Background()
client, err := openai.New(
	openai.WithAuthProvider(auth.StaticBearer(os.Getenv("OPENAI_API_KEY"))),
)
if err != nil {
	return err
}

response, err := client.Chat().Create(ctx, chat.Request{
	Model: "model",
	Messages: []chat.Message{{
		Role:    "user",
		Content: json.RawMessage(`"Hello"`),
	}},
})
if err != nil {
	return err
}
fmt.Println(string(response.Choices[0].Message.Content))
```

For an OpenAI-compatible local or hosted server, add
`openai.WithBaseURL("https://example.test/v1")`. OpenAI-compatible llama.cpp
routes use this client too.

## Anthropic Messages

```go
content, err := messages.TextContent("Explain SSE briefly.")
if err != nil {
	return err
}

client, err := anthropic.New(
	anthropic.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
)
if err != nil {
	return err
}

response, err := client.Messages().Create(ctx, messages.Request{
	Model:     "claude-3-5-sonnet-latest",
	MaxTokens: 256,
	Messages: []messages.Message{{Role: "user", Content: content}},
})
if err != nil {
	return err
}
fmt.Println(response.ID)
```

Use `anthropic.WithBearerToken` or a custom `WithAuthProvider` when the
compatible server does not use `x-api-key`. `RequestMeta.Domain` remains
available to the callback for routing credentials.

## Native llama.cpp completion

```go
prompt, err := completions.StringPrompt("Write one short sentence.")
if err != nil {
	return err
}

client, err := llama.New(llama.WithBaseURL("http://localhost:8080"))
if err != nil {
	return err
}
response, err := client.Completions().Create(ctx, completions.Request{
	Prompt:   prompt,
	NPredict: 32,
})
if err != nil {
	return err
}
fmt.Println(response.Content)
```

Use `llama.Tokenization()`, `Server()`, `Slots()`, `LoRA()`, `Models()`, and
`Metrics()` only for native routes. Do not send OpenAI `/v1/*` requests through
this root.

## Server environment metadata

Environment metadata does not require a client instance:

```go
for _, variable := range llama.EnvironmentScheme().Variables {
	fmt.Printf("%s (default %q): %s\n", variable.Name, variable.Default, variable.Description)
}
```

Use `ollama.EnvironmentScheme()` for Ollama server variables. The returned
entries are descriptions only; the library does not read or set environment
variables. `AllowedValues` is empty for documented scalar formats such as
paths, durations, integers, and comma-separated lists.

## Ollama NDJSON stream

```go
client, err := ollama.New()
if err != nil {
	return err
}

events, err := client.Chat().CreateStream(ctx, chat.Request{
	Model:    "llama3.2",
	Messages: []chat.Message{{Role: "user", Content: "Hello"}},
})
if err != nil {
	return err
}
defer events.Close()

for events.Next(ctx) {
	fmt.Print(events.Value().Message.Content)
}
return events.Err()
```

The Ollama `chat` import in this example is `go.osspkg.com/llm-client/ollama/chat`.
OpenAI and Anthropic streams use the same iterator lifecycle, but their event
types and SSE payloads are provider-specific.

## Dynamic auth with destination domain

```go
client, err := openai.New(openai.WithAuthProvider(func(ctx context.Context, meta auth.RequestMeta) (http.Header, error) {
	if meta.Domain == "internal.example.test" {
		return auth.StaticBearer(os.Getenv("INTERNAL_LLM_TOKEN"))(ctx, meta)
	}
	return make(http.Header), nil
}))
```

The callback is called for every request. Do not cache a mutable `http.Header`
map and do not put the selected token into a URL.

## OpenAI Realtime WebSocket

```go
client, err := openai.New(
	openai.WithAuthProvider(auth.StaticBearer(os.Getenv("OPENAI_API_KEY"))),
)
if err != nil {
	return err
}

session, err := client.Realtime().Connect(ctx, "gpt-realtime")
if err != nil {
	return err
}
defer session.Close()

if err := session.Send(ctx, realtime.Event{Type: "session.update"}); err != nil {
	return err
}
event, err := session.Receive(ctx)
if err != nil {
	return err
}
fmt.Println(event.Type)
```

Realtime owns one reader and serializes writes. Keep `Connect`, `Send`,
`Receive`, and `Close` under caller-controlled contexts and always close the
session.

## Typed errors

```go
var httpErr *llmerrors.HTTPError
if errors.As(err, &httpErr) {
	fmt.Println(httpErr.StatusCode, httpErr.RequestID)
}

if errors.Is(err, llmerrors.ErrBodyTooLarge) {
	// Increase the explicit limit or handle the bounded failure.
}
```

For a missing native llama.cpp route, inspect `*llmerrors.CapabilityError`; the
client never silently switches to another provider's endpoint.

## Streaming with a callback

```go
return stream.ForEach(ctx, events, func(event chat.StreamChunk) error {
	fmt.Print(event.Choices[0].Delta.Content)
	return nil
})
```

`stream.ForEach` closes the iterator on return. If manual iteration is used,
defer `Close` immediately after a successful stream creation.
