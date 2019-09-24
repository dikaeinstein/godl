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

package cmd

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
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
		bashSymDir := path.Join("/usr", "local", "etc", "bash_completion.d")
		zshSymDir := path.Join("/usr", "local", "share", "zsh", "site-functions")
		fl := fsLinker{}
		return completion(args[0], home, bashSymDir, zshSymDir, fl)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("provide shell to configure e.g bash or zsh")
		}
		return nil
	},
}

func completion(shell, home, bashSymDir, zshSymDir string, fl fileLinker) error {
	autocompleteDir := path.Join(home, ".godl", "autocomplete")
	bashDir := path.Join(autocompleteDir, "bash")
	zshDir := path.Join(autocompleteDir, "zsh")
	must(os.MkdirAll(bashDir, os.ModePerm))
	must(os.MkdirAll(zshDir, os.ModePerm))

	switch shell {
	case "bash":
		bashTarget := filepath.Join(bashDir, "godl")
		err := rootCmd.GenBashCompletionFile(bashTarget)
		if err != nil {
			return err
		}
		return fl.Symlink(bashTarget, filepath.Join(bashSymDir, "godl"))
	case "zsh":
		zshTarget := filepath.Join(zshDir, "_godl")
		err := rootCmd.GenZshCompletionFile(zshTarget)
		if err != nil {
			return err
		}
		return fl.Symlink(zshTarget, filepath.Join(zshSymDir, "_godl"))
	default:
		return errors.New("unknown shell passed")
	}
}

type fileLinker interface {
	Symlink(oldName, newName string) error
}

// fsLinker is a filesystem implementation of fileLinker
type fsLinker struct{}

func (fsLinker) Symlink(oldName, newName string) error {
	return os.Symlink(oldName, newName)
}
