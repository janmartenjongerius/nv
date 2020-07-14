/*
TODO:
	> Write unit tests
	> Document symbols
*/
package export

import (
	"io"
	"janmarten.name/env/config"
	"janmarten.name/env/config/encoding"
)

type Exporter interface {
	Export(variable *config.Variable) error
	ExportList(variables []*config.Variable) error
}

type fileExporter struct {
	output  io.Writer
	encoder encoding.Encoder
}

func (exporter fileExporter) write(payload string) error {
	if _, e := exporter.output.Write([]byte(payload)); e != nil {
		return e
	}

	return nil
}

func (exporter fileExporter) Export(variable *config.Variable) error {
	var (
		result []byte
		e      error
	)

	if result, e = exporter.encoder.Encode(variable); e != nil {
		return e
	}

	return exporter.write(string(result) + "\n")
}

func (exporter fileExporter) ExportList(variables []*config.Variable) error {
	var (
		encoded []byte
		result  string
		e       error
	)

	if encoded, e = exporter.encoder.Encode(variables...); e != nil {
		return e
	}

	result = string(encoded)

	if len(result) > 0 {
		result += "\n"
	}

	return exporter.write(result)
}

func New(format string, writer io.Writer) (Exporter, error) {
	var (
		encoder encoding.Encoder
		e       error
	)

	if encoder, e = encoding.New(format); e != nil {
		return nil, e
	}

	return fileExporter{writer, encoder}, nil
}

type Factory func(format string) (Exporter, error)

func NewFactory(writer io.Writer) Factory {
	return func(format string) (Exporter, error) {
		return New(format, writer)
	}
}
