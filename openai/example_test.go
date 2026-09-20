package openai_test

import (
	"context"

	"go.osspkg.com/llm-client/openai"
	"go.osspkg.com/llm-client/openai/chat"
	"go.osspkg.com/llm-client/pkg/auth"
)

func ExampleNew() {
	client, err := openai.New(openai.WithAuthProvider(auth.StaticBearer("token")))
	if err != nil {
		return
	}
	_, _ = client.Chat.Create(context.Background(), chat.Request{Model: "model"})
}
