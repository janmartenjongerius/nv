package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"janmarten.name/nv/debug"
	"text/tabwriter"
)

var (
	scopeWalkerFactory = func(writer io.Writer) debug.ScopeWalker {
		return func(scope debug.Scope) debug.MessageWalker {
			_, _ = fmt.Fprintln(writer, "")
			_, _ = fmt.Fprint(writer, scope)

			return messageWalkerFactory(writer)
		}
	}

	messageWalkerFactory = func(writer io.Writer) debug.MessageWalker {
		return func(group string, message interface{}) {
			_, _ = fmt.Fprintf(writer, "\t%v\t%v\n", group, message)
		}
	}

	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Show information of use when debugging Nv.",
		RunE: func(cmd *cobra.Command, args []string) error {
			writer := new(tabwriter.Writer)
			writer.Init(
				cmd.OutOrStdout(),
				24,
				8,
				2,
				[]byte(" ")[0],
				tabwriter.DiscardEmptyColumns)

			cmd.Println("Nv debug:")
			debug.Walk(scopeWalkerFactory(writer))

			return writer.Flush()
		},
	}
)

func init() {
	rootCmd.AddCommand(debugCmd)
}
