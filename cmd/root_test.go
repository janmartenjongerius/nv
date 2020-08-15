package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"janmarten.name/nv/config"
	"reflect"
	"testing"
)

func TestRootCmd_Run(t *testing.T) {
	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	rootCmd.Run(cmd, []string{})

	if out.Len() == 0 {
		t.Errorf("Nothing exported to buffer")
	}
}

func TestExecute(t *testing.T) {
	err := Execute("test")

	if err != nil {
		t.Errorf("Execute produced an error: %v", err)
	}
}

func TestFormatValue_Set(t *testing.T) {
	value := new(formatValue)
	cases := []struct{
		in string
		err bool
	}{
		{
			in:  "Foo",
			err: true,
		},
	}

	for _, format := range config.GetEncodings() {
		cases = append(
			cases,
			struct{
				in string
				err bool
			}{in: format, err: false})
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			err := value.Set(c.in)

			if reflect.DeepEqual(err, nil) == c.err {
				t.Errorf("Expected error? %v, got: %v", c.err, err)
			}

			if err == nil && value.String() != c.in {
				t.Errorf("Could not update value. Got %q, want %q", value.String(), c.in)
			}
		})
	}
}
