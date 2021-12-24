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

	testCases := []struct {
		name     string
		flags    string
		expected string
		useFlag  bool
	}{
		{
			name:     "CalledWithNoArgs",
			flags:    "",
			useFlag:  false,
			expected: "Error: provide version to download\n",
		},
		{
			name:     "Help",
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

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			godl := NewRootCmd()
			download := NewDownloadCmd(testClient)
			godl.RegisterSubCommands([]*cobra.Command{download})

			var errOutput string
			if tC.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "download", tC.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "download")
			}
			if errOutput != tC.expected {
				t.Errorf("godl download failed: expected: %s; got: %s", tC.expected, errOutput)
			}
		})
	}
}
