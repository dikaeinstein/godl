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

func TestCompletionCmdCalledWithInvalidArgs(t *testing.T) {
	godl := NewRootCmd()
	completion := NewCompletionCmd(godl)
	godl.RegisterSubCommands([]*cobra.Command{completion})

	_, errOutput := test.ExecuteCommand(t, true, godl.CobraCmd, "completion", "powershell")
	expected := "Error: unknown shell passed\n"
	if errOutput != expected {
		t.Errorf("godl completion failed: expected: %s; got: %s", expected, errOutput)
	}
}

func TestCompletionCmd(t *testing.T) {
	testCases := map[string]struct {
		shell string
	}{"CalledWithBash": {"bash"}, "CalledWithZsh": {"zsh"}, "CalledWithFish": {"fish"}}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			godl := NewRootCmd()
			completion := NewCompletionCmd(godl)
			godl.RegisterSubCommands([]*cobra.Command{completion})

			_, errOutput := test.ExecuteCommand(t, true, godl.CobraCmd, "completion", tc.shell)
			if errOutput != "" {
				t.Errorf("godl completion failed to generate completion")
			}
		})
	}
}
