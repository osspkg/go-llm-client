package messages_test

import (
	"encoding/json"
	"testing"

	"go.osspkg.com/llm-client/anthropic/messages"
)

func TestContentConstructorsReturnMarshalErrors(t *testing.T) {
	if _, err := messages.BlocksContent(messages.ContentBlock{Type: "tool_use", Input: json.RawMessage("{")}); err == nil {
		t.Fatal("BlocksContent returned nil error for invalid raw JSON")
	}
	if _, err := messages.ToolUseBlockValue("tool_1", "lookup", json.RawMessage("{")); err == nil {
		t.Fatal("ToolUseBlockValue returned nil error for invalid raw JSON")
	}
	if _, err := messages.ToolResultBlockValue("tool_1", json.RawMessage("{"), false); err == nil {
		t.Fatal("ToolResultBlockValue returned nil error for invalid raw JSON")
	}
}

func TestContentConstructorsEncodeTypedValues(t *testing.T) {
	text, err := messages.TextContent("hello")
	if err != nil {
		t.Fatal(err)
	}
	blocks, err := messages.BlocksContent(messages.TextBlockValue("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if got, err := text.DecodeText(); err != nil || got != "hello" {
		t.Fatalf("text = %q, err=%v", got, err)
	}
	decoded, err := blocks.DecodeBlocks()
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 1 || decoded[0].Type != "text" {
		t.Fatalf("blocks = %#v", decoded)
	}

	if _, err := messages.ImageBlockValue(messages.ImageSource{Type: "url", URL: "https://example.test/image"}); err != nil {
		t.Fatal(err)
	}
}
