package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestGodlCmd(t *testing.T) {
	godlCmd := NewRootCmd()
	_, errOutput := test.ExecuteCommand(t, true, godlCmd.CobraCmd)
	if errOutput != "" {
		t.Errorf("calling command without subcommands should not have error: %v", errOutput)
	}
}

func TestGodlExecuteUnknownCommand(t *testing.T) {
	godlCmd := NewRootCmd()
	// Register list subcommand so there's a list to filter an unknown command against.
	godlCmd.RegisterSubCommands([]*cobra.Command{NewListCmd()})

	_, errOutput := test.ExecuteCommand(t, true, godlCmd.CobraCmd, "unknown")
	expected := "Error: unknown command \"unknown\" for \"godl\"\nRun 'godl --help' for usage.\n"

	if errOutput != expected {
		t.Errorf("expected: %q\ngot:\n %q", expected, errOutput)
	}
}
