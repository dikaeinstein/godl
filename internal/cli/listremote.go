package cli

import (
	"fmt"
	"net/http"

	"github.com/dikaeinstein/godl/internal/godl/listremote"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/spf13/cobra"
)

// New returns the list-remote command
func NewListRemoteCmd() *cobra.Command {
	lsRemoteExAsc := "ls-remote -s asc or ls-remote -s=asc"
	lsRemoteExDesc := "ls-remote -s desc or ls-remote -s=desc"

	listRemoteCmd := &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Example: fmt.Sprintf("%11s\n%38s\n%40s", "ls-remote", lsRemoteExAsc, lsRemoteExDesc),
		Short:   "List the available remote versions.",
	}

	sd := listRemoteCmd.Flags().StringP("sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list-remote`. It sorts in ascending order by default.")

	listRemoteCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lsRemote := listremote.ListRemote{Client: http.DefaultClient}
		return lsRemote.Run(cmd.Context(), gv.SortDirection(*sd))
	}

	return listRemoteCmd
}
