package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"janmarten.name/nv/config"
	"reflect"
	"strconv"
	"testing"
)

func TestSearchCmd_RunE(t *testing.T) {
	cmd := &cobra.Command{}

	cases := []struct {
		name           string
		env            config.Variables
		numSuggestions uint64
		args           []string
		err            bool
		stdout         string
		stderr         string
	}{
		{
			name:           "Empty environment - With suggestions",
			args:           []string{"HOME"},
			numSuggestions: 5,
			err:            true,
			stderr:         "Could not find HOME.\n  No suggestions.\n",
		},
		{
			name:           "Empty environment - No suggestions",
			args:           []string{"HOME"},
			numSuggestions: 0,
			err:            true,
			stderr:         "Could not find HOME.\n  No suggestions.\n",
		},
		{
			name:           "Single variable - Match",
			args:           []string{"HOME"},
			numSuggestions: 0,
			stdout:         "HOME=/home/gopher\n",
			env: config.Variables{
				&config.Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
			},
		},
		{
			name:           "One match, one fail - No suggestions",
			args:           []string{"HOME", "USER"},
			numSuggestions: 0,
			err:            true,
			stdout:         "HOME=/home/gopher\n",
			stderr:         "Could not find USER.\n  No suggestions.\n",
			env: config.Variables{
				&config.Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
			},
		},
		{
			name:           "One match, one fail - Single suggestion",
			args:           []string{"HOME", "USER"},
			numSuggestions: 1,
			err:            true,
			stdout:         "HOME=/home/gopher\n",
			stderr:         "Could not find USER.\n  Suggestions:\n   - USER_LEVEL\n",
			env: config.Variables{
				&config.Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
				&config.Variable{
					Key:   "USER_LEVEL",
					Value: "1",
				},
			},
		},
		{
			name:           "One match, one fail - Multiple suggestions",
			args:           []string{"HOME", "USER"},
			numSuggestions: 1,
			err:            true,
			stdout:         "HOME=/home/gopher\n",
			stderr:         "Could not find USER.\n  Suggestions:\n   - USER_LEVEL\n",
			env: config.Variables{
				&config.Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
				&config.Variable{
					Key:   "USER_LEVEL",
					Value: "1",
				},
				&config.Variable{
					Key:   "USER_LEVELS",
					Value: "1,2,3,4",
				},
			},
		},
		{
			name:           "One match, one automatic correction",
			args:           []string{"HOME", "USR"},
			numSuggestions: 3,
			stdout:         "HOME=/home/gopher\nUSER=gopher\n",
			env: config.Variables{
				&config.Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
				&config.Variable{
					Key:   "USER",
					Value: "gopher",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config.Environment = c.env

			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			cmd.SetOut(stdout)
			cmd.SetErr(stderr)

			_ = searchCmd.
				Flag("num-suggestions").
				Value.
				Set(strconv.FormatUint(c.numSuggestions, 10))

			err := searchCmd.RunE(cmd, c.args)

			if reflect.DeepEqual(err, nil) == c.err {
				t.Errorf("Expected error? %v, got: %v", c.err, err)
			}

			if c.stdout != stdout.String() {
				t.Errorf(
					"Unexpected output.\nGot:\n%q\nWant:\n%q",
					stdout.String(),
					c.stdout,
				)
			}

			if c.stderr != stderr.String() {
				t.Errorf(
					"Unexpected error output.\nGot:\n%q\nWant:\n%q",
					stderr.String(),
					c.stderr,
				)
			}
		})
	}
}
