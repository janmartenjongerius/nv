/*
Nv

More than ever, modern development relies on environment variables.
To easily debug the local environment or export it, a chain of commands specific to your operating system would do.
However, nv wants to solve this in a modern way, cross platform.

Features

	* Get an environment variable, ensuring the environment variable exists.
	* Search for environment variables, interactively and programmatically.
	* Export a list of required environment variables to a DotEnv file.
	* Set, update and unset environment variables programmatically.
*/
package main

import (
	"janmarten.name/nv/cmd"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

var (
	// PluginLocations tells the application where to look for plugin files.
	PluginLocations = map[string]map[string][]string{
		"linux": {
			"amd64": {
				"/usr/local/lib/nv",
				"/usr/lib/nv",
			},
		},
	}

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
	for _, location := range PluginLocations[runtime.GOOS][runtime.GOARCH] {
		_ = filepath.Walk(location, pluginLoader)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		// Cobra already performs error handling, so nothing much for us to do here.
		os.Exit(1)
	}
}
