/*
Package config implements an Encoder and Decoder for a JSON representation of Variable.

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

The Encoder will output:
	{"HOME":"C:\\Users\\Gopher","USERNAME":"Gopher"}
*/
package config

import (
	"encoding/json"
	"fmt"
)

type jsonEncoder struct {
	Encoder
}

type buffer map[string]string

// Allows to Encode Variable structs into a byte sequence.
func (e jsonEncoder) Encode(variables ...*Variable) ([]byte, error) {
	payload := make(buffer)

	for _, v := range variables {
		payload[v.Key] = fmt.Sprintf("%v", v.Value)
	}

	return json.Marshal(payload)
}

type jsonDecoder struct {
	Decoder
}

// Allows to Decode a byte sequence into a list of Variables.
func (d jsonDecoder) Decode(payload []byte) (Variables, error) {
	var e error

	vars := make(buffer, 0)
	result := make([]*Variable, 0)

	if e = json.Unmarshal(payload, &vars); e == nil {
		for key, value := range vars {
			result = append(result, &Variable{
				Key:   key,
				Value: value,
			})
		}
	}

	return result, e
}

func init() {
	RegisterEncoding(
		"json",
		struct {
			jsonEncoder
			jsonDecoder
		}{},
	)
}
