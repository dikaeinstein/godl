package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestCompletionCmdCalledWithNoArgs(t *testing.T) {
	godl := NewRootCmd()
	completion := NewCompletionCmd(godl)
	godl.RegisterSubCommands([]*cobra.Command{completion})

	_, errOutput := test.ExecuteCommand(t, true, godl.CobraCmd, "completion")
	expected := "Error: provide shell to configure e.g bash or zsh\n"
	if errOutput != expected {
		t.Errorf("godl completion failed: expected: %s; got: %s", expected, errOutput)
	}
}
