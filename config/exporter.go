/*
Package config provides an Exporter to export Variable structs to a provided io.Writer.
 */
package config

import (
	"io"
)

type Exporter interface {
	Export(variable ...*Variable)
}

type ioExporter struct {
	output  io.Writer
	format string
}

func (exporter ioExporter) Export(variable ...*Variable) {
	WithEncoding(exporter.format, func(enc Encoding) {
		result, _ := enc.Encode(variable...)

		// If the list is not empty, ensure a newline at the end.
		if len(result) > 0 {
			result = append(result, []byte("\n")...)
		}

		_, _ = exporter.output.Write(result)
	})
}

func NewExporter(format string, writer io.Writer) Exporter {
	return ioExporter{writer, format}
}
