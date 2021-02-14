package godl_test

import (
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/internal/godl/version"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestGodlCmd(t *testing.T) {
	godlCmd := godl.New()
	_, errOutput := test.ExecuteCommand(t, true, godlCmd)
	if errOutput != "" {
		t.Errorf("calling command without subcommands should not have error: %v", errOutput)
	}
}

func TestGodlExecuteUnknownCommand(t *testing.T) {
	godlCmd := godl.New()
	// Register version subcommand so there's a list to filter an unknown command against.
	godlCmd.RegisterSubCommands([]*cobra.Command{version.New()})

	_, errOutput := test.ExecuteCommand(t, true, godlCmd, "unknown")
	expected := "Error: unknown command \"unknown\" for \"godl\"\nRun 'godl --help' for usage.\n"

	if errOutput != expected {
		t.Errorf("expected: %q\ngot:\n %q", expected, errOutput)
	}
}
