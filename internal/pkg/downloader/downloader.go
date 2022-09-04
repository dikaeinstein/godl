package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/dikaeinstein/downloader/pkg/downloader"

	"github.com/dikaeinstein/godl/internal/pkg/version"
)

const (
	postfix = "darwin-amd64.tar.gz"
	prefix  = "go"
)

type Downloader struct {
	BaseURL       string
	Client        *http.Client
	DownloadDir   string
	FS            fs.FS
	ForceDownload bool
	Hasher        downloader.Hasher
	HashVerifier  downloader.HashVerifier
}

// Download downloads a binary release of a given version.
func (d *Downloader) Download(ctx context.Context, ver string) error {
	dl, err := downloader.New(
		d.DownloadDir, d.Client, d.FS,
		d.Hasher, &downloader.ProgressBar{}, d.HashVerifier,
	)
	if err != nil {
		return err
	}

	exists, err := version.Exists(ver, d.DownloadDir)
	if err != nil {
		return err
	}

	if exists && !d.ForceDownload {
		fmt.Println("archive has already been downloaded")
		return nil
	}

	err = d.CheckIfExistsRemote(ctx, ver)
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("%s%s.%s", prefix, ver, postfix)
	goURL := d.versionURL(ver)

	return dl.Download(ctx, goURL, archiveName, goURL+".sha256")
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
