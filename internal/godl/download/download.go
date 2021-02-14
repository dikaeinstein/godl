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

package download

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fs/os"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/spf13/cobra"
)

var forceDownload bool
var timeout int64

// New returns the download command
func New() *cobra.Command {
	download := &cobra.Command{
		Use:   "download version",
		Short: "Download go binary archive.",
		Long: `Download the archive version from https://golang.org/dl/ and save to $HOME/godl/downloads.

	By default, if archive version already exists locally, godl doesn't attempt to download it again.
	To force it to download the version again pass the --force flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dlDir, err := godlutil.GetDownloadDir()
			if err != nil {
				return err
			}

			dl := &downloader.Downloader{
				BaseURL:       "https://storage.googleapis.com/golang/",
				Client:        &http.Client{},
				DownloadDir:   dlDir,
				Fsys:          os.FS{},
				ForceDownload: forceDownload,
				Hasher:        hash.NewRemoteHasher(http.DefaultClient),
				HashVerifier:  godlutil.VerifyHash,
			}

			d := downloadCmd{dl}
			return d.Run(cmd.Context(), args[0])
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("provide binary archive version to download")
			}
			return nil
		},
	}

	download.Flags().BoolVarP(&forceDownload, "force", "f", false,
		"Force download instead of using local version")
	download.Flags().Int64VarP(&timeout, "timeout", "t", 60000, "Set the download timeout.")

	return download
}

type downloadCmd struct {
	dl *downloader.Downloader
}

func (d *downloadCmd) Run(ctx context.Context, archiveVersion string) error {
	fmt.Printf("Downloading go binary %v\n", archiveVersion)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout*int64(time.Millisecond)))
	defer cancel()
	err := d.dl.Download(ctx, archiveVersion)
	if err != nil {
		return fmt.Errorf("error downloading %v: %v", archiveVersion, err)
	}

	fmt.Println("\nDownload complete")
	return nil
}
