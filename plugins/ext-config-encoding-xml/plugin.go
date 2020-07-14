/*
Package main implements an encoder and decoder for an XML representation of config.Variable.

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
	<Variable><Key>HOME</Key><Value>C:\Users\Gopher</Value></Variable><Variable><Key>USERNAME</Key><Value>Gopher</Value></Variable>
*/
package main

import (
	"encoding/xml"
	"janmarten.name/env/config"
	"janmarten.name/env/config/encoding"
)

type encoder struct {
	encoding.Encoder
}

// Allows to Encode config.Variable structs into a byte sequence.
func (e encoder) Encode(variables ...*config.Variable) ([]byte, error) {
	return xml.Marshal(variables)
}

type buffer []struct {
	Key   string
	Value string
}

type decoder struct {
	encoding.Decoder
}

// Allows to Decode a byte sequence into a list of config.Variable structs.
func (d decoder) Decode(payload []byte) ([]*config.Variable, error) {
	var e error

	vars := make(buffer, 0)
	result := make([]*config.Variable, 0)

	if e = xml.Unmarshal(payload, &vars); e != nil {
		for _, v := range vars {
			result = append(result, &config.Variable{
				Key:   v.Key,
				Value: v.Value,
			})
		}
	}

	return result, e
}

func init() {
	encoding.Register(
		"xml",
		struct {
			encoder
			decoder
		}{},
	)
}
