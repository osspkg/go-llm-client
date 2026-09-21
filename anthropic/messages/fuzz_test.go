package messages_test

import (
	"encoding/json"
	"testing"

	"go.osspkg.com/llm-client/anthropic/messages"
)

func FuzzStreamEventJSON(f *testing.F) {
	f.Add([]byte(`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var event messages.StreamEvent
		_ = json.Unmarshal(data, &event)
	})
}

func FuzzContentBlockJSON(f *testing.F) {
	f.Add([]byte(`{"type":"text","text":"hello"}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var block messages.TextBlock
		_ = json.Unmarshal(data, &block)
	})
}
