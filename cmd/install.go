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
	"path"

	"github.com/mholt/archiver"

	"github.com/spf13/cobra"
)

type unArchiver interface {
	Unarchive(source, target string) error
}

type gzipUnArchiver struct {
	z archiver.TarGz
}

func (gz gzipUnArchiver) Unarchive(source, target string) error {
	return gz.z.Unarchive(source, target)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Installs the specified go binary archive version or path into /usr/local/go",
	RunE: func(cmd *cobra.Command, args []string) error {
		gz := gzipUnArchiver{
			z: archiver.TarGz{
				Tar: &archiver.Tar{
					OverwriteExisting: true,
				},
				CompressionLevel: -1,
			},
		}
		fmt.Println("Installing binary into /usr/local")
		err := installGoBinary(args[0], gz)
		if err != nil {
			return err
		}
		fmt.Println("Installation Complete. Type `go version` to check installation")
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("provide binary archive version to install")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installGoBinary(archiveVersion string, ua unArchiver) error {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	godlDownloadDir, err := getDownloadDir()

	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := path.Join(godlDownloadDir, archiveName)

	target := path.Join("/usr", "local")
	return ua.Unarchive(downloadPath, target)
}
