package cli

import (
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestDownloadCmd(t *testing.T) {
	testCases := map[string]struct {
		flags    string
		expected string
		useFlag  bool
	}{
		"CalledWithNoArgs": {
			flags:    "",
			useFlag:  false,
			expected: "Error: provide version to download\n",
		},
		"Help": {
			flags:    "-h",
			expected: "",
			useFlag:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			godl := NewRootCmd()
			download := NewDownloadCmd()
			godl.RegisterSubCommands([]*cobra.Command{download})

			var errOutput string
			if tc.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "download", tc.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "download")
			}
			if errOutput != tc.expected {
				t.Errorf("godl download failed: expected: %s; got: %s", tc.expected, errOutput)
			}
		})
	}
}
