package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var (
	docOutputDir string

	docCmd = &cobra.Command{
		Use:       "doc [man|markdown|rst|yaml]",
		Short:     "Generate documentation",
		ValidArgs: []string{"man", "markdown", "rst", "yaml"},
		Args:      cobra.ExactValidArgs(1),
		Aliases:   []string{"docs", "documentation"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if docOutputDir == "" {
				return fmt.Errorf("no output directory provided")
			}

			return os.MkdirAll(docOutputDir, os.ModeDir)
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
			return nil
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
