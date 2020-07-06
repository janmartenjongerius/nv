package main

import (
	"context"
	"janmarten.name/env/command"
	"janmarten.name/env/config"
	"janmarten.name/env/search"
	"os"
)

func main() {
	cfg := config.Parse(os.Environ(), "=")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app := command.NewApplication(ctx, os.Stdin, os.Stdout, os.Stderr)
	app.Register(command.NewGetCommand(cfg))
	app.Register(command.NewSearchCommand(
		search.NewEngine(ctx, cfg.MapInterfaces())))

	if len(os.Args) < 2 {
		app.Usage(false)
		os.Exit(0)
	}

	os.Exit(app.Run(os.Args[1], os.Args[2:]))
}
