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
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type writeCloseNamer interface {
	io.WriteCloser
	Name() string
}

type fileCreator interface {
	Create(name string) (writeCloseNamer, error)
}

type fileCreatorRenamer interface {
	fileCreator
	Rename(oldPath, newPath string) error
}

type fsFileCreatorRenamer struct{}

func (ac fsFileCreatorRenamer) Create(name string) (writeCloseNamer, error) {
	return os.Create(name)
}

func (ac fsFileCreatorRenamer) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

var forceDownload bool

const (
	postfix = "darwin-amd64.tar.gz"
	prefix  = "go"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [version]",
	Short: "Download go binary archive.",
	Long: `Download the archive version from https://golang.org/dl/ and save to $HOME/godl/downloads.

By default, if archive version already exists locally, godl doesn't attempt to download it again.
To force it to download the version again pass the --force flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		archiveVersion := args[0]
		fcr := fsFileCreatorRenamer{}

		goBinDownloader := goBinaryDownloader{
			Client: &http.Client{},
			// BaseURL: "https://dl.google.com/go/",
			BaseURL: "https://storage.googleapis.com/golang/",
			fCR:     fcr,
		}

		godlDownloadDir, err := getDownloadDir()
		if err != nil {
			return err
		}

		fmt.Printf("Downloading go binary %v\n", archiveVersion)
		err = goBinDownloader.download(archiveVersion, godlDownloadDir, forceDownload)
		if err != nil {
			return fmt.Errorf("error downloading %v: %v", archiveVersion, err)
		}
		fmt.Println("\nDownload complete")
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
	bytesWritten       int64
	TotalExpectedBytes int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.bytesWritten += int64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *writeCounter) PrintProgress() {
	percentDownloaded := float64(wc.bytesWritten) / float64(wc.TotalExpectedBytes) * 100
	fmt.Printf("\rDownloading... %.0f%% complete", math.Round(percentDownloaded))
}

type hashGenerator func(url string) (string, error)
type hashVerifier func(file, wantHash string) error

type goBinaryDownloader struct {
	Client     *http.Client
	BaseURL    string
	fCR        fileCreatorRenamer
	genHash    hashGenerator
	verifyHash hashVerifier
}

func (goBinDown *goBinaryDownloader) download(version, downloadDir string, forceDownload bool) error {
	// Create download directory and its parent
	must(os.MkdirAll(downloadDir, os.ModePerm))

	exists, err := versionExists(version, downloadDir)
	// handle stat errors even when file exists
	if err != nil {
		return err
	}
	// return early if archive is already downloaded and forceDownload is false
	if exists && !forceDownload {
		fmt.Println("archive has already been downloaded")
		return nil
	}

	if err = checkIfExistsRemote(goBinDown.BaseURL, version); err != nil {
		return err
	}

	archiveName := fmt.Sprintf("%s%s.%s", prefix, version, postfix)
	downloadPath := filepath.Join(downloadDir, archiveName)

	// Create the file with tmp extension. So we don't overwrite until
	// the file is completely downloaded.
	tmpFile, err := goBinDown.fCR.Create(downloadPath + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	goURL := versionURL(goBinDown.BaseURL, version)
	res, err := goBinDown.Client.Get(goURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	// Create the writeCounter to be used with the writer
	wc := &writeCounter{TotalExpectedBytes: res.ContentLength}

	n, err := io.Copy(tmpFile, io.TeeReader(res.Body, wc))
	if err != nil {
		return err
	}
	if res.ContentLength != -1 && res.ContentLength != n {
		return fmt.Errorf("copied %v bytes; expected %v", n, res.ContentLength)
	}

	wantHex, err := goBinDown.genHash(goURL + ".sha256")
	if err != nil {
		return err
	}

	if err = goBinDown.verifyHash(tmpFile.Name(), wantHex); err != nil {
		return fmt.Errorf("error verifying SHA256 of %v: %v", tmpFile, err)
	}

	// Rename the temporary file once fully downloaded
	return goBinDown.fCR.Rename(downloadPath+".tmp", downloadPath)
}

func versionURL(baseURL, version string) string {
	return baseURL + prefix + version + "." + postfix
}

// getBinaryHash downloads the given URL and returns it as a string.
func getBinaryHash(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url, res.Status)
	}

	urlHash, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url, err)
	}

	return string(urlHash), nil
}

// verifyHash reports whether the named file has contents with
// SHA-256 of the given hex value.
func verifyHash(file, hex string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return err
	}
	if hex != fmt.Sprintf("%x", hash.Sum(nil)) {
		return fmt.Errorf("%s corrupt? does not have expected SHA-256 of %v", file, hex)
	}

	return nil
}

func checkIfExistsRemote(baseURL, version string) error {
	u := versionURL(baseURL, version)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no binary release of %v", version)
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}
