package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/dikaeinstein/godl/pkg/text"
)

const (
	ShellBash = "bash"
	ShellZsh  = "zsh"
	ShellFish = "fish"
)

// newCompletionCmd returns the a new instance of the completion command
func newCompletionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [" + ShellBash + "|" + ShellZsh + "|" + ShellFish + "]",
		Short: "Generate completion script.",
		Example: text.Indent(heredoc.Docf(example(),
			text.Bold("Bash"), text.Bold("Zsh"), text.Bold("Fish")), "  "),
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{ShellBash, ShellZsh, ShellFish},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	return completionCmd
}

func run(cmd *cobra.Command, args []string) error {
	rootCmd := cmd.Parent()
	shell := args[0]
	out := rootCmd.OutOrStdout()

	switch shell {
	case ShellBash:
		return rootCmd.GenBashCompletionV2(out, true)
	case ShellZsh:
		return rootCmd.GenZshCompletion(out)
	case ShellFish:
		return rootCmd.GenFishCompletion(out, true)
	default:
		return fmt.Errorf("unsupported shell: %q", shell)
	}
}

func example() string {
	return `
		%s:

		# Generate completion script
		$ godl completion bash > godl-completion.bash

		# Install the completion script
		$ sudo cp godl-completion.bash /etc/bash_completion.d/

		# Reload the shell
		$ source ~/.bashrc

		%s:

		# Generate completion script
		$ godl completion zsh > _godl

		# Install the completion script
		$ sudo cp _godl /usr/local/share/zsh/site-functions/

		# Reload the shell
		$ source ~/.zshrc

		%s:

		# Generate and install completion script
		$ godl completion fish > ~/.config/fish/completions/godl.fish

		# Reload the shell
		$ source ~/.config/fish/config.fish
	`
}
