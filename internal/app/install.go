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
	"context"
	"fmt"
	"path"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/version"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/text"
)

type Archiver interface {
	Unarchive(source, target string) error
}

type Install struct {
	Archiver Archiver
	Dl       *downloader.Downloader
	Timeout  time.Duration
}

// Run installs the go version.
func (i *Install) Run(ctx context.Context, ver string) error {
	archiveName := fmt.Sprintf("%s%s.%s", downloader.Prefix(), ver, downloader.Postfix())
	downloadPath := path.Join(i.Dl.DownloadDir, archiveName)

	fmt.Println(text.Green("Installing binary into /usr/local"))

	exists, err := version.Exists(ver, i.Dl.DownloadDir)
	if err != nil {
		return err
	}

	// download binary if it doesn't exist locally or the -forceDownload flag is passed
	if !exists || i.Dl.ForceDownload {
		fmt.Printf("%v not found locally.\n", ver)
		fmt.Println("fetching from remote...")

		ctx, cancel := context.WithTimeout(ctx, i.Timeout)
		defer cancel()
		err = i.Dl.Download(ctx, ver)
		if err != nil {
			return fmt.Errorf("error downloading %v: %v", ver, err)
		}
	}

	// clean install - remove existing go installation before installing
	// new version
	fmt.Println()
	fmt.Println("removing old installation...")
	err = fsys.RemoveAll(i.Dl.FS, path.Join("/usr", "local", "go"))
	if err != nil {
		return fmt.Errorf("error removing old installation: %v", err)
	}
	fmt.Println("old installation removed")

	fmt.Printf("unpacking %v ...\n", ver)
	target := path.Join("/usr", "local")
	err = i.Archiver.Unarchive(downloadPath, target)
	if err != nil {
		return err
	}

	fmt.Println("adding to $PATH...")
	pathsD := path.Join("/etc", "paths.d", "go")
	const perm = 0o644
	err = fsys.WriteFile(i.Dl.FS, pathsD, []byte("/usr/local/go/bin\n"), perm)
	if err != nil {
		return err
	}

	fmt.Println(text.Green("Installation successful. Type `go version` to check installation"))
	return nil
}
