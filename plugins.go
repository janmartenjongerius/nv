package main

import (
	"os"
	"path"
	"path/filepath"
	"plugin"
)

var (
	// PluginLocations tells the application where to look for plugin files.
	PluginLocations = []string{
		"/usr/local/lib/nv",
		"/usr/lib/nv",
		"plugins",
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
				PluginsLoaded = append(PluginsLoaded, path)

				break
			}
		}

		return err
	}

	// PluginsLoaded is a list of loaded plugin files.
	PluginsLoaded []string
)

func init() {
	wd, _ := os.Getwd()

	PluginsLoaded = make([]string, 0)

	for _, location := range PluginLocations {
		if !filepath.IsAbs(location) {
			location = path.Join(wd, location)
		}

		_ = filepath.Walk(location, pluginLoader)
	}
}
