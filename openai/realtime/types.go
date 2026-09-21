/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package realtime

import "encoding/json"

//go:generate easyjson -all types.go

// Event is a typed subset of the OpenAI Realtime event envelope. Fields not
// applicable to an event type remain zero-valued.
type Event struct {
	// Type identifies the object or event variant.
	Type string `json:"type"`
	// EventID identifies the client or server event.
	EventID string `json:"event_id,omitempty"`
	// Session contains the provider value used for the `Session` field.
	Session *SessionConfig `json:"session,omitempty"`
	// Response contains generated text or a generated response object.
	Response *Response `json:"response,omitempty"`
	// Item contains the item submitted to a grader.
	Item *Item `json:"item,omitempty"`
	// Delta contains incremental text or argument content.
	Delta string `json:"delta,omitempty"`
	// Audio contains base64-encoded audio data.
	Audio string `json:"audio,omitempty"`
	// Text contains transcribed, generated, or annotated text.
	Text string `json:"text,omitempty"`
	// ContentPart contains the provider value used for the `ContentPart` field.
	ContentPart *ContentPart `json:"content_part,omitempty"`
	// ConversationID identifies the conversation associated with the event.
	ConversationID string `json:"conversation_id,omitempty"`
	// OutputIndex identifies the output item affected by an event.
	OutputIndex int `json:"output_index,omitempty"`
	// ContentIndex identifies the content part affected by an event.
	ContentIndex int `json:"content_index,omitempty"`
	// ItemID identifies the conversation item affected by an event.
	ItemID string `json:"item_id,omitempty"`
	// Arguments contains the JSON arguments supplied to a function.
	Arguments string `json:"arguments,omitempty"`
	// ResponseStatus contains the provider value used for the `ResponseStatus` field.
	ResponseStatus string `json:"status,omitempty"`
	// Error contains provider error details.
	Error *Error `json:"error,omitempty"`
	// Data contains the returned records.
	Data json.RawMessage `json:"data,omitempty"`
}

// SessionConfig configures a Realtime session.
type SessionConfig struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id,omitempty"`
	// Model selects the model used for the operation.
	Model string `json:"model,omitempty"`
	// Voice selects the voice used for speech generation.
	Voice string `json:"voice,omitempty"`
	// Instructions contains system-level guidance for response generation.
	Instructions string `json:"instructions,omitempty"`
	// Tools declares tools the model may call.
	Tools []Tool `json:"tools,omitempty"`
	// ToolChoice contains the provider value used for the `ToolChoice` field.
	ToolChoice json.RawMessage `json:"tool_choice,omitempty"`
	// Modalities contains the provider value used for the `Modalities` field.
	Modalities []string `json:"modalities,omitempty"`
	// InputAudioFormat contains the provider value used for the `InputAudioFormat` field.
	InputAudioFormat string `json:"input_audio_format,omitempty"`
	// OutputAudioFormat contains the provider value used for the `OutputAudioFormat` field.
	OutputAudioFormat string `json:"output_audio_format,omitempty"`
}

// Response configures a response generation event.
type Response struct {
	// Modalities contains the provider value used for the `Modalities` field.
	Modalities []string `json:"modalities,omitempty"`
	// Instructions contains system-level guidance for response generation.
	Instructions string `json:"instructions,omitempty"`
	// Voice selects the voice used for speech generation.
	Voice string `json:"voice,omitempty"`
	// OutputAudioFormat contains the provider value used for the `OutputAudioFormat` field.
	OutputAudioFormat string `json:"output_audio_format,omitempty"`
	// Tools declares tools the model may call.
	Tools []Tool `json:"tools,omitempty"`
}

// Tool declares a Realtime function tool.
type Tool struct {
	// Type identifies the object or event variant.
	Type string `json:"type"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Description explains a tool or resource so the model can use it correctly.
	Description string `json:"description,omitempty"`
	// Parameters contains the JSON Schema or arguments accepted by a function.
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// Item is a conversation item.
type Item struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id,omitempty"`
	// Type identifies the object or event variant.
	Type string `json:"type,omitempty"`
	// Role identifies the organization or message role.
	Role string `json:"role,omitempty"`
	// Content contains the message or content-part text.
	Content []ContentPart `json:"content,omitempty"`
	// CallID contains the provider value used for the `CallID` field.
	CallID string `json:"call_id,omitempty"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Arguments contains the JSON arguments supplied to a function.
	Arguments string `json:"arguments,omitempty"`
	// Output contains structured items produced by the response.
	Output string `json:"output,omitempty"`
}

// ContentPart is one content item.
type ContentPart struct {
	// Type identifies the object or event variant.
	Type string `json:"type,omitempty"`
	// Text contains transcribed, generated, or annotated text.
	Text string `json:"text,omitempty"`
	// Audio contains base64-encoded audio data.
	Audio string `json:"audio,omitempty"`
	// Transcript contains the provider value used for the `Transcript` field.
	Transcript string `json:"transcript,omitempty"`
}

// Error is a Realtime protocol error.
type Error struct {
	// Type identifies the object or event variant.
	Type string `json:"type,omitempty"`
	// Code contains the provider error code.
	Code string `json:"code,omitempty"`
	// Message contains the provider error message.
	Message string `json:"message,omitempty"`
	// Param identifies the request parameter associated with the error.
	Param string `json:"param,omitempty"`
}
