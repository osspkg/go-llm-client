package audio

import "io"

//go:generate easyjson -all types.go

// Transcription is a transcription response.
type Transcription struct {
	// Text contains transcribed, generated, or annotated text.
	Text string `json:"text"`
	// Language identifies the spoken or transcribed language.
	Language string `json:"language,omitempty"`
	// Duration contains audio duration in seconds.
	Duration float64 `json:"duration,omitempty"`
}

// SpeechRequest creates speech audio.
type SpeechRequest struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// Input contains the text, tokens, or input items to process.
	Input string `json:"input"`
	// Voice selects the voice used for speech generation.
	Voice string `json:"voice"`
	// Format selects the model or response format.
	Format string `json:"response_format,omitempty"`
}

// Translation is a translation response.
type Translation struct {
	// Text contains transcribed, generated, or annotated text.
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
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
}

// VoiceConsent is a consent recording that authorizes a custom voice.
type VoiceConsent struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// Language identifies the spoken or transcribed language.
	Language string `json:"language"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
}

// VoiceConsentListResponse is a cursor page of voice-consent recordings.
type VoiceConsentListResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []VoiceConsent `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// DeleteVoiceConsentResponse confirms deletion of a voice-consent recording.
type DeleteVoiceConsentResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
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
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
}
