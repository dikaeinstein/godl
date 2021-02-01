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

	_, _, err := test.ExecuteCommand(godlCmd, "version")
	if err != nil {
		t.Errorf("godl version failed: %v", err)
	}
}
