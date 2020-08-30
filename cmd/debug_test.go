package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"regexp"
	"testing"
)

func TestDebugCmd_RunE(t *testing.T) {
	out := &bytes.Buffer{}
	cmd := &cobra.Command{}
	cmd.SetOut(out)

	err := debugCmd.RunE(cmd, []string{})

	if err != nil {
		t.Fatal(err)
	}

	// Only the cmd package is loaded during this test. Therefore, the output
	// only shows debugging information from commands that use packages which
	// register debugging callbacks.
	// This is desirable, because plugins will interfere with this output.
	//
	// When the whole project is tested at-once, the Main and Plugins context are
	// appended to the end of this output. Asserting their structure is not that
	// beneficial, so only their entries under Callbacks are tested.
	got := out.String()
	want := regexp.MustCompile(
		`(?sm)^Nv debug:

Config\s+Env\s+\[.+?]

Debug\s+Callbacks\s+\[Config Debug Encoding( Main Plugins)?]

Encoding\s+Default\s+text
\s+Formats\s+\[.+?]`)

	if !want.MatchString(got) {
		t.Errorf("Could not match output against %q:\n%v", want, got)
	}
}
