package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"janmarten.name/nv/config"
	"janmarten.name/nv/search"
	"time"
)

var (
	numSuggestions uint

	searchCmd = &cobra.Command{
		Use:     "search [query...]",
		Aliases: []string{"find"},
		Short:   "Search environment variables",
		Long:    "Find environment variables that closest resemble your query.",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			svc := search.NewService(config.Environment)
			svc.Suggestions = numSuggestions

			exp := config.NewExporter(format, cmd.OutOrStdout())

			matches := make(config.Variables, 0)
			misses := 0

			// Work around a bug in cobra.
			// See: https://github.com/spf13/cobra/pull/894
			var printErr = func(message string) {
				_, _ = fmt.Fprint(cmd.ErrOrStderr(), message)
			}

			for _, r := range svc.Search(args...) {
				if r.Match == nil && len(r.Suggestions) == 1 && svc.Suggestions > 1 {
					r = svc.Search(r.Suggestions[0])[0]
				}

				if r.Match == nil {
					cmd.PrintErr(
						fmt.Sprintf(
							"Could not find %s.\n",
							r.Request.Query,
						),
					)
					misses++
				}

				if r.Match == nil || len(r.Suggestions) > 1 {
					if len(r.Suggestions) > 0 {
						suggestion := "  Suggestions:\n"

						for _, s := range r.Suggestions {
							suggestion += fmt.Sprintf("   - %s\n", s)
						}

						printErr(suggestion)
						time.Sleep(time.Millisecond * 100)
					} else {
						printErr("  No suggestions.\n")
					}

					continue
				}

				matches = append(matches, r.Match)
			}

			exp.Export(matches...)

			if misses > 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("could not resolve %d queries", misses)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().UintVarP(
		&numSuggestions,
		"num-suggestions",
		"n",
		uint(5),
		"Provide the given number of suggestions when a variable could not be found")
}
