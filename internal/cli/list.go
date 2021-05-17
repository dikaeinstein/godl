package cli

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/dikaeinstein/godl/internal/godl/list"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/pkg/text"
	"github.com/spf13/cobra"
)

// New returns a new instance of the list command
func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the downloaded versions.",
		Example: text.Indent(heredoc.Doc(`
			$ ls -s asc or ls -s=asc
			$ ls -s desc or ls -s=desc
		`), "  "),
	}

	sortDirection := listCmd.Flags().StringP("sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list`. It sorts in ascending order by default.")

	listCmd.RunE = func(cmd *cobra.Command, args []string) error {
		d, err := godlutil.GetDownloadDir()
		if err != nil {
			return err
		}

		ls := list.List{}
		return ls.Run(d, gv.SortDirection(*sortDirection))
	}

	return listCmd
}
