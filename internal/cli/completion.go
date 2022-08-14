package cli

import (
	"io"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/text"
)

// newCompletionCmd returns the a new instance of the completion command
func newCompletionCmd() *cobra.Command {
	cCli := completionCli{}

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate completion script.",
		Example: text.Indent(heredoc.Docf(example(),
			text.Bold("Bash"), text.Bold("Zsh"), text.Bold("Fish")), "  "),
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		PreRunE:   cCli.setupConfig,
		RunE:      cCli.run,
	}

	completionCmd.Flags().BoolP(
		"default",
		"d",
		false,
		"Generate and load completion into default path based on shell",
	)

	return completionCmd
}

type completionConfig struct{ useDefault bool }

type completionCli struct {
	completionConfig
}

func (cCli completionCli) run(cmd *cobra.Command, args []string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	autocompleteDir := path.Join(home, ".godl", "autocomplete")

	var out io.Writer
	if cCli.useDefault {
		outFile, err := os.Create(
			app.CompletionMakeTarget(args[0], autocompleteDir),
		)
		if err != nil {
			return err
		}
		defer outFile.Close()

		out = outFile
	} else {
		out = os.Stdout
	}

	bashSymlinkDir := path.Join("/usr", "local", "etc", "bash_completion.d")
	zshSymlinkDir := path.Join(
		"/usr", "local", "share", "zsh", "site-functions",
	)
	fishSymlinkDir := path.Join(home, ".config", "fish", "completions")

	c := app.Completion{
		BashSymlinkDir:      bashSymlinkDir,
		FS:                  fsys.OsFS{},
		FishSymlinkDir:      fishSymlinkDir,
		HomeDir:             home,
		CompletionGenerator: cmd,
		ZshSymlinkDir:       zshSymlinkDir,
		AutocompleteDir:     autocompleteDir,
	}

	return c.Run(args[0], out, cCli.useDefault)
}

func (cCli completionCli) setupConfig(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	cCli.completionConfig.useDefault = viper.GetBool("default")

	return nil
}

func example() string {
	return `
		%s:

		$ source <(godl completion bash)

		# To load completions for each session, execute once:
		$ godl completion bash > /usr/local/etc/bash_completion.d/godl

		%s:

		# If shell completion is not already enabled in your environment,
		# you will need to enable it.  You can execute the following once:

		$ echo "autoload -U compinit; compinit" >> ~/.zshrc

		# To load completions for each session, execute once:
		$ godl completion zsh > "/usr/local/share/zsh/site-functions/_godl"

		# You will need to start a new shell for this setup to take effect.

		%s:

		$ godl completion fish | source

		# To load completions for each session, execute once:
		$ godl completion fish > ~/.config/fish/completions/godl.fish

		If you want 'godl' to generate and load the completion, just pass the --default(-d) flag:

		$ godl completion -d [bash|zsh|fish]
	`
}
