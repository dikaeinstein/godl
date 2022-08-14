package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/internal/pkg/version"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

const (
	postfix = "darwin-amd64.tar.gz"
	prefix  = "go"
)

// Hasher generates a hash from given path to file.
//
// Hash returns the hash of the file at the given path.
type Hasher interface {
	Hash(ctx context.Context, path string) (string, error)
}

// HashVerifier verifies an input io.Reader content against wantHash
type HashVerifier func(input io.Reader, wantHash string) error

type Downloader struct {
	BaseURL       string
	Client        *http.Client
	DownloadDir   string
	FS            fs.FS
	ForceDownload bool
	Hasher        Hasher
	HashVerifier  HashVerifier
}

// Download downloads a binary release of a given version.
func (d *Downloader) Download(ctx context.Context, ver string) error {
	// Create download directory and its parent
	godlutil.Must(os.MkdirAll(d.DownloadDir, os.ModePerm))

	exists, err := version.Exists(ver, d.DownloadDir)
	// handle stat errors even when file exists
	if err != nil {
		return err
	}
	// return early if archive is already downloaded and forceDownload is false
	if exists && !d.ForceDownload {
		fmt.Println("archive has already been downloaded")
		return nil
	}

	err = d.CheckIfExistsRemote(ctx, ver)
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("%s%s.%s", prefix, ver, postfix)
	downloadPath := filepath.Join(d.DownloadDir, archiveName)

	// Create the file with tmp extension. So we don't overwrite until
	// the file is completely downloaded.
	tmp, err := fsys.Create(d.FS, downloadPath+".tmp")
	if err != nil {
		return err
	}
	defer tmp.Close()

	goURL := d.versionURL(ver)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, goURL, http.NoBody)
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	// Create the countWriter used for counting response bytes
	cw := &countWriter{totalExpectedBytes: res.ContentLength}

	tmpFile, ok := tmp.(io.Writer)
	if !ok {
		return fmt.Errorf("invalid writer: %T", tmp)
	}

	n, err := io.Copy(tmpFile, io.TeeReader(res.Body, cw))
	if err != nil {
		return err
	}
	if res.ContentLength != -1 && res.ContentLength != n {
		return fmt.Errorf("copied %v bytes; expected %v", n, res.ContentLength)
	}

	wantHex, err := d.Hasher.Hash(ctx, goURL+".sha256")
	if err != nil {
		return err
	}

	fmt.Println("\nverifying checksum")
	f, err := d.FS.Open(downloadPath + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	err = d.HashVerifier(f, wantHex)
	if err != nil {
		return fmt.Errorf("error verifying SHA256 checksum of %v: %v", tmpFile, err)
	}

	fmt.Println("checksums matched!")

	// Rename the temporary file once fully downloaded
	return fsys.Rename(d.FS, downloadPath+".tmp", downloadPath)
}

func (d *Downloader) CheckIfExistsRemote(ctx context.Context, ver string) error {
	u := d.versionURL(ver)
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, u, http.NoBody)
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no binary release of %v", ver)
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}

func (d *Downloader) versionURL(ver string) string {
	return d.BaseURL + prefix + ver + "." + postfix
}

func Postfix() string {
	return postfix
}

func Prefix() string {
	return prefix
}

// countWriter counts the number of bytes written to it.
type countWriter struct {
	bytesWritten       int64
	totalExpectedBytes int64
}

func (wc *countWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.bytesWritten += int64(n)
	wc.Progress()
	return n, nil
}

// Progress prints the progress of bytes counted
func (wc *countWriter) Progress() {
	percentDownloaded := float64(wc.bytesWritten) / float64(wc.totalExpectedBytes) * 100
	fmt.Printf("\rDownloading... %.0f%% complete", math.Round(percentDownloaded))
}
