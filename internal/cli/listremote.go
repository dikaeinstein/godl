package cli

import (
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/pkg/text"
)

// newListRemoteCmd returns the list-remote command
func newListRemoteCmd(client *http.Client) *cobra.Command {
	lsRemoteCli := &lsRemoteCli{httpClient: client}

	lsRemoteCmd := &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Short:   "List the available remote versions.",
		Example: text.Indent(heredoc.Doc(`
			$ ls-remote
			$ ls-remote -s asc or ls-remote -s=asc
			$ ls-remote -s desc or ls-remote -s=desc
			$ ls-remote -t 15s or ls-remote --timeout=15s
		`), "  "),
		PreRunE: lsRemoteCli.setupConfig,
		RunE:    lsRemoteCli.run,
	}

	setupLsRemoteCliFlags(lsRemoteCmd)

	return lsRemoteCmd
}

type lsRemoteConfig struct {
	timeout       time.Duration
	sortDirection string
}

type lsRemoteCli struct {
	httpClient *http.Client
	lsRemoteConfig
}

func (c *lsRemoteCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	c.lsRemoteConfig.timeout = viper.GetDuration("timeout")
	c.lsRemoteConfig.sortDirection = viper.GetString("sortDirection")

	return nil
}

func (c *lsRemoteCli) run(cmd *cobra.Command, args []string) error {
	lsRemote := app.ListRemote{Client: c.httpClient, Timeout: c.timeout}

	return lsRemote.Run(cmd.Context(), c.sortDirection)
}

func setupLsRemoteCliFlags(cmd *cobra.Command) {
	const defaultTimeout = 60 * time.Second
	cmd.Flags().DurationP(
		"timeout", "t", defaultTimeout, "Set the download timeout.")
	cmd.Flags().StringP("sortDirection", "s", "",
		"Specify the sort direction of the output of `list-remote`. "+
			"It sorts in ascending order by default.")
}
