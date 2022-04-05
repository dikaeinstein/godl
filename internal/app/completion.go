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

package app

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

// Generator controls how the completion files should be generated
type CompletionGenerator interface {
	GenerateBashCompletion(io.Writer) error
	GenerateFishCompletion(out io.Writer, includeDesc bool) error
	GenerateZshCompletion(io.Writer) error
}

// Completion generates completion files for zsh/bash shell.
type Completion struct {
	AutocompleteDir string
	BashSymlinkDir  string
	FishSymlinkDir  string
	FS              fs.FS
	CompletionGenerator
	HomeDir       string
	ZshSymlinkDir string
}

// Run generates completion file for specified shell
func (c *Completion) Run(shell string, out io.Writer, useDefault bool) error {
	switch shell {
	case "bash":
		if err := c.GenerateBashCompletion(out); err != nil {
			return err
		}

		if useDefault {
			bashTarget := CompletionMakeTarget(shell, c.AutocompleteDir)
			return fsys.Symlink(c.FS, bashTarget, filepath.Join(c.BashSymlinkDir, "godl"))
		}

		return nil
	case "zsh":
		if err := c.GenerateZshCompletion(out); err != nil {
			return err
		}

		if useDefault {
			zshTarget := CompletionMakeTarget(shell, c.AutocompleteDir)
			return fsys.Symlink(c.FS, zshTarget, filepath.Join(c.ZshSymlinkDir, "_godl"))
		}
		return nil
	case "fish":
		if err := c.GenerateFishCompletion(out, true); err != nil {
			return err
		}

		if useDefault {
			zshTarget := CompletionMakeTarget(shell, c.AutocompleteDir)
			return fsys.Symlink(c.FS, zshTarget, filepath.Join(c.FishSymlinkDir, "godl.fish"))
		}
		return nil
	default:
		return errors.New("unknown shell passed")
	}
}

// CompletionMakeTarget creates the file and it's parent directories where the
// completion output can be written to.
func CompletionMakeTarget(shell, autocompleteDir string) string {
	bashDir := filepath.Join(autocompleteDir, shell)
	godlutil.Must(os.MkdirAll(bashDir, os.ModePerm))
	return filepath.Join(bashDir, "godl")
}
