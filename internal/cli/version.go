package cli

import (
	"github.com/spf13/cobra"

	"github.com/dikaeinstein/godl/internal/app"
)

// newVersionCmd returns the version command
func newVersionCmd(info app.BuildInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the godl version information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.NewVersion(info).Run(cmd.OutOrStdout())
		},
	}
}
