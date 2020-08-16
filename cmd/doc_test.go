package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDocCmd_PreRunE(t *testing.T) {
	cmd := &cobra.Command{}
	args := make([]string, 0)
	outputDir := docCmd.Flag("output-dir").Value
	tmp := os.TempDir()

	cases := []struct {
		in  string
		err bool
	}{
		{
			in:  filepath.Join(tmp, "nv-test-foo"),
			err: false,
		},
		{
			in:  "",
			err: true,
		},
	}

	// Clean up after the test ends.
	defer func() {
		for _, c := range cases {
			if _, err := os.Stat(c.in); !os.IsNotExist(err) {
				_ = os.RemoveAll(c.in)
			}
		}
	}()

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			err := outputDir.Set(c.in)

			if err != nil {
				t.Errorf("Could not set output dir to %q: %v", c.in, err)
			}

			err = docCmd.PreRunE(cmd, args)

			if !c.err {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if _, err := os.Stat(c.in); os.IsNotExist(err) {
					t.Errorf("Failed creating output directory %q: %v", c.in, err)
				}
			}
		})
	}
}

func TestDocCmd_RunE(t *testing.T) {
	cmd := &cobra.Command{Use: "nv-test"}
	dir := filepath.Join(os.TempDir(), "nv-test")

	// Set the output directory to a safe location.
	_ = docCmd.Flag("output-dir").Value.Set(dir)

	// Clean up after the test ends.
	defer func() {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			_ = os.RemoveAll(dir)
		}
	}()

	// Ensure the directory is made.
	_ = docCmd.PreRunE(cmd, []string{})

	cases := []struct {
		in  string
		err bool
	}{
		{
			in:  "man",
			err: false,
		},
		{
			in:  "markdown",
			err: false,
		},
		{
			in:  "rst",
			err: false,
		},
		{
			in:  "yaml",
			err: false,
		},
		{
			in:  "foo",
			err: true,
		},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			args := []string{c.in}
			err := docCmd.RunE(cmd, args)

			if reflect.DeepEqual(err, nil) == c.err {
				t.Errorf("Expected error? %v, got: %v", c.err, err)
			}
		})
	}
}
