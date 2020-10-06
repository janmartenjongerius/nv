package main

import (
	"errors"
	"io/ioutil"
	"janmarten.name/nv/debug"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("OS is not expected to be able to load plugins: %v", runtime.GOOS)
	}

	want := make([]string, 0)
	got := regexp.MustCompile(
		"^(" +
			strings.Replace(
				strings.Join(pluginsLoaded, "|"),
				"\\",
				"\\\\",
				-1,
			) +
			")$",
	)

	for _, pattern := range pluginPatterns {
		if files, err := filepath.Glob(string(pattern)); err == nil {
			want = append(want, files...)
		}
	}

	for _, f := range want {
		t.Run(f, func(t *testing.T) {
			if !got.MatchString(f) {
				t.Errorf(
					"Could not match %q against %q. Available %q. Rejected %q",
					f,
					got.String(),
					pluginsLoaded,
					pluginsRejected)
			}
		})
	}

	if len(pluginsLoaded) > len(want) {
		t.Errorf(
			"More plugins loaded than expected. Want %v, got %v",
			want,
			got)
	}
}

func TestInit_debug(t *testing.T) {
	got := debug.Scope("nv.init.1").GetMessages()

	if !reflect.DeepEqual(pluginsLoaded, got["Loaded"]) {
		t.Errorf(
			"Unexpected loaded plugins. Want %v, got %v",
			pluginsLoaded,
			got["Loaded"])
	}

	if !reflect.DeepEqual(pluginsRejected, got["Rejected"]) {
		t.Errorf(
			"Unexpected rejected plugins. Want %v, got %v",
			pluginsRejected,
			got["Rejected"])
	}

	fileMap, ok := got["FileMap"].(map[pluginPattern][]string)

	if !ok {
		t.Fatalf(
			"Expected FileMap to be map[pluginPattern][]string: %q",
			fileMap)
	}

	found := 0

	for _, files := range fileMap {
		found += len(files)
	}

	if found != len(pluginsLoaded)+len(pluginsRejected) {
		t.Errorf(
			"Not all plugins were mapped. Want %q, got %v + %v",
			fileMap,
			pluginsLoaded,
			pluginsRejected)
	}
}

func TestPluginPattern_GetFiles(t *testing.T) {
	pattern := pluginPattern("[.so")
	got := make([]string, 0)
	want := got

	_, err := pattern.getFiles()

	if !errors.Is(err, filepath.ErrBadPattern) {
		t.Errorf(
			"Expected pattern %q to produce filepath.ErrBadPattern. Got %v",
			pattern,
			err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf(
			"Expected %v, got %v", want, got)
	}

	dir, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.RemoveAll(dir) }()

	pattern = pluginPattern(filepath.Join(dir, "*.so"))
	want = []string{
		filepath.Join(dir, "foo.so"),
		filepath.Join(dir, "bar.so"),
		filepath.Join(dir, "baz.so"),
	}

	for _, file := range want {
		if _, err := os.Create(file); err != nil {
			t.Fatal(err)
		}
	}

	got = pattern.mustGetFiles()

	sort.Strings(want)
	sort.Strings(got)

	if !reflect.DeepEqual(want, got) {
		t.Errorf(
			"Got unexpected files. Want %v, got %v", want, got)
	}
}

func TestPluginPattern_OpenPlugins(t *testing.T) {
	pattern := pluginPattern("[.so")
	_, _, err := pattern.openPlugins()

	if !errors.Is(err, filepath.ErrBadPattern) {
		t.Errorf(
			"Expected pattern %q to produce filepath.ErrBadPattern. Got %v",
			pattern,
			err)
	}

	dir, err := ioutil.TempDir(os.TempDir(), t.Name())

	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.RemoveAll(dir) }()

	if runtime.GOOS == "windows" {
		t.Skipf("OS is not expected to be able to load plugins: %v", runtime.GOOS)
	}

	dummy, err := ioutil.ReadFile(filepath.Join("plugins", "dummy.so"))

	if err != nil {
		t.Skipf("Could not open dummy plugin: %v", err)
	}

	pattern = pluginPattern(filepath.Join(dir, "*.so"))
	want := []string{
		filepath.Join(dir, "foo.so"),
		filepath.Join(dir, "bar.so"),
		filepath.Join(dir, "baz.so"),
	}

	for _, file := range want {
		fh, err := os.Create(file)

		if err != nil {
			t.Fatal(err)
		}

		_, err = fh.Write(dummy)

		if err != nil {
			t.Fatal(err)
		}
	}

	_, rejected := pluginLoader{pattern}.mustLoad()

	// Because the test plugins inherit from the dummy, which should have
	// already been loaded, the test will assert that the error is exactly such.
	if len(rejected) != len(want) {
		t.Fatalf(
			"Expected all plugins to be rejected. Want %v, got %v",
			want,
			rejected)
	}

	for _, reject := range rejected {
		if !strings.HasSuffix(reject.Error(), "plugin already loaded") {
			t.Errorf("Expected plugin to be already loaded: %v", reject)
		}
	}
}

func TestPluginLoader_MustLoad(t *testing.T) {
	loader := pluginLoader{"[.so"}

	defer func() {
		// Try to recover from the panic state.
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	loader.mustLoad()
}
