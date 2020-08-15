package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(nv completion bash)

# To load completions for each session, execute once:
Linux:
  $ nv completion bash > /etc/bash_completion.d/nv
MacOS:
  $ nv completion bash > /usr/local/etc/bash_completion.d/nv

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ nv completion zsh > "${fpath[1]}/_nv"

# You will need to start a new shell for this setup to take effect.

Fish:

$ nv completion fish | source

# To load completions for each session, execute once:
$ nv completion fish > ~/.config/fish/completions/nv.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
		}
		return fmt.Errorf("unknown shell %q requested", args[0])
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
