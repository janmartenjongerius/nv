// +build !windows

/*
Package main implements an config.Encoder and config.Decoder for an XML representation of config.Variable.

For the following input:
	config.Variables{
		&config.Variable{
			Key: "HOME",
			Value: "C:\Users\Gopher",
		},
		&config.Variable{
			Key: "USERNAME",
			Value: "Gopher",
		},
	}

The config.Encoder will output:
	<Variable><Key>HOME</Key><Value>C:\Users\Gopher</Value></Variable><Variable><Key>USERNAME</Key><Value>Gopher</Value></Variable>
*/
package main

import (
	"encoding/xml"
	"janmarten.name/env/config"
)

type xmlEncoder struct {
	config.Encoder
}

// Allows to Encode config.Variable structs into a byte sequence.
func (e xmlEncoder) Encode(variables ...*config.Variable) ([]byte, error) {
	return xml.Marshal(variables)
}

type xmlDecoder struct {
	config.Decoder
}

// Allows to Decode a byte sequence into a list of config.Variables.
func (d xmlDecoder) Decode(payload []byte) (config.Variables, error) {
	result := make([]*config.Variable, 0)
	e := xml.Unmarshal(payload, &result)

	return result, e
}

func init() {
	config.RegisterEncoding(
		"xml",
		struct {
			xmlEncoder
			xmlDecoder
		}{},
	)
}
