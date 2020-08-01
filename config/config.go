package config

import (
	"os"
)

// The environment variables expressed in a map of Variable entries.
var Environment = make(Variables, 0)

// A struct representing a configuration entry by Key and Value.
type Variable struct {
	Key   string
	Value interface{}
}

// A map of configuration Variable entries.
type Variables []*Variable

// Initialize the environment variables.
func init() {
	defer WithEncoding("text", func(enc Encoding) {
		for _, line := range os.Environ() {
			vars, _ := enc.Decode([]byte(line))
			Environment = append(Environment, vars...)
		}
	})
}
