package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
)

var (
	rootCmd = &cobra.Command{Use: "nv"}
)

// Executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	tag, err := exec.Command("git", "describe", "--tags").Output()

	if err != nil {
		log.Println(err)
		rootCmd.Version = "Unknown"
		return
	}

	rootCmd.Version = strings.TrimSpace(string(tag))
}
