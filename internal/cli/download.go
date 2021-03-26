package cli

import (
	"errors"
	"net/http"
	"time"

	"github.com/dikaeinstein/godl/internal/godl/download"
	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/spf13/cobra"
)

// New returns a new instance of the download command
func NewDownloadCmd(client *http.Client) *cobra.Command {
	downloadCmd := &cobra.Command{
		Use:   "download version",
		Short: "Download go binary archive.",
		Long: `Download the archive version from https://golang.org/dl/ and save to $HOME/godl/downloads.

	By default, if archive version already exists locally, godl doesn't attempt to download it again.
	To force it to download the version again pass the --force flag.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("provide version to download")
			}
			return nil
		},
	}

	const defaultTimeout = 60 * time.Second
	forceDownload := downloadCmd.Flags().BoolP("force", "f", false, "Force download instead of using local version")
	timeout := downloadCmd.Flags().DurationP("timeout", "t", defaultTimeout, "Set the download timeout.")

	downloadCmd.RunE = func(cmd *cobra.Command, args []string) error {
		dlDir, err := godlutil.GetDownloadDir()
		if err != nil {
			return err
		}

		dl := &downloader.Downloader{
			BaseURL:       "https://storage.googleapis.com/golang/",
			Client:        client,
			DownloadDir:   dlDir,
			FS:            fsys.OsFS{},
			ForceDownload: *forceDownload,
			Hasher:        hash.NewRemoteHasher(http.DefaultClient),
			HashVerifier:  godlutil.VerifyHash,
		}

		d := download.Download{
			Dl:      dl,
			Timeout: *timeout,
		}
		return d.Run(cmd.Context(), args[0])
	}

	return downloadCmd
}
