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

package install

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

type Archiver interface {
	Unarchive(source, target string) error
}

type Install struct {
	Archiver Archiver
	Dl       *downloader.Downloader
	Timeout  time.Duration
}

func (i *Install) Run(ctx context.Context, version string) error {
	archiveName := fmt.Sprintf("%s%s.%s", downloader.Prefix(), version, downloader.Postfix())
	downloadPath := path.Join(i.Dl.DownloadDir, archiveName)

	exists, err := gv.VersionExists(version, i.Dl.DownloadDir)
	if err != nil {
		return err
	}

	// download binary if it doesn't exist locally or the -forceDownload flag is passed
	if !exists || i.Dl.ForceDownload {
		fmt.Printf("%v not found locally.\n", version)
		fmt.Println("fetching from remote...")

		ctx, cancel := context.WithTimeout(ctx, i.Timeout)
		defer cancel()
		err = i.Dl.Download(ctx, version)
		if err != nil {
			return fmt.Errorf("error downloading %v: %v", version, err)
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

	fmt.Printf("unpacking %v ...\n", version)
	target := path.Join("/usr", "local")
	err = i.Archiver.Unarchive(downloadPath, target)
	if err != nil {
		return err
	}

	fmt.Println("adding to $PATH...")
	pathsD := path.Join("/etc", "paths.d", "go")
	err = fsys.WriteFile(i.Dl.FS, pathsD, []byte("/usr/local/go/bin\n"), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Installation successful. Type `go version` to check installation")
	return nil
}
