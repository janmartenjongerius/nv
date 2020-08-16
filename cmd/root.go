package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"janmarten.name/nv/config"
	"strings"
)

var (
	format  = config.DefaultEncoding
	rootCmd = &cobra.Command{
		Use:              "nv",
		Short:            "List environment variables",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			format, _ = cmd.Root().PersistentFlags().GetString("format")

			if len(format) == 0 {
				format = config.DefaultEncoding
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			config.NewExporter(format, cmd.OutOrStdout()).Export(config.Environment...)
		},
	}
)

// -- string pflag.Value
type formatValue string

// Set the value for the format flag.
func (v *formatValue) Set(val string) error {
	if !config.HasEncoding(val) {
		return fmt.Errorf("unknown encoding %q", val)
	}

	*v = formatValue(val)
	return nil
}

// Type gets the type of the format flag.
func (v *formatValue) Type() string {
	return "string"
}

// String gets the format flag value as a string.
func (v *formatValue) String() string {
	return string(*v)
}

// Execute runs the root command.
func Execute(version string) error {
	rootCmd.Version = version

	// Since encodings can be registered with plugins, the format flag has to be registered at runtime.
	rootCmd.PersistentFlags().StringP(
		"format",
		"f",
		config.DefaultEncoding,
		"One of: "+strings.Join(config.GetEncodings(), ", "))

	// Ensure only known formats are allowed.
	rootCmd.Flag("format").Value = new(formatValue)

	return rootCmd.Execute()
}
