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
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

type fileCreator interface {
	Create(name string) (io.WriteCloser, error)
}

type fileCreatorRenamer interface {
	fileCreator
	Rename(oldPath, newPath string) error
}

type fsFileCreatorRenamer struct{}

func (ac fsFileCreatorRenamer) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (ac fsFileCreatorRenamer) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

var forceDownload bool

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [version]",
	Short: "Download go binary archive",
	Long: `Download the archive version from https://golang.org/dl/ and save to $HOME/godl/downloads.

By default, if archive version already exists locally, godl doesn't attempt to download it again.
To force it to download the version again pass the --force flag`,
	RunE: func(cmd *cobra.Command, args []string) error {
		archiveVersion := args[0]
		fcr := fsFileCreatorRenamer{}

		goBinDownloader := goBinaryDownloader{
			Client:  &http.Client{},
			BaseURL: "https://dl.google.com/go/",
			fCR:     fcr,
		}

		fmt.Printf("Downloading go binary %v\n", archiveVersion)
		err := goBinDownloader.download(archiveVersion, forceDownload)
		if err != nil {
			return err
		}
		fmt.Println("Download complete")
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("provide binary archive version to download")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolVarP(&forceDownload, "force", "f", false, "Force download")
}

// writeCounter counts the number of bytes written to it.
type writeCounter struct {
	bytesWritten       uint64
	TotalExpectedBytes uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.bytesWritten += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *writeCounter) PrintProgress() {
	percentDownloaded := float64(wc.bytesWritten) / float64(wc.TotalExpectedBytes) * 100
	fmt.Printf("\rDownloading... %.0f%% complete", math.Round(percentDownloaded))
}

type goBinaryDownloader struct {
	Client  *http.Client
	BaseURL string
	fCR     fileCreatorRenamer
}

func (goBinDown *goBinaryDownloader) download(archiveVersion string, forceDownload bool) error {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	godlDownloadDir, err := getDownloadDir()
	if err != nil {
		return err
	}

	// Create download directory and its parent
	must(os.MkdirAll(godlDownloadDir, os.ModePerm))

	exists, err := versionExists(archiveVersion)
	// handle stat errors even when file exists
	if err != nil {
		return err
	}
	// return early if archive is already downloaded and forceDownload is false
	if exists && !forceDownload {
		fmt.Println("archive has already been downloaded")
		return nil
	}

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := filepath.Join(godlDownloadDir, archiveName)

	// Create the file with tmp extension. So we don't overwrite until
	// the file is completely downloaded.
	tmpFile, err := goBinDown.fCR.Create(downloadPath + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	response, err := goBinDown.Client.Get(goBinDown.BaseURL + archiveName)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fileSize, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	// Create the writeCounter to be used with the writer
	wc := &writeCounter{TotalExpectedBytes: uint64(fileSize)}

	_, err = io.Copy(tmpFile, io.TeeReader(response.Body, wc))
	if err != nil {
		return err
	}

	fmt.Println()

	// Rename the temporary file once fully downloaded
	return goBinDown.fCR.Rename(downloadPath+".tmp", downloadPath)
}
