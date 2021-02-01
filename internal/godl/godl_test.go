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
	_, _, err := test.ExecuteCommand(godlCmd)
	if err != nil {
		t.Errorf("Calling command without subcommands should not have error: %v", err)
	}
}

func TestGodlExecuteUnknownCommand(t *testing.T) {
	godlCmd := godl.New()
	// Register version subcommand so there's a list to filter an unknown command against.
	version := version.New()
	godlCmd.RegisterSubCommands([]*cobra.Command{version})

	output, errOutput, _ := test.ExecuteCommand(godlCmd, "unknown")
	expected := "Error: unknown command \"unknown\" for \"godl\"\nRun 'godl --help' for usage.\n"

	if output != expected {
		t.Errorf("Expected:\n %q\nGot:\n %q\n %q", expected, output, errOutput)
	}
}
