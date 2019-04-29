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
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type fileCreator interface {
	Create(name string) (io.WriteCloser, error)
	Rename(oldPath, newPath string) error
}

type fsFileCreator struct{}

func (ac *fsFileCreator) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (ac *fsFileCreator) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [version] [downloadDir]",
	Short: "Download go binary archive",
	Long: `Download the archive version from https://golang.org/dl/ and save to specified directory.

If no download directory is specified, godl downloads the archive into
$HOME/godl/downloads.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("%v: home directory cannot be detected", err)
		}

		var downloadDir string
		if len(args) > 1 && len(args) <= 2 {
			downloadDir = args[1]
		} else {
			godlDownloadDir := path.Join(home, "godl", "downloads")
			must(os.MkdirAll(godlDownloadDir, os.ModePerm))
			downloadDir = godlDownloadDir
		}

		archiveVersion := args[0]
		fc := &fsFileCreator{}

		goBinDownloader := goBinaryDownloader{
			Client:  &http.Client{},
			BaseURL: "https://dl.google.com/go/",
			FC:      fc,
		}

		fmt.Printf("Downloading go binary %v\n", archiveVersion)
		goBinDownloader.download(archiveVersion, downloadDir)
		fmt.Println("Download complete")
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
	FC      fileCreator
}

func (goBinDown *goBinaryDownloader) download(archiveVersion, downloadDir string) error {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := path.Join(downloadDir, archiveName)

	// Create the file with tmp extension. So we don't overwrite until
	// the file is completely downloaded.
	tmpFile, err := goBinDown.FC.Create(downloadPath + ".tmp")
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
	return goBinDown.FC.Rename(downloadPath+".tmp", downloadPath)
}
