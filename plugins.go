package main

import (
	"janmarten.name/nv/debug"
	"os"
	"path/filepath"
	"plugin"
)

var (
	pluginPatterns = pluginLoader{
		"/usr/local/lib/nv/*.so",
		"/usr/lib/nv/*.so",
	}
	pluginsRejected = make([]error, 0)
	pluginsLoaded   = make([]string, 0)
)

type pluginPattern string

type pluginLoader []pluginPattern

func (pattern pluginPattern) walk(walk func(file string)) error {
	plugins, patternErr := filepath.Glob(string(pattern))

	if patternErr == nil {
		for _, p := range plugins {
			walk(p)
		}
	}

	return patternErr
}

func (pattern pluginPattern) getFiles() (files []string, err error) {
	files = make([]string, 0)

	err = pattern.walk(func(file string) {
		files = append(files, file)
	})

	return files, err
}

func (pattern pluginPattern) mustGetFiles() (files []string) {
	files, _ = pattern.getFiles()
	return files
}

func (pattern pluginPattern) openPlugins() (loaded []string, rejected []error, err error) {
	loaded = make([]string, 0)
	rejected = make([]error, 0)

	err = pattern.walk(func(file string) {
		if _, err := plugin.Open(file); err != nil {
			rejected = append(rejected, err)
			return
		}

		loaded = append(loaded, file)
	})

	return loaded, rejected, err
}

func (loader pluginLoader) load() (loaded []string, rejected []error, err error) {
	loaded = make([]string, 0)
	rejected = make([]error, 0)

	for _, pattern := range loader {
		l, r, err := pattern.openPlugins()
		loaded = append(loaded, l...)
		rejected = append(rejected, r...)

		if err != nil {
			return nil, nil, err
		}
	}

	return loaded, rejected, err
}

func (loader pluginLoader) mustLoad() (loaded []string, rejected []error) {
	loaded, rejected, err := loader.load()

	if err != nil {
		panic(err)
	}

	return loaded, rejected
}

func init() {
	// If we're in developer mode, prepend the plugins directory in our working
	// directory to the list of plugin patterns.
	if isDevMode() {
		wd, _ := os.Getwd()
		pluginPatterns = append(
			pluginLoader{
				pluginPattern(
					filepath.Join(wd, "plugins", "*.so"),
				),
			},
			pluginPatterns...,
		)
	}

	pluginsLoaded, pluginsRejected = pluginPatterns.mustLoad()

	debug.RegisterCallback("Plugins", func() debug.Messages {
		return debug.Messages{
			"Loaded":   pluginsLoaded,
			"Rejected": pluginsRejected,
			"FileMap": func(patterns []pluginPattern) map[pluginPattern][]string {
				fileMap := make(map[pluginPattern][]string)

				for _, pattern := range patterns {
					fileMap[pattern] = pattern.mustGetFiles()
				}

				return fileMap
			}(pluginPatterns),
		}
	})
}
