package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestCompletionCmd_RunE(t *testing.T) {
	cmd := &cobra.Command{}
	cases := []struct {
		in   string
		want string
	}{
		{
			in:   "bash",
			want: "# bash completion for",
		},
		{
			in:   "zsh",
			want: "#compdef _",
		},
		{
			in:   "fish",
			want: "# fish completion for",
		},
		{
			in:   "powershell",
			want: "using namespace System.Management.Automation",
		},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			err := completionCmd.RunE(cmd, []string{c.in})

			got := buf.String()

			if !strings.HasPrefix(got, c.want) {
				t.Errorf("Got: %q\nWant: %q ...", got, c.want)
			}

			if err != nil {
				t.Errorf("Command provided unexpected error: %v", err)
			}
		})
	}

	err := completionCmd.RunE(cmd, []string{"unknown-shell"})

	if err == nil {
		t.Errorf("Expected program to produce error for args: %q", []string{"unknown-shell"})
	}
}
