package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestUpdateCmd(t *testing.T) {
	testCases := []struct {
		desc        string
		errOutput   string
		godlVersion string
		output      string
	}{
		{
			desc:        "Writes correct message to stdout when no update",
			errOutput:   "",
			godlVersion: "v0.11.6",
			output:      "No update available.\n",
		},
		{
			desc:        "Writes correct message to stdout when update is available",
			errOutput:   "",
			godlVersion: "v0.11.5",
			output: `Your version of Godl is out of date! The latest version
 is v0.11.6. You can update by downloading from https://github.com/dikaeinstein/godl/releases
`,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			update := NewUpdateCmd()
			godl := NewRootCmd()
			godl.RegisterSubCommands([]*cobra.Command{update})

			godlVersion = tC.godlVersion
			output, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "update")
			if errOutput != tC.errOutput {
				t.Errorf("godl update failed: expected errOutput %s; got %s", tC.errOutput, errOutput)
			}
			if output != tC.output {
				t.Errorf("godl update failed: expected output %s; got %s", tC.output, output)
			}
		})
	}
}
