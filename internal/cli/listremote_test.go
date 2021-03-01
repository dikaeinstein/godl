package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestListRemoteCommand(t *testing.T) {
	lsRemote := NewListRemoteCmd()
	godl := NewRootCmd()
	godl.RegisterSubCommands([]*cobra.Command{lsRemote})

	_, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "list-remote")
	expected := ""
	if errOutput != expected {
		t.Errorf("godl list failed: expected %s; got %s", expected, errOutput)
	}
}
