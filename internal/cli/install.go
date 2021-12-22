package cli

import (
	"errors"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/dikaeinstein/godl/internal/app/install"
	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

// NewInstallCmd returns the install command
func NewInstallCmd(client *http.Client) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install version",
		Short: "Installs the specified go binary version from local or remote.",
		Long: heredoc.Doc(`
			Installs the specified go binary version from local or remote.
			It fetches the version from the remote if not found locally before installing it.
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("provide version to install")
			}
			return nil
		},
	}

	const defaultTimeout = 60 * time.Second
	forceDownload := installCmd.Flags().BoolP("force", "f", false, "Force download instead of using local version")
	timeout := installCmd.Flags().DurationP("timeout", "t", defaultTimeout, "Set the download timeout.")

	installCmd.RunE = func(cmd *cobra.Command, args []string) error {
		dlDir, err := godlutil.GetDownloadDir()
		if err != nil {
			return err
		}

		install := install.Install{
			Archiver: &archiver.TarGz{
				Tar: &archiver.Tar{
					OverwriteExisting: true,
				},
				CompressionLevel: -1,
			},
			Dl: &downloader.Downloader{
				BaseURL:       "https://storage.googleapis.com/golang/",
				Client:        client,
				DownloadDir:   dlDir,
				FS:            fsys.OsFS{},
				ForceDownload: *forceDownload,
				Hasher:        hash.NewRemoteHasher(http.DefaultClient),
				HashVerifier:  godlutil.VerifyHash,
			},
			Timeout: *timeout,
		}

		return install.Run(cmd.Context(), args[0])
	}

	return installCmd
}
