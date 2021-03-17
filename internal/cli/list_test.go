package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestListCmd(t *testing.T) {
	ls := NewListCmd()
	godl := NewRootCmd()
	godl.RegisterSubCommands([]*cobra.Command{ls})

	_, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "list")
	expected := ""
	if errOutput != expected {
		t.Errorf("godl list failed: expected %s; got %s", expected, errOutput)
	}
}
