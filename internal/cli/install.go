package cli

import (
	"net/http"
	"runtime"
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
	os            string
	arch          string
}

type installCli struct {
	httpClient *http.Client
	dl         *downloader.Downloader
	cfg        installConfig
	installer  *app.Install
}

func (iCli *installCli) run(cmd *cobra.Command, args []string) error {
	iCli.dl.Configure(iCli.cfg.forceDownload)
	iCli.installer.Configure(iCli.cfg.timeout)
	return iCli.installer.Run(
		cmd.Context(),
		args[0],
		iCli.cfg.os,
		iCli.cfg.arch,
		iCli.cfg.forceDownload,
	)
}

func (iCli *installCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	iCli.cfg.timeout = viper.GetDuration("timeout")
	iCli.cfg.forceDownload = viper.GetBool("force")
	iCli.cfg.os = viper.GetString("os")
	iCli.cfg.arch = viper.GetString("arch")

	return nil
}

func setupInstallCliFlags(cmd *cobra.Command) {
	const defaultTimeout = 60 * time.Second
	cmd.Flags().BoolP("force", "f", false,
		"Force download instead of using local version.")
	cmd.Flags().DurationP("timeout", "t", defaultTimeout,
		"Set the download timeout.")
	cmd.Flags().StringP("os", "o", runtime.GOOS,
		`Set the target OS. One of darwin, freebsd, linux, and so on.
To view possible combinations of GOOS and GOARCH, run "go tool dist list".`)
	cmd.Flags().StringP("arch", "a", runtime.GOARCH, `Set the target architecture.
GOARCH is the running program's architecture target: one of 386, amd64, arm, s390x, and so on.`)
}
