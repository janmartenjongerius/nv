/*
Exporting variables
Using an Exporter, Variable structs can be written a provided io.Writer, for a given format.
 */
package config

import (
	"io"
)

// Exporter describes structs that are able to export Variable entries.
type Exporter interface {
	Export(variable ...*Variable)
}

type ioExporter struct {
	output  io.Writer
	format string
}

// Export the given Variable entries.
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

// Create a new Exporter for the given format, using the given io.Writer as target.
func NewExporter(format string, writer io.Writer) Exporter {
	return ioExporter{writer, format}
}
