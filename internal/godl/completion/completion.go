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

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fs"
)

// Generator controls how the completion files should be generated
type Generator interface {
	GenerateBashCompletionFile(string) error
	GenerateZshCompletionFile(string) error
}

// Completion generates completion files for zsh/bash shell.
type Completion struct {
	BashSymlinkDir string
	Generator
	FSys          fs.FS
	HomeDir       string
	ZshSymlinkDir string
}

// Run generates completion file for specified shell
func (c *Completion) Run(shell string) error {
	autocompleteDir := path.Join(c.HomeDir, ".godl", "autocomplete")
	bashDir := path.Join(autocompleteDir, "bash")
	zshDir := path.Join(autocompleteDir, "zsh")
	godlutil.Must(os.MkdirAll(bashDir, os.ModePerm))
	godlutil.Must(os.MkdirAll(zshDir, os.ModePerm))

	switch shell {
	case "bash":
		bashTarget := filepath.Join(bashDir, "godl")
		err := c.GenerateBashCompletionFile(bashTarget)
		if err != nil {
			return err
		}
		return fs.Symlink(c.FSys, bashTarget, filepath.Join(c.BashSymlinkDir, "godl"))
	case "zsh":
		zshTarget := filepath.Join(zshDir, "_godl")
		err := c.GenerateZshCompletionFile(zshTarget)
		if err != nil {
			return err
		}
		return fs.Symlink(c.FSys, zshTarget, filepath.Join(c.ZshSymlinkDir, "_godl"))
	default:
		return errors.New("unknown shell passed")
	}
}
