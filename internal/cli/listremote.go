package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dikaeinstein/godl/internal/godl/listremote"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/spf13/cobra"
)

// New returns the list-remote command
func NewListRemoteCmd(client *http.Client) *cobra.Command {
	lsRemoteExAsc := "ls-remote -s asc or ls-remote -s=asc"
	lsRemoteExDesc := "ls-remote -s desc or ls-remote -s=desc"
	lsRemoteExTimeout := "ls-remote -t 15s or ls-remote --timeout=15s"

	listRemoteCmd := &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Example: fmt.Sprintf("%11s\n%38s\n%40s\n%45s", "ls-remote", lsRemoteExAsc, lsRemoteExDesc, lsRemoteExTimeout),
		Short:   "List the available remote versions.",
	}

	const defaultTimeout = 60 * time.Second
	sd := listRemoteCmd.Flags().StringP("sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list-remote`. It sorts in ascending order by default.")
	timeout := listRemoteCmd.Flags().DurationP("timeout", "t", defaultTimeout, "Set the download timeout.")

	listRemoteCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lsRemote := listremote.ListRemote{Client: client, Timeout: *timeout}
		return lsRemote.Run(cmd.Context(), gv.SortDirection(*sd))
	}

	return listRemoteCmd
}
