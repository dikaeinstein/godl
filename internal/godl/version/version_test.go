package version

import (
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestVersion(t *testing.T) {
	godlCmd := godl.New()
	version := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{version})

	test.ExecuteCommand(t, false, godlCmd, "version")
}
