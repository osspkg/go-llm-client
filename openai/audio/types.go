package audio

import "io"

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

// VoiceConsentInput uploads a voice-consent recording.
//
//easyjson:skip
type VoiceConsentInput struct {
	Name      string
	Language  string
	Filename  string
	Recording io.Reader
}

// UpdateVoiceConsentRequest changes a voice-consent label.
type UpdateVoiceConsentRequest struct {
	Name string `json:"name"`
}

// VoiceConsent is a consent recording that authorizes a custom voice.
type VoiceConsent struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Language  string `json:"language"`
	CreatedAt int64  `json:"created_at"`
}

// VoiceConsentListResponse is a cursor page of voice-consent recordings.
type VoiceConsentListResponse struct {
	Object  string         `json:"object"`
	Data    []VoiceConsent `json:"data"`
	FirstID string         `json:"first_id,omitempty"`
	LastID  string         `json:"last_id,omitempty"`
	HasMore bool           `json:"has_more"`
}

// DeleteVoiceConsentResponse confirms deletion of a voice-consent recording.
type DeleteVoiceConsentResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// VoiceInput creates a custom voice from a consented audio sample.
//
//easyjson:skip
type VoiceInput struct {
	Name        string
	ConsentID   string
	Filename    string
	AudioSample io.Reader
}

// Voice is a custom voice resource.
type Voice struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}
