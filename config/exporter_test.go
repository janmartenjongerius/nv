package config

import (
	"os"
)

func ExampleNewExporter() {
	exp := NewExporter("text", os.Stdout)
	exp.Export(&Variable{
		Key:   "FOO",
		Value: "Foo",
	})
	// Output:
	// FOO=Foo
}

func ExampleNewExporter_json() {
	exp := NewExporter("json", os.Stdout)
	exp.Export(&Variable{
		Key:   "FOO",
		Value: "Foo",
	})
	// Output:
	// {
	// 	"FOO": "Foo"
	// }
}
