package cli

import (
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/downloader/pkg/hash"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

// newDownloadCmd returns a new instance of the download command
func newDownloadCmd(client *http.Client) *cobra.Command {
	dCli := &downloadCli{httpClient: client}

	downloadCmd := &cobra.Command{
		Use:   "download version",
		Short: "Download go binary archive.",
		Long: heredoc.Doc(`
			Download the archive version from https://golang.org/dl/ and save to $HOME/godl/downloads.

			By default, if archive version already exists locally, godl doesn't attempt to download it again.
			To force it to download the version again pass the --force flag.
		`),
		Args:    cobra.ExactArgs(1),
		PreRunE: dCli.setupConfig,
		RunE:    dCli.run,
	}

	setupDownloadFlags(downloadCmd)

	return downloadCmd
}

type downloadConfig struct {
	timeout       time.Duration
	forceDownload bool
}

type downloadCli struct {
	httpClient *http.Client
	downloadConfig
}

func (dCli *downloadCli) run(cmd *cobra.Command, args []string) error {
	dlDir, err := godlutil.GetDownloadDir()
	if err != nil {
		return err
	}

	dl := &downloader.Downloader{
		BaseURL:       "https://storage.googleapis.com/golang/",
		Client:        dCli.httpClient,
		DownloadDir:   dlDir,
		FS:            fsys.OsFS{},
		ForceDownload: dCli.downloadConfig.forceDownload,
		Hasher:        hash.NewRemoteHasher(dCli.httpClient),
		HashVerifier:  hash.Verifier{},
	}

	d := app.Download{
		Dl:      dl,
		Timeout: dCli.downloadConfig.timeout,
	}

	return d.Run(cmd.Context(), args[0])
}

func (dCli *downloadCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	dCli.downloadConfig.timeout = viper.GetDuration("timeout")
	dCli.downloadConfig.forceDownload = viper.GetBool("force")

	return nil
}

func setupDownloadFlags(cmd *cobra.Command) {
	const defaultTimeout = 60 * time.Second
	cmd.Flags().BoolP("force", "f", false,
		"Force download instead of using local version")
	cmd.Flags().DurationP("timeout", "t", defaultTimeout,
		"Set the download timeout.")
}
