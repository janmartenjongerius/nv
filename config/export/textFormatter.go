package export

import (
	"fmt"
	"janmarten.name/env/config"
	"strings"
)

const FormatText = "text"

type TextFormatter struct {
	Formatter
}

func (formatter TextFormatter) Format(variable *config.Variable) (string, error) {
	return fmt.Sprintf("%s=%s", variable.Key, variable.Value), nil
}

func (formatter TextFormatter) FormatList(variables []*config.Variable) (string, error) {
	var e error
	result := make([]string, len(variables))

	for i, v := range variables {
		if result[i], e = formatter.Format(v); e != nil {
			return "", e
		}
	}

	return strings.Join(result, "\n"), nil
}

func init() {
	RegisterFormatter(FormatText, TextFormatter{})
}
