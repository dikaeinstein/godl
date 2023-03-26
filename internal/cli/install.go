package cli

import (
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/downloader/pkg/hash"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

// newInstallCmd returns the install command
func newInstallCmd(client *http.Client) *cobra.Command {
	iCli := &installCli{httpClient: client}

	installCmd := &cobra.Command{
		Use:   "install version",
		Short: "Installs the specified go binary version from local or remote.",
		Long: heredoc.Doc(`
			Installs the specified go binary version from local or remote.
			It fetches the version from the remote if not found locally before installing it.
		`),
		Args:    cobra.ExactArgs(1),
		PreRunE: iCli.setupConfig,
		RunE:    iCli.run,
	}

	setupInstallCliFlags(installCmd)

	return installCmd
}

type installConfig struct {
	timeout       time.Duration
	forceDownload bool
}

type installCli struct {
	httpClient *http.Client
	cfg        installConfig
}

func (iCli *installCli) run(cmd *cobra.Command, args []string) error {
	dlDir, err := godlutil.GetDownloadDir()
	if err != nil {
		return err
	}

	install := app.Install{
		Archiver: &archiver.TarGz{
			Tar: &archiver.Tar{
				OverwriteExisting: true,
			},
			CompressionLevel: -1,
		},
		Dl: &downloader.Downloader{
			BaseURL:       "https://storage.googleapis.com/golang/",
			Client:        iCli.httpClient,
			DownloadDir:   dlDir,
			FS:            fsys.OsFS{},
			ForceDownload: iCli.cfg.forceDownload,
			Hasher:        hash.NewRemoteHasher(http.DefaultClient),
			HashVerifier:  hash.Verifier{},
		},
		Timeout: iCli.cfg.timeout,
	}

	return install.Run(cmd.Context(), args[0])
}

func (iCli *installCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	iCli.cfg.timeout = viper.GetDuration("timeout")
	iCli.cfg.forceDownload = viper.GetBool("force")

	return nil
}

func setupInstallCliFlags(cmd *cobra.Command) {
	const defaultTimeout = 60 * time.Second
	cmd.Flags().BoolP("force", "f", false,
		"Force download instead of using local version.")
	cmd.Flags().DurationP("timeout", "t", defaultTimeout,
		"Set the download timeout.")
}
