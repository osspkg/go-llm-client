package blobs

//go:generate easyjson -all types.go

// Upload identifies an Ollama blob upload.
type Upload struct {
	Digest string `json:"digest"`
}
