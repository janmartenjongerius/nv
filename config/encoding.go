package config

import (
	"errors"
	"sort"
)

var (
	// ErrUnknownEncoding is returned when an unknown encoding is requested using NewEncoding.
	ErrUnknownEncoding = errors.New("unknown encoding requested")

	// ErrIllFormattedVariable is returned when a payload contains a line that cannot be resolved to a config.Variable.
	ErrIllFormattedVariable = errors.New("variable is incorrectly formatted")

	// Registry of Encoding instances.
	encodings = make(map[string]Encoding)

	// Map of Callbacks for when an Encoding gets registered.
	encodingCallbacks = make(map[string][]EncodingCallback)
)

// EncodingCallback is a callback for when an Encoding has become available.
type EncodingCallback func(enc Encoding)

// Encoder allows to Encode config.Variable structs into a byte sequence.
type Encoder interface {
	Encode(variables ...*Variable) ([]byte, error)
}

// Decoder allows to Decode a byte sequence into a list of config.Variable structs.
type Decoder interface {
	Decode(payload []byte) (Variables, error)
}

// Encoding is a shared interface for Encoder and Decoder.
type Encoding interface {
	Encoder
	Decoder
}

// RegisterEncoding registers the given Encoding in a registry for NewEncoding instances.
// Additionally, it triggers available EncodingCallback entries for the given name.
func RegisterEncoding(name string, encoding Encoding) {
	encodings[name] = encoding
	processEncodingCallbacks(name)
}

// WithEncoding registers an EncodingCallback to execute when the Encoding with
// the given name is/will be registered.
func WithEncoding(name string, callback EncodingCallback) {
	encodingCallbacks[name] = append(encodingCallbacks[name], callback)
	processEncodingCallbacks(name)
}

// Process the callbacks that are currently registered for the Encoding with the
// given name, if it exists.
func processEncodingCallbacks(name string) {
	enc, err := NewEncoding(name)

	if err != nil {
		return
	}

	defer func(Callbacks []EncodingCallback) {
		for _, Callback := range Callbacks {
			Callback(enc)
		}
	}(encodingCallbacks[name])
	encodingCallbacks[name] = nil
}

// NewEncoding creates a new Encoding for the given format.
func NewEncoding(format string) (Encoding, error) {
	encoding, ok := encodings[format]

	if !ok {
		return nil, ErrUnknownEncoding
	}

	return encoding, nil
}

// GetEncodings gets a list of available Encoding formats.
func GetEncodings() []string {
	keys := make([]string, 0)

	for k := range encodings {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// HasEncoding determines whether an encoding for the given format has been registered.
func HasEncoding(format string) bool {
	_, has := encodings[format]
	return has
}
