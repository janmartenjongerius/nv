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

nv/config\.WithEncoding\s+Env\s+\[.+?]

nv/config\.init\.1\s+Default\s+text
\s+Formats\s+\[.+?]

nv/debug\.init\.0\s+Callbacks\s+\[nv/config\.WithEncoding nv/config\.init\.1 nv/debug\.init\.0( nv\.init\.0 nv\.init\.1)?]`)

	if !want.MatchString(got) {
		t.Errorf("Could not match output against %q:\n%v", want, got)
	}
}
