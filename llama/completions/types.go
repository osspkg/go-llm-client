// Package completions provides typed native completion operations.
package completions

import (
	"encoding/json"
	"errors"
	"fmt"
)

//go:generate easyjson -all types.go

// Prompt is the native prompt union: text, token IDs, mixed values, or multimodal data.
// Use StringPrompt, TokenPrompt, MixedPrompt, or MultimodalPrompt to construct it.
type Prompt []byte

// StringPrompt creates a text prompt and returns any JSON encoding error.
func StringPrompt(value string) (Prompt, error) {
	return marshalPrompt(value)
}

// TokenPrompt creates a prompt from token IDs and returns any JSON encoding error.
func TokenPrompt(tokens []int) (Prompt, error) {
	return marshalPrompt(tokens)
}

// MixedPrompt creates a provider-defined mixed token/string prompt and validates
// its raw JSON values.
func MixedPrompt(values []json.RawMessage) (Prompt, error) {
	return marshalPrompt(values)
}

// MultimodalPrompt creates a prompt object containing text and base64 media
// payloads and returns any JSON encoding error.
func MultimodalPrompt(text string, media []string) (Prompt, error) {
	return marshalPrompt(struct {
		// PromptString is the text containing media marker placeholders.
		PromptString string `json:"prompt_string"`
		// MultimodalData contains base64-encoded media in marker order.
		MultimodalData []string `json:"multimodal_data,omitempty"`
	}{PromptString: text, MultimodalData: media})
}

func marshalPrompt(value any) (Prompt, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal llama prompt: %w", err)
	}
	return Prompt(data), nil
}

// MarshalJSON validates and returns the prompt union unchanged.
func (prompt Prompt) MarshalJSON() ([]byte, error) {
	if len(prompt) == 0 {
		return []byte("null"), nil
	}
	if !json.Valid(prompt) {
		return nil, errors.New("invalid llama prompt JSON")
	}
	return append([]byte(nil), prompt...), nil
}

// UnmarshalJSON stores a validated prompt union.
func (prompt *Prompt) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		return errors.New("invalid llama prompt JSON")
	}
	*prompt = append((*prompt)[:0], data...)
	return nil
}

// Request contains native completion and sampling options.
type Request struct {
	// Prompt is text, token IDs, a mixed token/string sequence, or a multimodal object.
	Prompt Prompt `json:"prompt"`
	// Model selects a model in router mode.
	Model string `json:"model,omitempty"`
	// Stream requests incremental SSE responses.
	Stream bool `json:"stream,omitempty"`
	// NPredict limits the number of generated tokens; -1 means no explicit limit.
	NPredict int `json:"n_predict,omitempty"`
	// NCmpl requests multiple completions for each prompt.
	NCmpl int `json:"n_cmpl,omitempty"`
	// Temperature controls sampling randomness.
	Temperature float64 `json:"temperature,omitempty"`
	// DynamicTemperatureRange enables dynamic temperature around Temperature.
	DynamicTemperatureRange float64 `json:"dynatemp_range,omitempty"`
	// DynamicTemperatureExponent controls the dynamic temperature curve.
	DynamicTemperatureExponent float64 `json:"dynatemp_exponent,omitempty"`
	// TopK limits sampling to the K most likely tokens.
	TopK int `json:"top_k,omitempty"`
	// TopP applies nucleus sampling.
	TopP float64 `json:"top_p,omitempty"`
	// MinP rejects tokens below a probability relative to the best token.
	MinP float64 `json:"min_p,omitempty"`
	// Stop contains strings that terminate generation and are omitted from output.
	Stop []string `json:"stop,omitempty"`
	// TypicalP enables locally typical sampling.
	TypicalP float64 `json:"typical_p,omitempty"`
	// RepeatPenalty controls repetition suppression.
	RepeatPenalty float64 `json:"repeat_penalty,omitempty"`
	// RepeatLastN is the context window used by repetition suppression.
	RepeatLastN int `json:"repeat_last_n,omitempty"`
	// PresencePenalty penalizes tokens that already appeared.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`
	// FrequencyPenalty scales penalties by token frequency.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	// DryMultiplier controls the DRY repetition penalty.
	DryMultiplier float64 `json:"dry_multiplier,omitempty"`
	// DryBase is the exponential base used by DRY.
	DryBase float64 `json:"dry_base,omitempty"`
	// DryAllowedLength is the repetition length before DRY increases its penalty.
	DryAllowedLength int `json:"dry_allowed_length,omitempty"`
	// DryPenaltyLastN limits the input range scanned by DRY.
	DryPenaltyLastN int `json:"dry_penalty_last_n,omitempty"`
	// DrySequenceBreakers contains strings that break repeated sequences.
	DrySequenceBreakers []string `json:"dry_sequence_breakers,omitempty"`
	// Mirostat selects the Mirostat sampler variant, or zero to disable it.
	Mirostat int `json:"mirostat,omitempty"`
	// MirostatTau sets the target entropy for Mirostat.
	MirostatTau float64 `json:"mirostat_tau,omitempty"`
	// MirostatEta sets the learning rate for Mirostat.
	MirostatEta float64 `json:"mirostat_eta,omitempty"`
	// Grammar is the grammar used to constrain generated output.
	Grammar string `json:"grammar,omitempty"`
	// JSONSchema is the schema-defined JSON grammar constraint.
	JSONSchema json.RawMessage `json:"json_schema,omitempty"`
	// Seed selects the random seed; -1 delegates seed selection to the server.
	Seed int64 `json:"seed,omitempty"`
	// NProbs requests token probability details in the response.
	NProbs int `json:"n_probs,omitempty"`
	// ReturnTokens asks the server to include generated token IDs.
	ReturnTokens bool `json:"return_tokens,omitempty"`
	// ResponseFields limits the fields emitted in each response.
	ResponseFields []string `json:"response_fields,omitempty"`
	// CachePrompt reuses a matching prompt prefix from the slot cache.
	CachePrompt bool `json:"cache_prompt,omitempty"`
	// IDSlot selects a specific server slot.
	IDSlot int `json:"id_slot,omitempty"`
	// LoRA overrides adapter scales for this request.
	LoRA []LoRAAdapter `json:"lora,omitempty"`
}

