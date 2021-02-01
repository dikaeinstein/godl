// Copyright © 2019 Onyedikachi Solomon Okwa <solozyokwa@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package completion

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// New returns the completion command
func New(godl *godl.GodlCmd) *cobra.Command {
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

			completion := completionCmd{
				bashSymlinkDir: bashSymlinkDir,
				fsys:           fsys,
				homeDir:        home,
				rootCmd:        godl,
				zshSymlinkDir:  zshSymlinkDir,
			}

			return completion.Run(args[0])
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("provide shell to configure e.g bash or zsh")
			}
			return nil
		},
	}
}

// osLinker is an os based filesystem that can symlink files
type osLinkerFS struct{}

func (osLinkerFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (osLinkerFS) Symlink(oldName, newName string) error {
	return os.Symlink(oldName, newName)
}

type completionCmd struct {
	bashSymlinkDir string
	fsys           fs.FS
	homeDir        string
	rootCmd        *godl.GodlCmd
	zshSymlinkDir  string
}

func (c *completionCmd) Run(shell string) error {
	autocompleteDir := path.Join(c.homeDir, ".godl", "autocomplete")
	bashDir := path.Join(autocompleteDir, "bash")
	zshDir := path.Join(autocompleteDir, "zsh")
	godlutil.Must(os.MkdirAll(bashDir, os.ModePerm))
	godlutil.Must(os.MkdirAll(zshDir, os.ModePerm))

	switch shell {
	case "bash":
		bashTarget := filepath.Join(bashDir, "godl")
		err := c.rootCmd.GenerateBashCompletionFile(bashTarget)
		if err != nil {
			return err
		}
		return fs.Symlink(c.fsys, bashTarget, filepath.Join(c.bashSymlinkDir, "godl"))
	case "zsh":
		zshTarget := filepath.Join(zshDir, "_godl")
		err := c.rootCmd.GenerateZshCompletionFile(zshTarget)
		if err != nil {
			return err
		}
		return fs.Symlink(c.fsys, zshTarget, filepath.Join(c.zshSymlinkDir, "_godl"))
	default:
		return errors.New("unknown shell passed")
	}
}
