// Package slots provides typed slot API operations.
package slots

import (
	"encoding/json"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

const (
	// ActionSave stores the current slot state in the requested file.
	ActionSave = "save"
	// ActionRestore loads slot state from the requested file.
	ActionRestore = "restore"
	// ActionErase removes the slot state from the server cache.
	ActionErase = "erase"
)

//go:generate easyjson -all types.go

// Slot reports the current generation/cache state of one server slot.
type Slot struct {
	// ID is the slot identifier used in save, restore, and erase operations.
	ID int `json:"id_slot"`
	// State describes the server's current state for the slot.
	State string `json:"state,omitempty"`
	// PromptTokens is the slot context capacity in tokens.
	PromptTokens int `json:"n_ctx"`
	// Timings contains the provider-defined timing counters for this slot.
	Timings json.RawMessage `json:"timings,omitempty"`
}

// ListResponse is the native endpoint's array of slot states.
type ListResponse []Slot

// MarshalJSON encodes the native slot array.
func (response ListResponse) MarshalJSON() ([]byte, error) {
	type responseAlias []Slot
	return json.Marshal(responseAlias(response))
}

// UnmarshalJSON decodes the native slot array.
func (response *ListResponse) UnmarshalJSON(data []byte) error {
	type responseAlias []Slot
	var value responseAlias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*response = ListResponse(value)
	return nil
}

// MarshalEasyJSON implements easyjson.Marshaler for the slot array.
func (response ListResponse) MarshalEasyJSON(writer *jwriter.Writer) {
	data, err := response.MarshalJSON()
	writer.Raw(data, err)
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler for the slot array.
func (response *ListResponse) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	data := lexer.Raw()
	if lexer.Ok() {
		lexer.AddError(response.UnmarshalJSON(data))
	}
}

// ActionResponse reports the result of a slot cache operation.
type ActionResponse struct {
	// ID is the slot identifier affected by the operation.
	ID int `json:"id_slot"`
	// Filename is the cache file used by save or restore.
	Filename string `json:"filename,omitempty"`
	// Saved is the number of tokens written by a save operation.
	Saved int `json:"n_saved,omitempty"`
	// Restored is the number of tokens read by a restore operation.
	Restored int `json:"n_restored,omitempty"`
	// Erased is the number of tokens removed by an erase operation.
	Erased int `json:"n_erased,omitempty"`
	// BytesWritten is the number of bytes written to a save file.
	BytesWritten int64 `json:"n_written,omitempty"`
	// BytesRead is the number of bytes read from a restore file.
	BytesRead int64 `json:"n_read,omitempty"`
	// Timings contains provider-defined save or restore timing counters.
	Timings json.RawMessage `json:"timings,omitempty"`
}
