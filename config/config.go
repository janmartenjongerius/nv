/*
TODO:
	> Write unit tests
 */

package config

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// A struct representing a configuration entry by Key and Value.
type Variable struct {
	Key   string
	Value interface{}
}

// A map of configuration Variable entries.
type Config map[string]*Variable

// Get the Variable for the given Variable.Key.
func (c Config) Variable(key string) (v *Variable, e error) {
	if v, ok := c[key]; ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("No variable '%s' found", key))
}

// Set a Variable on the current Config.
func (c Config) SetVariable(v *Variable) {
	c[v.Key] = v
}

// Delete the Variable that matches the given Variable.
func (c Config) DeleteVariable(v *Variable) {
	delete(c, v.Key)
}

// Get the keys registered in the current Config.
func (c Config) Keys() []string {
	var keys []string

	for _, k := range reflect.ValueOf(c).MapKeys() {
		keys = append(keys, k.String())
	}

	sort.Strings(keys)

	return keys
}

// Create a new Config with the given Variable objects.
func NewConfig(vars ...*Variable) *Config {
	c := Config{}

	for _, v := range vars {
		c.SetVariable(v)
	}

	return &c
}

// Parse the given variable strings into a Config object.
func Parse(variables []string, separator string) *Config {
	config := NewConfig()

	for _, v := range variables {
		var tuple = strings.SplitN(v, separator, 2)

		config.SetVariable(&Variable{tuple[0], tuple[1]})
	}

	return config
}

// Map the Config to a map of interface{} types.
func (c Config) MapInterfaces() map[string]*interface{} {
	keys := c.Keys()
	exported := make(map[string]*interface{}, len(keys))

	for _, key := range keys {
		i := reflect.ValueOf(c[key]).Interface()
		exported[key] = &i
	}

	return exported
}
