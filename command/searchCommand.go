package command

import (
	"flag"
	"fmt"
	"janmarten.name/env/config"
	"janmarten.name/env/search"
	"time"
)

type searchCommand struct {
	BaseCommand
	engine      search.Engine
	suggestions int
	interactive bool
}

func (cmd searchCommand) Run(args []string, io IO) int {
	var line string

	for {
		if cmd.interactive == false && len(args) == 0 {
			break
		}

		io.Write("$ ")

		if len(args) > 0 {
			line, args = args[0], args[1:]
			io.WriteLn(line)
		} else {
			line = io.ReadLn()
		}

		if len(line) == 0 {
			if len(args) > 0 {
				continue
			}

			break
		}

		r := cmd.engine.Query(line, uint(cmd.suggestions)).Result()

		if r.Match == nil && len(r.Suggestions) == 1 {
			r = cmd.engine.Query(r.Suggestions[0], 0).Result()
		}

		if r.Match == nil {
			io.WriteError(
				fmt.Sprintf("Could not find %s.\n", r.Request.Query))
		}

		if r.Match == nil || len(r.Suggestions) > 1 {
			if len(r.Suggestions) > 0 {
				suggestion := "Suggestions:\n"

				for _, s := range r.Suggestions {
					suggestion += fmt.Sprintf("  - %s\n", s)
				}

				io.WriteError(suggestion)
				time.Sleep(time.Millisecond * 100)
			}

			continue
		}

		v := (*r.Match).(*config.Variable)
		io.Write(fmt.Sprintf("%s=%+v\n", v.Key, v.Value))
	}

	return 0
}

func NewSearchCommand(engine search.Engine) Command {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	cmd := searchCommand{
		BaseCommand: BaseCommand{
			flags:       flags,
			description: "Fuzzy search for environment variables.",
		},
		engine:      engine,
		suggestions: 5,
		interactive: false,
	}

	flags.BoolVar(
		&cmd.interactive,
		"interactive",
		cmd.interactive,
		"Interactively search for variables",
	)

	flags.IntVar(
		&cmd.suggestions,
		"num-suggestions",
		cmd.suggestions,
		"Set the number of suggestions returned when an entry could not be found",
	)

	return &cmd
}
