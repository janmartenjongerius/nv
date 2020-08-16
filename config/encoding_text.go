package config

import (
	"fmt"
	"strings"
)

type textEncoder struct {
	Encoder
}

/*
Encode allows to encode Variable structs into a byte sequence.

Text encoding
By default, an encoder and decoder for a text representation of Variable is
registered.

For the following input:
	Variables{
		&Variable{
			Key: "HOME",
			Value: "C:\Users\Gopher",
		},
		&Variable{
			Key: "USERNAME",
			Value: "Gopher",
		},
	}

The encoder will output:
	HOME=C:\Users\Gopher
	USERNAME=Gopher
 */
func (e textEncoder) Encode(variables ...*Variable) ([]byte, error) {
	result := make([]string, len(variables))

	for i, v := range variables {
		result[i] = fmt.Sprintf("%s=%s", v.Key, v.Value)
	}

	return []byte(strings.Join(result, "\n")), nil
}

type textDecoder struct {
	Decoder
}

// Decode allows to decode a byte sequence into a list of Variables.
func (d textDecoder) Decode(payload []byte) (Variables, error) {
	variables := make([]*Variable, 0)

	for _, line := range strings.Split(string(payload), "\n") {
		if len(line) == 0 {
			continue
		}

		components := strings.SplitN(line, "=", 2)

		if len(components) != 2 {
			return nil, IllFormattedVariable
		}

		variables = append(variables, &Variable{
			Key:   components[0],
			Value: components[1],
		})
	}

	return variables, nil
}

func init() {
	RegisterEncoding(
		DefaultEncoding,
		struct {
			textEncoder
			textDecoder
		}{},
	)
}
