package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestVersion(t *testing.T) {
	godl := NewRootCmd()
	version := NewVersionCmd()
	godl.RegisterSubCommands([]*cobra.Command{version})

	test.ExecuteCommand(t, false, godl.CobraCmd, "version")
}
