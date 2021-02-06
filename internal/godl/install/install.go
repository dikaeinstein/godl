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
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fs"
	"github.com/dikaeinstein/godl/pkg/fs/os"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

var forceDownload bool

// New returns the install command
func New() *cobra.Command {
	install := &cobra.Command{
		Use:   "install version",
		Short: "Installs the specified go binary version from local or remote.",
		Long: `Installs the specified go binary version from local or remote.
	It fetches the version from the remote if not found locally before installing it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dlDir, err := godlutil.GetDownloadDir()
			if err != nil {
				return err
			}

			install := installCmd{
				archiver: &archiver.TarGz{
					Tar: &archiver.Tar{
						OverwriteExisting: true,
					},
					CompressionLevel: -1,
				},
				dl: &downloader.Downloader{
					BaseURL:       "https://storage.googleapis.com/golang/",
					Client:        &http.Client{},
					DownloadDir:   dlDir,
					Fsys:          os.FS{},
					ForceDownload: forceDownload,
					Hasher:        hash.NewRemoteHasher(http.DefaultClient),
					HashVerifier:  godlutil.VerifyHash,
				},
			}

			fmt.Println("Installing binary into /usr/local")
			return install.Run(args[0])
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("provide binary archive version to install")
			}
			return nil
		},
	}

	install.Flags().BoolVarP(&forceDownload, "force", "f", false,
		"Force download instead of using local version")

	return install
}

type Archiver interface {
	Unarchive(source, target string) error
}

type installCmd struct {
	archiver Archiver
	dl       *downloader.Downloader
}

func (i *installCmd) Run(archiveVersion string) error {
	archiveName := fmt.Sprintf("%s%s.%s", downloader.Prefix(), archiveVersion, downloader.Postfix())
	downloadPath := path.Join(i.dl.DownloadDir, archiveName)

	exists, err := godlutil.VersionExists(archiveVersion, i.dl.DownloadDir)
	if err != nil {
		return err
	}

	// download binary if it doesn't exist locally or the -forceDownload flag is passed
	if !exists || i.dl.ForceDownload {
		fmt.Printf("%v not found locally.\n", archiveVersion)
		fmt.Println("fetching from remote...")
		if err := i.dl.Download(archiveVersion); err != nil {
			return fmt.Errorf("error downloading %v: %v", archiveVersion, err)
		}
	}

	// clean install - remove existing go installation before installing
	// new version
	fmt.Println()
	fmt.Println("removing old installation...")
	err = fs.RemoveAll(i.dl.Fsys, path.Join("/usr", "local", "go"))
	if err != nil {
		return fmt.Errorf("error removing old installation: %v", err)
	}
	fmt.Println("old installation removed")

	fmt.Printf("unpacking %v ...\n", archiveVersion)
	target := path.Join("/usr", "local")
	if err := i.archiver.Unarchive(downloadPath, target); err != nil {
		return err
	}

	fmt.Println("Installation successful. Type `go version` to check installation")
	return nil
}
