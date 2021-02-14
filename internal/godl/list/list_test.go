package list

import (
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestListDownloadedBinaryArchives(t *testing.T) {
	tmpFile, tmpDir := test.CreateTempGoBinaryArchive(t, "1.13")
	defer tmpFile.Close()

	ls := listCmd{}
	if err := ls.Run(tmpDir); err != nil {
		t.Errorf("Error listing downloaded archive versions: %v", err)
	}
}

func TestListCommand(t *testing.T) {
	list := New()
	godlCmd := godl.New()
	godlCmd.RegisterSubCommands([]*cobra.Command{list})

	test.ExecuteCommand(t, false, godlCmd, "list")
}
