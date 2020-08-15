package main

import (
	"os"
	"path"
	"path/filepath"
	"plugin"
	"runtime"
)

var (
	// PluginLocations tells the application where to look for plugin files.
	PluginLocations = map[string]map[string][]string{
		"linux": {
			"amd64": {
				"/usr/local/lib/nv",
				"/usr/lib/nv",
				"plugins",
			},
		},
	}

	// PluginExtensions tells the plugin loader which extensions to match while loading plugins.
	PluginExtensions = []string{".so"}

	// Conditionally load the plugin that matches the given path.
	pluginLoader filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		for _, extension := range PluginExtensions {
			if !info.IsDir() && filepath.Ext(info.Name()) == extension {
				_, err = plugin.Open(path)
				break
			}
		}

		return err
	}
)

func init() {
	wd, _ := os.Getwd()

	for _, location := range PluginLocations[runtime.GOOS][runtime.GOARCH] {
		if !filepath.IsAbs(location) {
			location = path.Join(wd, location)
		}

		_ = filepath.Walk(location, pluginLoader)
	}
}
