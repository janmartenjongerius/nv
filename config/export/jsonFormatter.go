package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"janmarten.name/env/config"
	"strings"
)

const (
	FormatJson = "json"
	FormatJsonPretty = "json:pretty"
)

type JsonFormatter struct {
	Formatter
	indent string
	indentSize uint8
	prefix string
}

func (formatter JsonFormatter) Format(variable *config.Variable) (string, error) {
	return formatter.FormatList([]*config.Variable{variable})
}

func (formatter JsonFormatter) FormatList(variables []*config.Variable) (string, error) {
	result := make(map[string]string, len(variables))

	for _, v := range variables {
		result[v.Key] = fmt.Sprintf("%v", v.Value)
	}

	marshalled, e := json.Marshal(result)

	if e == nil && formatter.indentSize > 0 && len(formatter.indent) > 0 {
		var buffer bytes.Buffer
		e = json.Indent(
			&buffer,
			marshalled,
			formatter.prefix,
			strings.Repeat(
				formatter.indent,
				int(formatter.indentSize),
			),
		)

		if e == nil {
			marshalled = buffer.Bytes()
		}
	}

	return string(marshalled), e
}

func init () {
	RegisterFormatter(FormatJson, JsonFormatter{})
	RegisterFormatter(FormatJsonPretty, JsonFormatter{indent: "\t", indentSize: 1})
}
