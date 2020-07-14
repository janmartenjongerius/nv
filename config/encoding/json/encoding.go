/*
Package json implements an encoder and decoder for a JSON representation of config.Variable.

For the following input:
	[]*config.Variable{
		&config.Variable{
			Key: "HOME",
			Value: "C:\Users\Gopher",
		},
		&config.Variable{
			Key: "USERNAME",
			Value: "Gopher",
		},
	}

The encoder will output:
	{"HOME":"C:\\Users\\Gopher","USERNAME":"Gopher"}
*/
package json

import (
	"encoding/json"
	"fmt"
	"janmarten.name/env/config"
	"janmarten.name/env/config/encoding"
)

type encoder struct {
	encoding.Encoder
}

type buffer map[string]string

// Allows to Encode config.Variable structs into a byte sequence.
func (e encoder) Encode(variables ...*config.Variable) ([]byte, error) {
	payload := make(buffer)

	for _, v := range variables {
		payload[v.Key] = fmt.Sprintf("%v", v.Value)
	}

	return json.Marshal(payload)
}

type decoder struct {
	encoding.Decoder
}

// Allows to Decode a byte sequence into a list of config.Variable structs.
func (d decoder) Decode(payload []byte) ([]*config.Variable, error) {
	var e error

	vars := make(buffer, 0)
	result := make([]*config.Variable, 0)

	if e = json.Unmarshal(payload, &vars); e == nil {
		for key, value := range vars {
			result = append(result, &config.Variable{
				Key:   key,
				Value: value,
			})
		}
	}

	return result, e
}

func init() {
	encoding.Register(
		"json",
		struct {
			encoder
			decoder
		}{},
	)
}
