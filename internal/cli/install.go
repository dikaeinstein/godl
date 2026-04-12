package cli

import (
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/downloader"
)

// newInstallCmd returns the install command
func newInstallCmd(
	client *http.Client,
	dl *downloader.Downloader,
	installer *app.Install,
) *cobra.Command {
	iCli := &installCli{httpClient: client, dl: dl, installer: installer}

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
	dl         *downloader.Downloader
	cfg        installConfig
	installer  *app.Install
}

func (iCli *installCli) run(cmd *cobra.Command, args []string) error {
	iCli.dl.Configure(iCli.cfg.forceDownload)
	iCli.installer.Configure(iCli.dl, iCli.cfg.timeout)
	return iCli.installer.Run(cmd.Context(), args[0], iCli.cfg.forceDownload)
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
