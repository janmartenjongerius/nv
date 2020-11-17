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
	<Config>
		<Variable key="HOME">C:\Users\Gopher</Variable>
		<Variable key="USERNAME">Gopher</Variable>
	</Config>

*/
package config

import (
	"encoding/xml"
)

// config.Encoder for config.Variables represented in XML.
type xmlEncoder struct {
	Encoder
}

// A struct representing configuration Variable entries.
type xmlConfig struct {
	XMLName   xml.Name     `xml:"Config"`
	Variables xmlVariables `xml:"Variable"`
}

type xmlVariables []*xmlVariable

type xmlVariable struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

// Encode allows to encode config.Variable structs into a byte sequence.
func (e xmlEncoder) Encode(variables ...*Variable) ([]byte, error) {
	vars := xmlVariables{}

	for _, v := range variables {
		vars = append(vars, &xmlVariable{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return xml.MarshalIndent(xmlConfig{Variables: vars}, "", "\t")
}

// config.Decoder for config.Variables represented in XML.
type xmlDecoder struct {
	Decoder
}

// Decode allows to decode a byte sequence into a list of config.Variables.
func (d xmlDecoder) Decode(payload []byte) (Variables, error) {
	cfg := xmlConfig{}
	result := Variables{}
	e := xml.Unmarshal(payload, &cfg)

	for _, v := range cfg.Variables {
		result = append(result, &Variable{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return result, e
}

func init() {
	RegisterEncoding(
		"xml",
		struct {
			xmlEncoder
			xmlDecoder
		}{},
	)
}
