package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	godlVersion = "unknown version"
	gitHash     = "unknown commit"
	goVersion   = "unknown go version"
	buildDate   = "unknown build date"
)

// New returns the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the godl version information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s\n",
				godlVersion, goVersion, gitHash, buildDate)
		},
	}
}
