/*
Package text implements an encoder and decoder for a text based representation of config.Variable.

For the following input:
	[]*config.Variable{
		&config.Variable{
			Key: "HOME",
			Value: "C:\Users\Gopher",
		},
		&config.Variable{
			Key: "USERNAME",
			Value: "GOPHER",
		},
	}

The encoder will output:
	HOME=C:\Users\Gopher
	USERNAME=Gopher
 */
package text

import (
	"fmt"
	"janmarten.name/env/config"
	"janmarten.name/env/config/encoding"
	"strings"
)

type encoder struct {
	encoding.Encoder
}

// Allows to Encode config.Variable structs into a byte sequence.
func (e encoder) Encode(variables ...*config.Variable) ([]byte, error) {
	result := make([]string, len(variables))

	for i, v := range variables {
		result[i] = fmt.Sprintf("%s=%s", v.Key, v.Value)
	}

	return []byte(strings.Join(result, "\n")), nil
}

type decoder struct {
	encoding.Decoder
}

// Allows to Decode a byte sequence into a list of config.Variable structs.
func (d decoder) Decode(payload []byte) ([]*config.Variable, error) {
	variables := make([]*config.Variable, 0)

	for _, line := range strings.Split(string(payload), "\n") {
		if len(line) == 0 {
			continue
		}

		components := strings.SplitN(line, "=", 2)

		if len(components) != 2 {
			return nil, encoding.IllFormattedVariable
		}

		variables = append(variables, &config.Variable{
			Key:   components[0],
			Value: components[1],
		})
	}

	return variables, nil
}

func init() {
	encoding.Register(
		"text",
		struct {
			encoder
			decoder
		}{},
	)
}
