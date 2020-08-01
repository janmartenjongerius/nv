package config

import "bytes"

func ExampleExportJsonBuffer() {
	buf := bytes.Buffer{}
	exp := NewExporter("json", &buf)
	exp.Export(&Variable{Key: "FOO", Value: "Foo"})
}
