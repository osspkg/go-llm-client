package llama_test

import (
	"encoding/json"
	"testing"

	"go.osspkg.com/llm-client/llama/completions"
	"go.osspkg.com/llm-client/llama/tokenization"
)

func FuzzNativePromptJSON(f *testing.F) {
	f.Add([]byte(`"hello"`))
	f.Add([]byte(`[1,"text",2]`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var prompt completions.Prompt
		_ = json.Unmarshal(data, &prompt)
	})
}

func FuzzTokenPieceJSON(f *testing.F) {
	f.Add([]byte(`{"id":1,"piece":"hello"}`))
	f.Add([]byte(`{"id":2,"piece":[195,161]}`))
	f.Add([]byte(`3`))
	f.Fuzz(func(t *testing.T, data []byte) {
		var piece tokenization.TokenPiece
		_ = json.Unmarshal(data, &piece)
	})
}
