package app

import (
	"testing"

	"github.com/dikaeinstein/godl/internal/pkg/version"
	"github.com/dikaeinstein/godl/test"
)

func TestListDownloadedBinaryArchives(t *testing.T) {
	tmpFile, tmpDir := test.CreateTempGoBinaryArchive(t, "1.13")
	defer tmpFile.Close()

	ls := List{}
	if err := ls.Run(tmpDir, version.SortAsc); err != nil {
		t.Errorf("Error listing downloaded archive versions: %v", err)
	}
}
