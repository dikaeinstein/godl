package cli

import (
	"errors"
	"os"
	"path"

	"github.com/dikaeinstein/godl/internal/godl/completion"
	"github.com/dikaeinstein/godl/pkg/exitcode"
	"github.com/dikaeinstein/godl/pkg/fs"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// osLinker is an os based filesystem that can symlink files
type osLinkerFS struct{}

func (osLinkerFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (osLinkerFS) Symlink(oldName, newName string) error {
	return os.Symlink(oldName, newName)
}

// New returns the a new instance of the completion command
func NewCompletionCmd(godl completion.Generator) *cobra.Command {
	return &cobra.Command{
		Use:   "completion <bash|zsh>",
		Short: "Generates completion scripts for bash or zsh.",
		Long: `To generate completion script run

	# godl completion <bash|zsh>

	To configure your bash shell to load completions for each session add to your bashrc

	# ~/.bashrc or ~/.bash_profile
	. "/usr/local/etc/profile.d/bash_completion.sh"
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := homedir.Dir()
			if err != nil {
				return err
			}
			bashSymlinkDir := path.Join("/usr", "local", "etc", "bash_completion.d")
			zshSymlinkDir := path.Join("/usr", "local", "share", "zsh", "site-functions")
			fsys := osLinkerFS{}

			c := completion.Completion{
				BashSymlinkDir: bashSymlinkDir,
				FSys:           fsys,
				HomeDir:        home,
				Generator:      godl,
				ZshSymlinkDir:  zshSymlinkDir,
			}
			return exitcode.NewError(c.Run(args[0]), 1)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return exitcode.NewError(errors.New("provide shell to configure e.g bash or zsh"), 1)
			}
			return nil
		},
	}
}
