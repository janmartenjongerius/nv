package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	wd, _ := os.Getwd()
	want := make([]string, 0)

	// The init func runs implicitly. We can directly test its output.
	for _, location := range PluginLocations {
		if !filepath.IsAbs(location) {
			location = path.Join(wd, location)
		}

		for _, extension := range PluginExtensions {
			if matches, err := filepath.Glob(fmt.Sprintf("%s/*%s", location, extension)); err == nil {
				want = append(want, matches...)
			}
		}
	}

	got := regexp.MustCompile(
		"^(" +
			strings.Replace(
				strings.Join(PluginsLoaded, "|"),
				"\\",
				"\\\\",
				-1,
			) +
			")$",
	)

	for _, f := range want {
		t.Run(f, func(t *testing.T) {
			if !got.MatchString(f) {
				t.Errorf(
					"Could not match %q against %q. Available %q",
					f,
					got.String(),
					PluginsLoaded)
			}
		})
	}
}
