package command

import (
	"flag"
	"fmt"
	"janmarten.name/env/config"
	"strings"
)

type getCommand struct {
	BaseCommand
	cfg *config.Config
}

func (cmd getCommand) Run(args []string, io IO) int {
	if len(args) == 0 {
		return 1
	}

	var missing []string

	prependName := len(args) > 1

	for _, a := range args {
		if v, e := cmd.cfg.Variable(a); e == nil {
			io.WriteLn(cmd.formatResult(v, prependName))
			continue
		}

		missing = append(missing, a)
	}

	if len(missing) > 0 {
		io.WriteLn(
			fmt.Sprintf(
				"Could not find %v.\n",
				strings.Join(missing, ", ")))
		return 2
	}

	return 0
}

func (cmd getCommand) formatResult(v *config.Variable, prependName bool) string {
	result := ""

	if prependName {
		result += fmt.Sprintf("%v=", v.Key)
	}

	result += fmt.Sprintf("%v", v.Value)

	return result
}

func NewGetCommand(cfg *config.Config) Command {
	return &getCommand{
		BaseCommand: BaseCommand{
			flags: flag.NewFlagSet("get", flag.ContinueOnError),
		},
		cfg: cfg,
	}
}
