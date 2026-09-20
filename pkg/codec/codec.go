// Package codec provides the repository's easyjson-only model codec boundary.
package codec

import (
	"errors"
	"fmt"

	"github.com/mailru/easyjson"
)

// Marshal encodes a generated easyjson model.
func Marshal(value easyjson.Marshaler) ([]byte, error) {
	if value == nil {
		return nil, nil
	}
	data, err := easyjson.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal model: %w", err)
	}
	return data, nil
}

// Unmarshal decodes a generated easyjson model.
func Unmarshal(data []byte, value easyjson.Unmarshaler) error {
	if value == nil {
		return errors.New("unmarshal model: nil target")
	}
	if err := easyjson.Unmarshal(data, value); err != nil {
		return fmt.Errorf("unmarshal model: %w", err)
	}
	return nil
}
