package config

import (
	"janmarten.name/nv/debug"
	"os"
	"sort"
)

// Environment holds the environment variables expressed in Variables.
var Environment = make(Variables, 0)

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
	debug.RegisterCallback("Config", func() debug.Messages {
		return debug.Messages{
			"Env": func(vars Variables) []string {
				result := make([]string, 0)

				for _, v := range vars {
					result = append(result, v.Key)
				}

				sort.Strings(result)

				return result
			}(Environment),
		}
	})
}
