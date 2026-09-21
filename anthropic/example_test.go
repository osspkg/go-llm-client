package anthropic_test

import (
	"context"

	"go.osspkg.com/llm-client/anthropic"
	"go.osspkg.com/llm-client/anthropic/messages"
)

func ExampleClient_messages() {
	client, err := anthropic.New(anthropic.WithBaseURL("http://127.0.0.1:18080/v1"))
	if err != nil {
		return
	}
	_, _ = client.Messages.Create(context.Background(), messages.Request{
		Model:     "model",
		MaxTokens: 32,
		Messages:  []messages.Message{{Role: "user", Content: messages.TextContent("Hello")}},
	})
}
