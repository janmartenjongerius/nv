/*
Package encoding provides interfaces for configuration encoders and decoders.

Additional encoders can be added through the RegisterEncoding function.
This allows for separation of concerns internally, but also for plugins to implement their own encoding.
*/
package config

import (
	"errors"
	"sort"
)

var (
	// UnknownEncodingError is returned when an unknown encoding is requested using NewEncoding.
	UnknownEncodingError = errors.New("unknown encoding requested")

	// IllFormattedVariable is returned when a payload contains a line that cannot be resolved to a config.Variable.
	IllFormattedVariable = errors.New("variable is incorrectly formatted")

	// Registry of Encoding instances.
	encodings = make(map[string]Encoding)

	// Map of Callbacks for when an Encoding gets registered.
	encodingCallbacks = make(map[string][]EncodingCallback)
)

// Callback for when an Encoding has become available.
type EncodingCallback func (enc Encoding)

// Allows to Encode config.Variable structs into a byte sequence.
type Encoder interface {
	Encode(variables ...*Variable) ([]byte, error)
}

// Allows to Decode a byte sequence into a list of config.Variable structs.
type Decoder interface {
	Decode(payload []byte) (Variables, error)
}

// Shared interface Encoding for Encoder and Decoder.
type Encoding interface {
	Encoder
	Decoder
}

// Register the given Encoding in a registry for NewEncoding instances.
func RegisterEncoding(name string, encoding Encoding) {
	encodings[name] = encoding
	processEncodingCallbacks(name)
}

// Register an EncodingCallback to execute when the Encoding with the given name is/will be registered.
func WithEncoding(name string, callback EncodingCallback) {
	encodingCallbacks[name] = append(encodingCallbacks[name], callback)
	processEncodingCallbacks(name)
}

// Process the callbacks that are currently registered for the Encoding with the given name, if it exists.
func processEncodingCallbacks(name string) {
	enc, err := NewEncoding(name)

	if err != nil {
		return
	}

	defer func (Callbacks []EncodingCallback) {
		for _, Callback := range Callbacks {
			Callback(enc)
		}
	}(encodingCallbacks[name])
	encodingCallbacks[name] = nil
}

// Create a new Encoding for the given format.
func NewEncoding(format string) (Encoding, error) {
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
