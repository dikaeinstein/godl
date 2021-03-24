package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

func TestDownloadCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestDownloadCmd in short mode.")
	}

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
			download := NewDownloadCmd(testClient)
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
