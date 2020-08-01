package config

import (
	"io"
)

type Exporter interface {
	Export(variable *Variable) error
	ExportList(variables []*Variable) error
}

type fileExporter struct {
	output  io.Writer
	encoder Encoder
}

func (exporter fileExporter) write(payload string) error {
	if _, e := exporter.output.Write([]byte(payload)); e != nil {
		return e
	}

	return nil
}

func (exporter fileExporter) Export(variable *Variable) error {
	var (
		result []byte
		e      error
	)

	if result, e = exporter.encoder.Encode(variable); e != nil {
		return e
	}

	return exporter.write(string(result) + "\n")
}

func (exporter fileExporter) ExportList(variables []*Variable) error {
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

func NewExporter(format string, writer io.Writer) (Exporter, error) {
	var (
		encoder Encoder
		e       error
	)

	if encoder, e = NewEncoding(format); e != nil {
		return nil, e
	}

	return fileExporter{writer, encoder}, nil
}

type ExporterFactory func(format string) (Exporter, error)

func NewExporterFactory(writer io.Writer) ExporterFactory {
	return func(format string) (Exporter, error) {
		return NewExporter(format, writer)
	}
}
