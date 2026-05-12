package downloader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/dikaeinstein/downloader/pkg/downloader"

	"github.com/dikaeinstein/godl/internal/godlutil"
	"github.com/dikaeinstein/godl/internal/version"
)

type Downloader struct {
	FS            fs.FS
	Hasher        downloader.Hasher
	HashVerifier  downloader.HashVerifier
	Client        *http.Client
	BaseURL       string
	DownloadDir   string
	ForceDownload bool
	dl            *downloader.Downloader
}

func New(
	fsys fs.FS,
	hasher downloader.Hasher,
	hashVerifier downloader.HashVerifier,
	client *http.Client,
	baseURL string,
	downloadDir string,
	forceDownload bool,
) (*Downloader, error) {
	dl, err := downloader.New(
		downloadDir, client, fsys,
		hasher, &downloader.ProgressBar{}, hashVerifier,
	)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		FS:            fsys,
		Client:        client,
		BaseURL:       baseURL,
		DownloadDir:   downloadDir,
		ForceDownload: forceDownload,
		dl:            dl,
	}, nil
}

func (d *Downloader) Configure(forceDownload bool) {
	d.ForceDownload = forceDownload
}

// Download downloads a binary release of a given version.
func (d *Downloader) Download(ctx context.Context, ver, os, arch string) error {
	archiveName := godlutil.ArchiveName(ver, os, arch)

	exists, err := version.Exists(archiveName, d.DownloadDir)
	if err != nil {
		return err
	}

	if exists && !d.ForceDownload {
		fmt.Println("archive has already been downloaded")
		return nil
	}

	goReleaseURL := d.BaseURL + archiveName

	err = d.CheckIfExistsRemote(ctx, goReleaseURL)
	if err != nil {
		return err
	}

	return d.dl.Download(ctx, goReleaseURL, archiveName, goReleaseURL+".sha256")
}

func (d *Downloader) CheckIfExistsRemote(ctx context.Context, goReleaseURL string) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, goReleaseURL, http.NoBody)
	res, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no binary release of %v", goReleaseURL)
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}
