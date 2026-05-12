// Copyright © 2019 Onyedikachi Solomon Okwa <onyedikachi.okwa@gmail.com>
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
	"context"
	"fmt"
	"io/fs"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/dikaeinstein/godl/internal/godlutil"
	"github.com/dikaeinstein/godl/internal/version"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/text"
)

type Archiver interface {
	Unarchive(ctx context.Context, source, target string) error
}

type Install struct {
	Archiver    Archiver
	Dl          Downloader
	Timeout     time.Duration
	FS          fs.FS
	DownloadDir string
}

func (i *Install) Configure(timeout time.Duration) {
	i.Timeout = timeout
}

// Run installs the go version.
func (i *Install) Run(
	ctx context.Context,
	ver, os, arch string,
	forceDownload bool,
) error {
	archiveName := godlutil.ArchiveName(ver, os, arch)
	downloadPath := path.Join(i.DownloadDir, archiveName)

	fmt.Println("Installing binary into /usr/local/bin...")

	exists, err := version.Exists(archiveName, i.DownloadDir)
	if err != nil {
		return err
	}

	// download binary if it doesn't exist locally or the -forceDownload flag is passed
	if !exists || forceDownload {
		fmt.Printf("%v not found locally.\n", ver)
		fmt.Println("fetching from remote...")

		ctx, cancel := context.WithTimeout(ctx, i.Timeout)
		defer cancel()
		err = i.Dl.Download(ctx, ver, os, arch)
		if err != nil {
			return fmt.Errorf("error downloading %v: %v", ver, err)
		}
	}

	fmt.Printf("unpacking %v ...\n", archiveName)
	target, err := installDir(ver)
	if err != nil {
		return fmt.Errorf("error getting install directory: %v", err)
	}

	err = i.Archiver.Unarchive(ctx, downloadPath, target)
	if err != nil {
		return err
	}

	goBinDir := path.Join(target, "go", "bin")
	if err := fsys.SymlinkDir(i.FS, goBinDir, "/usr/local/bin/"); err != nil {
		return err
	}

	fmt.Println(text.Green("Installation successful."))
	fmt.Println("Make sure to add /usr/local/bin to your PATH if it's not already there.")
	fmt.Println("You may need to restart your terminal for changes to take effect.")
	fmt.Println("Type `go version` to check installation")
	return nil
}

// installDir returns the `godl` installations directory
func installDir(v string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("home directory cannot be detected: %v", err)
	}

	return path.Join(home, ".godl", "installations", v), nil
}