// LoRAAdapter selects an adapter and its per-request scale.
type LoRAAdapter struct {
	// ID is the adapter ID returned by the LoRA endpoint.
	ID int `json:"id"`
	// Scale is the adapter multiplier for this completion.
	Scale float64 `json:"scale"`
}

// Response is one native completion result or streamed token event.
type Response struct {
	// Content is the generated text, or the next text fragment in a stream.
	Content string `json:"content,omitempty"`
	// Tokens contains generated token IDs when requested or streaming is enabled.
	Tokens []int `json:"tokens,omitempty"`
	// Stop reports that the streamed generation has ended.
	Stop bool `json:"stop,omitempty"`
	// Model is the model alias used for generation.
	Model string `json:"model,omitempty"`
	// Prompt is the processed prompt when the server returns it.
	Prompt Prompt `json:"prompt,omitempty"`
	// StopType identifies why generation stopped: eos, limit, or word.
	StopType string `json:"stop_type,omitempty"`
	// StoppingWord is the stop sequence that ended generation, when applicable.
	StoppingWord string `json:"stopping_word,omitempty"`
	// GenerationSettings contains the effective server-side sampling settings.
	GenerationSettings GenerationSettings `json:"generation_settings,omitempty"`
	// Timings contains measured prompt and generation performance counters.
	Timings Timings `json:"timings,omitempty"`
	// CompletionProbabilities contains per-token probability details when requested.
	CompletionProbabilities []TokenProbability `json:"completion_probabilities,omitempty"`
	// Probabilities is the legacy/native short name used by some server versions.
	Probabilities []TokenProbability `json:"probs,omitempty"`
	// TokensCached is the number of prompt tokens reused from the slot cache.
	TokensCached int `json:"tokens_cached,omitempty"`
	// TokensEvaluated is the total number of prompt tokens evaluated.
	TokensEvaluated int `json:"tokens_evaluated,omitempty"`
	// Truncated reports that the prompt exceeded the context window.
	Truncated bool `json:"truncated,omitempty"`
}

// GenerationSettings describes effective completion parameters returned by llama.cpp.
type GenerationSettings struct {
	// NPredict is the effective generation token limit.
	NPredict int `json:"n_predict,omitempty"`
	// NCtx is the context size used by the slot.
	NCtx int `json:"n_ctx,omitempty"`
	// Model is the effective model alias.
	Model string `json:"model,omitempty"`
	// Temperature is the effective sampling temperature.
	Temperature float64 `json:"temperature,omitempty"`
	// TopK is the effective top-k value.
	TopK int `json:"top_k,omitempty"`
	// TopP is the effective nucleus threshold.
	TopP float64 `json:"top_p,omitempty"`
	// MinP is the effective minimum probability threshold.
	MinP float64 `json:"min_p,omitempty"`
	// Stop contains effective stopping strings.
	Stop []string `json:"stop,omitempty"`
	// Stream reports whether the effective request used streaming.
	Stream bool `json:"stream,omitempty"`
	// Samplers contains the effective sampler sequence.
	Samplers []string `json:"samplers,omitempty"`
}

// Timings contains native prompt and generation timing counters.
type Timings struct {
	// PromptN is the number of prompt tokens processed.
	PromptN int `json:"prompt_n,omitempty"`
	// PromptMS is prompt processing time in milliseconds.
	PromptMS float64 `json:"prompt_ms,omitempty"`
	// PromptPerSecond is prompt processing throughput.
	PromptPerSecond float64 `json:"prompt_per_second,omitempty"`
	// PredictedN is the number of generated tokens.
	PredictedN int `json:"predicted_n,omitempty"`
	// PredictedMS is generation time in milliseconds.
	PredictedMS float64 `json:"predicted_ms,omitempty"`
	// PredictedPerSecond is generation throughput.
	PredictedPerSecond float64 `json:"predicted_per_second,omitempty"`
}

// TokenProbability contains the probability and text representation of one token.
type TokenProbability struct {
	// ID is the vocabulary token ID.
	ID int `json:"id"`
	// LogProbability is the token log probability when log probabilities are requested.
	LogProbability float64 `json:"logprob,omitempty"`
	// Probability is the token probability when post-sampling probabilities are requested.
	Probability float64 `json:"prob,omitempty"`
	// Token is the UTF-8 token text.
	Token string `json:"token,omitempty"`
	// Bytes contains the raw token bytes for non-text-safe tokens.
	Bytes []int `json:"bytes,omitempty"`
	// TopLogProbabilities contains alternative token probabilities at this position.
	TopLogProbabilities []TokenProbability `json:"top_logprobs,omitempty"`
	// TopProbabilities contains alternative post-sampling probabilities.
	TopProbabilities []TokenProbability `json:"top_probs,omitempty"`
}
