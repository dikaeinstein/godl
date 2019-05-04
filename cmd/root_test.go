package cmd

import "testing"

func TestRootCmd(t *testing.T) {
	_, err := executeCommand(rootCmd)
	if err != nil {
		t.Errorf("Calling command without subcommands should not have error: %v", err)
	}
}

func TestRootExecuteUnknownCommand(t *testing.T) {
	output, _ := executeCommand(rootCmd, "unknown")
	expected := "Error: unknown command \"unknown\" for \"godl\"\nRun 'godl --help' for usage.\n"

	if output != expected {
		t.Errorf("Expected:\n %q\nGot:\n %q\n", expected, output)
	}
}
