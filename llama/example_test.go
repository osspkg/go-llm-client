package llama_test

import (
	"context"

	"go.osspkg.com/llm-client/llama"
	"go.osspkg.com/llm-client/llama/completions"
)

func ExampleClient_completions() {
	client, err := llama.New(llama.WithBaseURL("http://127.0.0.1:18081"))
	if err != nil {
		return
	}
	_, _ = client.Completions.Create(context.Background(), completions.Request{
		Prompt:   completions.StringPrompt("Hello"),
		NPredict: 32,
	})
}
