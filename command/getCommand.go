package command

import (
	"flag"
	"fmt"
	"janmarten.name/env/config"
	"janmarten.name/env/config/export"
	"janmarten.name/env/search"
	"strings"
)

type getCommand struct {
	BaseCommand
	engine search.Engine
	output string
	exportFactory export.Factory
}

func (cmd getCommand) Run(args []string, io IO) int {
	if len(args) == 0 {
		return 1
	}

	missing := make([]string, 0)
	found := make([]*config.Variable, 0)
	exporter, _ := cmd.exportFactory(cmd.output)

	cmd.engine.QueryAll(args, 0)

	for _, r := range cmd.engine.Results() {
		if r.Match == nil {
			missing = append(missing, r.Request.Query)
			continue
		}

		if v, ok := (*r.Match).(*config.Variable); ok == true {
			found = append(found, v)
		}
	}

	_ = exporter.ExportList(found)

	if len(missing) > 0 {
		io.WriteErrorLn(
			fmt.Sprintf(
				"Could not find %v.",
				strings.Join(missing, ", "),
			),
		)
		return 2
	}

	return 0
}

func NewGetCommand(engine search.Engine, exportFactory export.Factory) Command {
	flags := flag.NewFlagSet("get", flag.ContinueOnError)
	cmd := &getCommand{
		BaseCommand: BaseCommand{
			flags: flags,
			description: "Get the value of the requested environment variables.",
		},
		engine: engine,
		output: "text",
		exportFactory: exportFactory,
	}

	flags.StringVar(
		&cmd.output,
		"output",
		cmd.output,
		fmt.Sprintf("The format to use for output. Can be: %v", strings.Join(export.GetFormats(), ", ")),
	)

	return cmd
}
