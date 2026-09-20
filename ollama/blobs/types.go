package blobs

//go:generate easyjson -all types.go

// Upload identifies an Ollama blob upload.
type Upload struct {
	// Digest identifies an immutable blob or model content by SHA-256 digest.
	Digest string `json:"digest"`
}
