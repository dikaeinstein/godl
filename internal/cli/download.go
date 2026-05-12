package cli

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/downloader/pkg/hash"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/downloader"
	"github.com/dikaeinstein/godl/internal/godlutil"
	"github.com/dikaeinstein/godl/pkg/fsys"
)

// newDownloadCmd returns a new instance of the download command
func newDownloadCmd(client *http.Client) *cobra.Command {
	dCli := &downloadCli{httpClient: client}

	downloadCmd := &cobra.Command{
		Use:   "download version",
		Short: "Download go binary archive.",
		Long: heredoc.Doc(`
			Download the archive version from https://golang.org/dl/ and save to $HOME/.godl/downloads.

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
	os            string
	arch          string
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

	dl, err := downloader.New(
		fsys.OsFS{},
		hash.NewRemoteHasher(dCli.httpClient),
		hash.Verifier{},
		dCli.httpClient,
		"https://dl.google.com/go/",
		dlDir,
		dCli.downloadConfig.forceDownload,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize downloader: %w", err)
	}

	d := app.Download{
		Dl:      dl,
		Timeout: dCli.downloadConfig.timeout,
	}

	return d.Run(cmd.Context(), args[0], dCli.downloadConfig.os, dCli.downloadConfig.arch)
}

func (dCli *downloadCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	dCli.downloadConfig.timeout = viper.GetDuration("timeout")
	dCli.downloadConfig.forceDownload = viper.GetBool("force")
	dCli.downloadConfig.os = viper.GetString("os")
	dCli.downloadConfig.arch = viper.GetString("arch")

	return nil
}

func setupDownloadFlags(cmd *cobra.Command) {
	const defaultTimeout = 60 * time.Second
	cmd.Flags().BoolP("force", "f", false,
		"Force download instead of using local version")
	cmd.Flags().DurationP("timeout", "t", defaultTimeout,
		"Set the download timeout.")
	cmd.Flags().StringP("os", "o", runtime.GOOS,
		`Set the target OS. One of darwin, freebsd, linux, and so on.
To view possible combinations of GOOS and GOARCH, run "go tool dist list".`)
	cmd.Flags().StringP("arch", "a", runtime.GOARCH, `Set the target architecture.
It is the running program's architecture target: one of 386, amd64, arm, s390x, and so on.`)
}
