/*
TODO:
	> Write documentation
*/
package main

import (
	"context"
	"janmarten.name/env/command"
	"janmarten.name/env/config"
	"janmarten.name/env/config/export"
	"janmarten.name/env/search"
	"os"
	"runtime"
)

func main() {
	cfg := config.Parse(os.Environ(), "=")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, search.KeyParallel, runtime.GOMAXPROCS(0))
	defer cancel()

	exportFactory := export.NewFactory(os.Stdout)
	engine := search.New(ctx, cfg.MapInterfaces())

	app := command.NewApplication(ctx, os.Stdin, os.Stdout, os.Stderr)
	app.SetDescription("Look up and manage environment variables.")
	app.Register(command.NewGetCommand(engine, exportFactory))
	app.Register(command.NewSearchCommand(engine))

	os.Exit(
		app.Run(
			command.ExtractCommandArgs(os.Args[1:]),
		),
	)
}
