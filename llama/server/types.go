// Package server provides typed native server administration operations.
package server

import "encoding/json"

//go:generate easyjson -all types.go

// HealthResponse is the readiness response returned by the server.
type HealthResponse struct {
	// Status is normally "ok" when the model is ready.
	Status string `json:"status"`
}

// Modalities reports the model input modalities exposed by the server.
type Modalities struct {
	// Vision reports whether image input is supported.
	Vision bool `json:"vision,omitempty"`
	// Audio reports whether audio input is supported.
	Audio bool `json:"audio,omitempty"`
}

// GenerationSettings contains the effective defaults for native completion.
type GenerationSettings struct {
	// NPredict is the default maximum number of generated tokens.
	NPredict int `json:"n_predict,omitempty"`
	// NCtx is the default context size.
	NCtx int `json:"n_ctx,omitempty"`
	// Temperature is the default sampling temperature.
	Temperature float64 `json:"temperature,omitempty"`
	// TopK is the default top-k sampling limit.
	TopK int `json:"top_k,omitempty"`
	// TopP is the default nucleus sampling threshold.
	TopP float64 `json:"top_p,omitempty"`
	// MinP is the default minimum probability threshold.
	MinP float64 `json:"min_p,omitempty"`
	// Stream reports the default streaming setting.
	Stream bool `json:"stream,omitempty"`
	// Stop contains the default stop sequences.
	Stop []string `json:"stop,omitempty"`
	// Samplers contains the default sampler order.
	Samplers []string `json:"samplers,omitempty"`
}

// Props contains the server's global model and generation properties.
type Props struct {
	// ModelPath is the loaded model file path.
	ModelPath string `json:"model_path,omitempty"`
	// ChatTemplate is the Jinja template used to format chat messages.
	ChatTemplate string `json:"chat_template,omitempty"`
	// TotalSlots is the number of concurrent processing slots.
	TotalSlots int `json:"total_slots,omitempty"`
	// Modalities lists the model capabilities such as vision and audio input.
	Modalities Modalities `json:"modalities,omitempty"`
	// IsSleeping reports whether the model is suspended because it is idle.
	IsSleeping bool `json:"is_sleeping,omitempty"`
	// DefaultGenerationSettings contains defaults applied to native completion requests.
	DefaultGenerationSettings GenerationSettings `json:"default_generation_settings,omitempty"`
	// ChatTemplateCapabilities contains provider-defined chat-template capability flags.
	ChatTemplateCapabilities json.RawMessage `json:"chat_template_caps,omitempty"`
	// MediaMarker is the placeholder expected for multimodal media payloads.
	MediaMarker string `json:"media_marker,omitempty"`
	// BuildInfo identifies the server build.
	BuildInfo string `json:"build_info,omitempty"`
}
