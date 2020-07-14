/*
TODO:
	> Write unit tests
	> Document symbols
	> Create YAML formatter as plugin (.so)
 */
package export

import (
	"errors"
	"io"
	"janmarten.name/env/config"
	"sort"
)

const UnknownFormatterError = "unknown format requested"

var formatters = make(map[string]Formatter)

type Exporter interface {
	Export(variable *config.Variable) error
	ExportList(variables []*config.Variable) error
}

/*
TODO:
	> Rename Formatter -> Encoding and put in new package: env/config/encoding
	> Create YAML encoding as plugin (.so)
	> Create ASCII encoding as plugin (.so)
 */
type Formatter interface {
	Format(variable *config.Variable) (string, error)
	FormatList(variables []*config.Variable) (string, error)
}

type fileExporter struct {
	output    io.Writer
	formatter Formatter
}

func (exporter fileExporter) write(payload string) error {
	if _, e := exporter.output.Write([]byte(payload)); e != nil {
		return e
	}

	return nil
}

func (exporter fileExporter) Export(variable *config.Variable) error {
	var (
		result string
		e      error
	)

	if result, e = exporter.formatter.Format(variable); e != nil {
		return e
	}

	return exporter.write(result + "\n")
}

func (exporter fileExporter) ExportList(variables []*config.Variable) error {
	var (
		result string
		e      error
	)

	if result, e = exporter.formatter.FormatList(variables); e != nil {
		return e
	}

	if len(result) > 0 {
		result += "\n"
	}

	return exporter.write(result)
}

func RegisterFormatter(format string, formatter Formatter) {
	formatters[format] = formatter
}

func GetFormats() []string {
	keys := make([]string, 0)

	for k := range formatters {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func New(format string, writer io.Writer) (Exporter, error) {
	var (
		formatter Formatter
		ok        bool
	)

	if formatter, ok = formatters[format]; ok == false {
		return nil, errors.New(UnknownFormatterError)
	}

	return fileExporter{writer, formatter}, nil
}

type Factory func(format string) (Exporter, error)

func NewFactory(writer io.Writer) Factory {
	return func(format string) (Exporter, error) {
		return New(format, writer)
	}
}
