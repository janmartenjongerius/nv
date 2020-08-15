/*
Encoding
Variables can be converted to a textual representation and back, using Decoder and Encoder instances.

An Encoding is the combination of an Encoder and Decoder, which can be registered against a given format.

Additional encoders can be added through the RegisterEncoding function.
This allows for separation of concerns internally, but also for plugins to implement their own encoding.

The function NewEncoding allows to get the encoder that is registered against the provided format. This function will
return an error when no encoding is known for the given format, even if the encoding is registered later on.

To process Variables for a given format, without being concerned about whether the format exists, the function
WithEncoding can be used. Since an Encoding may be registered at runtime, this strategy allows to queue processing of
Variables before the Encoding is available. The moment the Encoding is registered, all corresponding EncodingCallback
functions are invoked. When the Encoding is already present, this happens the moment WithEncoding is invoked.

When explicitly checking against an available Encoding, the functions HasEncoding and GetEncodings will be of use.
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
type EncodingCallback func(enc Encoding)

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

	defer func(Callbacks []EncodingCallback) {
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

// Get whether an encoding for the given format has been registered.
func HasEncoding(format string) bool {
	_, has := encodings[format]
	return has
}
