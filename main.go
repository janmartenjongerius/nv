/*
Env

More than ever, modern development relies on environment variables.
To easily debug the local environment or export it, a chain of commands specific to your operating system would do.
However, env wants to solve this in a modern way, cross platform.

Features
	* Get an environment variable, ensuring the environment variable exists.
	* Search for environment variables, interactively and programmatically.
	* Export a list of required environment variables to a DotEnv file.
	* Set, update and unset environment variables programmatically.
*/
package main

import (
	"context"
	"janmarten.name/env/command"
	"janmarten.name/env/config"
	_ "janmarten.name/env/config/encoding/json"
	_ "janmarten.name/env/config/encoding/text"
	"janmarten.name/env/config/export"
	"janmarten.name/env/search"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

var (
	// PluginLocations tells the application where to look for plugin files.
	PluginLocations = []string{"/usr/local/lib/env"}

	// PluginExtensions tells the plugin loader which extensions to match while loading plugins.
	PluginExtensions = []string{".so"}

	// Conditionally load the plugin that matches the given path.
	pluginLoader filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		for _, extension := range PluginExtensions {
			if err != nil {
				break
			}

			if strings.HasSuffix(info.Name(), extension) {
				_, err = plugin.Open(path)
				break
			}
		}

		return err
	}
)

func init() {
	for _, location := range PluginLocations {
		_ = filepath.Walk(location, pluginLoader)
	}
}

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
