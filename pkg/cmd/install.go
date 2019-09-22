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
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(install)
}

// installCmd represents the install command
var install = &cobra.Command{
	Use:   "install version",
	Short: "Installs the specified go binary version from local or remote.",
	Long: `Installs the specified go binary version from local or remote.
It fetches the version from the remote if not found locally before installing it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fcr := fsFileCreatorRenamer{}
		gz := gzipUnArchiver{
			z: archiver.TarGz{
				Tar: &archiver.Tar{
					OverwriteExisting: true,
				},
				CompressionLevel: -1,
			},
		}

		dlDir, err := getDownloadDir()
		if err != nil {
			return err
		}

		dl := &goBinaryDownloader{
			baseURL:       "https://storage.googleapis.com/golang/",
			client:        &http.Client{},
			downloadDir:   dlDir,
			fCR:           fcr,
			forceDownload: forceDownload,
			genHash:       getBinaryHash,
			verifyHash:    verifyHash,
		}

		return installRelease(args[0], dlDir, gz, fsRemover{}, dl)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("provide binary archive version to install")
		}
		return nil
	},
}

func installRelease(version, dlDir string, gz unArchiver, dr dirRemover, dl *goBinaryDownloader) error {
	fmt.Println("Installing binary into /usr/local")
	gi := goInstaller{dlDir, gz, dr, dl}

	if err := gi.install(version); err != nil {
		return err
	}

	return nil
}

type goInstaller struct {
	dlDir string
	ua    unArchiver
	dr    dirRemover
	dl    *goBinaryDownloader
}

func (gi goInstaller) install(archiveVersion string) error {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := path.Join(gi.dlDir, archiveName)

	exists, err := versionExists(archiveVersion, gi.dlDir)
	if err != nil {
		return err
	}

	// download binary if it doesn't exist locally
	if !exists {
		fmt.Printf("%v not found locally.\n", archiveVersion)
		fmt.Println("fetching from remote...")
		if err := gi.dl.download(archiveVersion); err != nil {
			return fmt.Errorf("error downloading %v: %v", archiveVersion, err)
		}
	}

	// clean install - remove existing go installation before installing
	// new version
	fmt.Println("removing old installation...")
	err = removeGo(gi.dr)
	if err != nil {
		return fmt.Errorf("error removing old installation: %v", err)
	}
	fmt.Println("old installation removed")

	fmt.Printf("unpacking %v ...", archiveVersion)
	target := path.Join("/usr", "local")
	if err := gi.ua.Unarchive(downloadPath, target); err != nil {
		return err
	}

	fmt.Println("Installation successful. Type `go version` to check installation")
	return nil
}

type unArchiver interface {
	Unarchive(source, target string) error
}

type gzipUnArchiver struct {
	z archiver.TarGz
}

func (gz gzipUnArchiver) Unarchive(source, target string) error {
	return gz.z.Unarchive(source, target)
}

type dirRemover interface {
	RemoveAll(path string) error
}

type fsRemover struct{}

func (f fsRemover) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func removeGo(dr dirRemover) error {
	err := dr.RemoveAll(path.Join("/usr", "local", "go"))
	if err != nil {
		return err
	}
	return nil
}
