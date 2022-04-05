package cli

import (
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/pkg/text"
	"github.com/spf13/cobra"
)

// NewListRemoteCmd returns the list-remote command
func NewListRemoteCmd(client *http.Client) *cobra.Command {
	listRemoteCmd := &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Short:   "List the available remote versions.",
		Example: text.Indent(heredoc.Doc(`
			$ ls-remote
			$ ls-remote -s asc or ls-remote -s=asc
			$ ls-remote -s desc or ls-remote -s=desc
			$ ls-remote -t 15s or ls-remote --timeout=15s
		`), "  "),
	}

	const defaultTimeout = 60 * time.Second
	sd := listRemoteCmd.Flags().StringP("sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list-remote`. It sorts in ascending order by default.")
	timeout := listRemoteCmd.Flags().DurationP("timeout", "t", defaultTimeout, "Set the download timeout.")

	listRemoteCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lsRemote := app.ListRemote{Client: client, Timeout: *timeout}
		return lsRemote.Run(cmd.Context(), gv.SortDirection(*sd))
	}

	return listRemoteCmd
}
