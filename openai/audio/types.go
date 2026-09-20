package audio

//go:generate easyjson -all types.go

// Transcription is a transcription response.
type Transcription struct {
	Text     string  `json:"text"`
	Language string  `json:"language,omitempty"`
	Duration float64 `json:"duration,omitempty"`
}

// SpeechRequest creates speech audio.
type SpeechRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Voice  string `json:"voice"`
	Format string `json:"response_format,omitempty"`
}

// Translation is a translation response.
type Translation struct {
	Text string `json:"text"`
}
