package list

import (
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestListDownloadedBinaryArchives(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := test.CreateTempGoBinaryArchive(tmpDir, "1.13")
	defer tmpFile.Close()

	if listDownloadedBinaryArchives(tmpDir); err != nil {
		t.Errorf("Error listing downloaded archive versions: %v", err)
	}
}

func TestListCommand(t *testing.T) {
	list := New()
	godlCmd := godl.New()
	godlCmd.RegisterSubCommands([]*cobra.Command{list})

	_, _, err := test.ExecuteCommand(godlCmd, "list")
	if err != nil {
		t.Errorf("godl list failed: %v", err)
	}
}
