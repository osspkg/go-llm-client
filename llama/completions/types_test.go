package completions_test

import (
	"encoding/json"
	"testing"

	"go.osspkg.com/llm-client/llama/completions"
)

func TestPromptConstructorsReturnMarshalErrors(t *testing.T) {
	if _, err := completions.MixedPrompt([]json.RawMessage{json.RawMessage("{")}); err == nil {
		t.Fatal("MixedPrompt returned nil error for invalid raw JSON")
	}
}

func TestPromptConstructorsEncodeValues(t *testing.T) {
	text, err := completions.StringPrompt("hello")
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := completions.TokenPrompt([]int{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	multimodal, err := completions.MultimodalPrompt("<image>", []string{"AQI="})
	if err != nil {
		t.Fatal(err)
	}
	for name, prompt := range map[string]completions.Prompt{
		"text": text, "tokens": tokens, "multimodal": multimodal,
	} {
		if !json.Valid(prompt) {
			t.Fatalf("%s prompt is invalid JSON: %s", name, prompt)
		}
	}
}
