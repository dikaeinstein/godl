package cli

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/dikaeinstein/godl/internal/app/completion"
	"github.com/dikaeinstein/godl/pkg/exitcode"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/text"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// NewCompletionCmd returns the a new instance of the completion command
func NewCompletionCmd(godl completion.Generator) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate completion script.",
		Example: text.Indent(heredoc.Docf(example(),
			text.Bold("Bash"), text.Bold("Zsh"), text.Bold("Fish")), "  "),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return exitcode.NewError(errors.New("provide shell to configure e.g bash or zsh"), 1)
			}
			return nil
		},
		ValidArgs: []string{"bash", "zsh", "fish"},
	}

	useDefault := completionCmd.Flags().BoolP("default", "d", false, "Generate and load completion into default path based on shell")
	completionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		autocompleteDir := path.Join(home, ".godl", "autocomplete")

		var out io.Writer
		if *useDefault {
			outFile, err := os.Create(completion.MakeTarget(args[0], autocompleteDir))
			if err != nil {
				return err
			}
			defer outFile.Close()

			out = outFile
		} else {
			out = os.Stdout
		}

		bashSymlinkDir := path.Join("/usr", "local", "etc", "bash_completion.d")
		zshSymlinkDir := path.Join("/usr", "local", "share", "zsh", "site-functions")
		fishSymlinkDir := path.Join(home, ".config", "fish", "completions")

		c := completion.Completion{
			BashSymlinkDir:  bashSymlinkDir,
			FS:              fsys.OsFS{},
			FishSymlinkDir:  fishSymlinkDir,
			HomeDir:         home,
			Generator:       godl,
			ZshSymlinkDir:   zshSymlinkDir,
			AutocompleteDir: autocompleteDir,
		}

		return c.Run(args[0], out, *useDefault)
	}

	return completionCmd
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
