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
	"os"
	"os/exec"
	"strings"
)

var version = ""

func init() {
	if len(version) == 0 {
		tag, err := exec.Command("git", "describe", "--tags").Output()

		if err != nil {
			version = "dev-unknown"
			return
		}

		version = fmt.Sprintf("dev-%s", strings.TrimSpace(string(tag)))
	}
}

func main() {
	if err := cmd.Execute(version); err != nil {
		// Cobra already performs error handling, so nothing much for us to do here.
		os.Exit(1)
	}
}
