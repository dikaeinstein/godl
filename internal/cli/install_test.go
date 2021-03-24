package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestInstallCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestDownloadCmd in short mode.")
	}

	testCases := map[string]struct {
		expected string
		flags    string
		useFlag  bool
	}{
		"CalledWithNoArgs": {
			expected: "Error: provide version to install\n",
			flags:    "",
			useFlag:  false,
		},
		"Help": {
			expected: "",
			flags:    "-h",
			useFlag:  true,
		},
	}

	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "test", "testdata", "listbucketresult.xml"))
		if err != nil {
			panic(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}
	}))

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			godl := NewRootCmd()
			install := NewInstallCmd(testClient)
			godl.RegisterSubCommands([]*cobra.Command{install})

			var errOutput string
			if tc.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "install", tc.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "install")
			}
			if errOutput != tc.expected {
				t.Errorf("godl install failed: expected: %s; got: %s", tc.expected, errOutput)
			}
		})
	}
}
