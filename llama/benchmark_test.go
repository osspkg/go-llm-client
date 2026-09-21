package llama_test

import (
	"encoding/json"
	"testing"

	"go.osspkg.com/llm-client/llama/completions"
	"go.osspkg.com/llm-client/llama/embeddings"
)

func BenchmarkNativePromptMarshal(b *testing.B) {
	prompt, err := completions.StringPrompt("benchmark prompt")
	if err != nil {
		b.Fatal(err)
	}
	request := completions.Request{Prompt: prompt, NPredict: 32}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(request); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNativeEmbeddingUnmarshal(b *testing.B) {
	data := []byte(`[{"index":0,"embedding":[[0.1,0.2,0.3]]}]`)
	b.ReportAllocs()
	for b.Loop() {
		var response embeddings.Response
		if err := json.Unmarshal(data, &response); err != nil {
			b.Fatal(err)
		}
	}
}
