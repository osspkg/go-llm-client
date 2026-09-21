package ollama_test

import (
	"context"

	"go.osspkg.com/llm-client/ollama"
	"go.osspkg.com/llm-client/ollama/generate"
)

func ExampleNew() {
	client, err := ollama.New()
	if err != nil {
		return
	}
	_, _ = client.Generate().Create(context.Background(), generate.Request{Model: "llama3"})
}
