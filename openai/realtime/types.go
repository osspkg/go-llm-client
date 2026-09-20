package realtime

import "encoding/json"

//go:generate easyjson -all types.go

// Event is a typed subset of the OpenAI Realtime event envelope. Fields not
// applicable to an event type remain zero-valued.
type Event struct {
	Type           string          `json:"type"`
	EventID        string          `json:"event_id,omitempty"`
	Session        *SessionConfig  `json:"session,omitempty"`
	Response       *Response       `json:"response,omitempty"`
	Item           *Item           `json:"item,omitempty"`
	Delta          string          `json:"delta,omitempty"`
	Audio          string          `json:"audio,omitempty"`
	Text           string          `json:"text,omitempty"`
	ContentPart    *ContentPart    `json:"content_part,omitempty"`
	ConversationID string          `json:"conversation_id,omitempty"`
	OutputIndex    int             `json:"output_index,omitempty"`
	ContentIndex   int             `json:"content_index,omitempty"`
	ItemID         string          `json:"item_id,omitempty"`
	Arguments      string          `json:"arguments,omitempty"`
	ResponseStatus string          `json:"status,omitempty"`
	Error          *Error          `json:"error,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
}

// SessionConfig configures a Realtime session.
type SessionConfig struct {
	ID                string          `json:"id,omitempty"`
	Model             string          `json:"model,omitempty"`
	Voice             string          `json:"voice,omitempty"`
	Instructions      string          `json:"instructions,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        json.RawMessage `json:"tool_choice,omitempty"`
	Modalities        []string        `json:"modalities,omitempty"`
	InputAudioFormat  string          `json:"input_audio_format,omitempty"`
	OutputAudioFormat string          `json:"output_audio_format,omitempty"`
}

// Response configures a response generation event.
type Response struct {
	Modalities        []string `json:"modalities,omitempty"`
	Instructions      string   `json:"instructions,omitempty"`
	Voice             string   `json:"voice,omitempty"`
	OutputAudioFormat string   `json:"output_audio_format,omitempty"`
	Tools             []Tool   `json:"tools,omitempty"`
}

// Tool declares a Realtime function tool.
type Tool struct {
	Type        string          `json:"type"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// Item is a conversation item.
type Item struct {
	ID        string        `json:"id,omitempty"`
	Type      string        `json:"type,omitempty"`
	Role      string        `json:"role,omitempty"`
	Content   []ContentPart `json:"content,omitempty"`
	CallID    string        `json:"call_id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Arguments string        `json:"arguments,omitempty"`
	Output    string        `json:"output,omitempty"`
}

// ContentPart is one content item.
type ContentPart struct {
	Type       string `json:"type,omitempty"`
	Text       string `json:"text,omitempty"`
	Audio      string `json:"audio,omitempty"`
	Transcript string `json:"transcript,omitempty"`
}

// Error is a Realtime protocol error.
type Error struct {
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   string `json:"param,omitempty"`
}
