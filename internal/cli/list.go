package cli

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/internal/pkg/version"
	"github.com/dikaeinstein/godl/pkg/text"
)

// newListCmd returns a new instance of the list command
func newListCmd() *cobra.Command {
	lsCli := &lsCli{}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the downloaded versions.",
		Example: text.Indent(heredoc.Doc(`
			$ ls -s asc or ls -s=asc
			$ ls -s desc or ls -s=desc
		`), "  "),
		PreRunE: lsCli.setupConfig,
		RunE:    lsCli.run,
	}

	listCmd.Flags().StringP("sortDirection", "s", version.SortAsc,
		"Specify the sort direction of the output of `list`. "+
			"It sorts in ascending order by default.")

	return listCmd
}

type lsConfig struct{ sortDirection string }

type lsCli struct{ lsConfig }

func (c *lsCli) setupConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	c.lsConfig.sortDirection = viper.GetString("sortDirection")

	return nil
}

func (c *lsCli) run(cmd *cobra.Command, args []string) error {
	d, err := godlutil.GetDownloadDir()
	if err != nil {
		return err
	}

	ls := app.List{}
	return ls.Run(d, c.lsConfig.sortDirection)
}
