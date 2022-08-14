package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/test"
)

func TestCompletionCmdCalledWithNoArgs(t *testing.T) {
	godl := newRootCmd()
	completion := newCompletionCmd()
	registerSubCommands(godl, []*cobra.Command{completion})

	_, errOutput := test.ExecuteCommand(t, true, godl, "completion")
	expected := "Error: accepts 1 arg(s), received 0\n"
	require.Equal(t, expected, errOutput)
}

func TestCompletionCmdCalledWithInvalidArgs(t *testing.T) {
	godl := newRootCmd()
	completion := newCompletionCmd()
	registerSubCommands(godl, []*cobra.Command{completion})

	_, errOutput := test.ExecuteCommand(t, true, godl, "completion", "powershell")
	expected := `Error: invalid argument "powershell" for "godl completion"
`
	require.Equal(t, expected, errOutput)
}

func TestCompletionCmd(t *testing.T) {
	testCases := []struct {
		name  string
		shell string
	}{
		{"CalledWithBash", "bash"},
		{"CalledWithZsh", "zsh"},
		{"CalledWithFish", "fish"},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			godl := newRootCmd()
			completion := newCompletionCmd()
			registerSubCommands(godl, []*cobra.Command{completion})

			_, errOutput := test.ExecuteCommand(t, true, godl, "completion", tC.shell)
			require.Equal(t, "", errOutput)
		})
	}
}
