package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type VersionOption struct {
	BuildDate   string
	GitHash     string
	GodlVersion string
	GoVersion   string
}

// newVersionCmd returns the version command
func newVersionCmd(v VersionOption) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the godl version information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s\n",
				v.GodlVersion, v.GoVersion, v.GitHash, v.BuildDate)
		},
	}
}
