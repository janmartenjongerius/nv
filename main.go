/*
Nv

Modern development relies on environment variables to keep track of variances between environments.
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
	"fmt"
	"janmarten.name/nv/cmd"
	"janmarten.name/nv/debug"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

var version = ""

func isDevMode() bool {
	return version == "" || strings.HasPrefix(version, "dev")
}

func init() {
	if len(version) == 0 {
		tag, _ := exec.Command("git", "describe", "--tags").Output()
		version = fmt.Sprintf("dev-%s", strings.TrimSpace(string(tag)))
	}

	debug.RegisterCallback(func() debug.Messages {
		wd, _ := os.Getwd()
		host, _ := os.Hostname()
		usr, _ := user.Current()

		return debug.Messages{
			"Version":           version,
			"Platform":          fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			"Runtime":           runtime.Version(),
			"Compiler":          runtime.Compiler,
			"Threads":           uint8(runtime.GOMAXPROCS(0)),
			"Working directory": wd,
			"Hostname":          host,
			"User":              usr.Username,
			"Args":              os.Args,
			"DevMode":           isDevMode(),
		}
	})
}

func main() {
	if err := cmd.Execute(version); err != nil {
		// Cobra already performs error handling, so nothing much for us to do here.
		os.Exit(1)
	}
}
