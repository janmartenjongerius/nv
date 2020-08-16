package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
	"regexp"
	"runtime"
)

var (
	docOutputDir string

	// The following environments do not support the chown call:
	//
	// runtime.GOOS: "windows"
	// runtime.GOOS: "plan9"
	//
	// See https://golang.org/pkg/os/#Chown
	supportsChown = !regexp.
		MustCompile("^(windows|plan9)$").
		MatchString(runtime.GOOS)

	docCmd = &cobra.Command{
		Use:       "doc [man|markdown|rst|yaml]",
		Short:     "Generate documentation",
		ValidArgs: []string{"man", "markdown", "rst", "yaml"},
		Args:      cobra.ExactValidArgs(1),
		Aliases:   []string{"docs", "documentation"},
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if docOutputDir == "" {
				return fmt.Errorf("no output directory provided")
			}

			// Ensure the directory exists.
			err = os.MkdirAll(docOutputDir, os.ModeDir | 0755)

			// Ensure the directory belongs to the current user.
			// This is required if the ownership is inherited from a parent
			// directory and the current user would not be allowed to create new
			// files because of that inheritance.
			if err == nil && supportsChown {
				err = os.Chown(docOutputDir, os.Getuid(), os.Getgid())
			}

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "man":
				return doc.GenManTree(cmd.Root(), nil, docOutputDir)
			case "markdown":
				return doc.GenMarkdownTree(cmd.Root(), docOutputDir)
			case "rst":
				return doc.GenReSTTree(cmd.Root(), docOutputDir)
			case "yaml":
				return doc.GenYamlTree(cmd.Root(), docOutputDir)
			}

			return fmt.Errorf("unexpected document type %q", args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.Flags().StringVarP(
		&docOutputDir,
		"output-dir",
		"o",
		"",
		"Directory in which to generate documentation (required)")
}
