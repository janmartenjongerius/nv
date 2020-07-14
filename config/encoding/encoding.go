/*
Package encoding provides interfaces for configuration encoders and decoders.

Additional encoders can be added through the Register function.
This allows for separation of concerns internally, but also for plugins to implement their own encoding.
*/
package encoding

import (
	"errors"
	"janmarten.name/env/config"
	"sort"
)

var (
	// UnknownEncodingError is returned when an unknown encoding is requested using New.
	UnknownEncodingError = errors.New("unknown encoding requested")

	// IllFormattedVariable is returned when a payload contains a line that cannot be resolved to a config.Variable.
	IllFormattedVariable = errors.New("variable is incorrectly formatted")

	// Registry of Encoding instances.
	encodings = make(map[string]Encoding)
)

// Allows to Encode config.Variable structs into a byte sequence.
type Encoder interface {
	Encode(variables ...*config.Variable) ([]byte, error)
}

// Allows to Decode a byte sequence into a list of config.Variable structs.
type Decoder interface {
	Decode(payload []byte) ([]*config.Variable, error)
}

// Shared interface Encoding for Encoder and Decoder.
type Encoding interface {
	Encoder
	Decoder
}

// Register the given Encoding in a registry for New instances.
func Register(name string, encoding Encoding) {
	encodings[name] = encoding
}

// Create a new Encoding for the given format.
func New(format string) (Encoding, error) {
	encoding, ok := encodings[format]

	if !ok {
		return nil, UnknownEncodingError
	}

	return encoding, nil
}

// Get available Encoding names.
func GetEncodings() []string {
	keys := make([]string, 0)

	for k := range encodings {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
