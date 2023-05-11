package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/dikaeinstein/godl/test"
)

func TestGodlCmd(t *testing.T) {
	godl := newRootCmd()
	_, errOutput := test.ExecuteCommand(t, true, godl)
	if errOutput != "" {
		t.Errorf(
			"calling command without subcommands should not have error: %v",
			errOutput,
		)
	}
}

func TestGodlExecuteUnknownCommand(t *testing.T) {
	godl := newRootCmd()
	ls := newListCmd()

	// register list subcommand so there's a list to filter an unknown command against
	registerSubCommands(godl, []*cobra.Command{ls})

	_, errOutput := test.ExecuteCommand(t, true, godl, "unknown")
	expected := "Error: unknown command \"unknown\" for \"godl\"\nRun 'godl --help' for usage.\n"

	if errOutput != expected {
		t.Errorf("expected: %q\ngot:\n %q", expected, errOutput)
	}
}
