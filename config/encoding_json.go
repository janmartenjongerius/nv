package config

import (
	"encoding/json"
)

type jsonEncoder struct {
	Encoder
}

type buffer map[string]string

/*
Encode allows to encode Variable structs into a byte sequence.

JSON encoding
By default, an encoder and decoder for a JSON representation of Variable is
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
	{
		"HOME":"C:\\Users\\Gopher",
		"USERNAME":"Gopher"
	}
*/
func (e jsonEncoder) Encode(variables ...*Variable) ([]byte, error) {
	payload := make(buffer)

	for _, v := range variables {
		payload[v.Key] = v.Value
	}

	return json.MarshalIndent(payload, "", "\t")
}

type jsonDecoder struct {
	Decoder
}

// Decode allows to decode a byte sequence into a list of Variables.
func (d jsonDecoder) Decode(payload []byte) (result Variables, err error) {
	vars := make(buffer, 0)
	result = make(Variables, 0)

	if err = json.Unmarshal(payload, &vars); err == nil {
		for key, value := range vars {
			result = append(result, &Variable{
				Key:   key,
				Value: value,
			})
		}
	}

	return result, err
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
