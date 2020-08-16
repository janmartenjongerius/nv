/*
Package config describes the structure of a Variable and lists variables for the current Environment.

During the init-phase of the application, the Environment variable is populated with the result of os.Environ.
*/
package config

import (
	"os"
)

// Environment holds the environment variables expressed in Variables.
var Environment = make(Variables, 0)

// DefaultEncoding is the default encoding format.
var DefaultEncoding = "text"

// Variable is a struct representing a configuration entry by Key and Value.
type Variable struct {
	Key   string
	Value string
}

// Variables are a list of configuration Variable entries.
type Variables []*Variable

// Initialize the environment variables.
func init() {
	WithEncoding(DefaultEncoding, func(enc Encoding) {
		for _, line := range os.Environ() {
			vars, _ := enc.Decode([]byte(line))
			Environment = append(Environment, vars...)
		}
	})
}
