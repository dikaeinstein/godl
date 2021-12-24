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

	testCases := []struct {
		name     string
		expected string
		flags    string
		useFlag  bool
	}{
		{
			name:     "CalledWithNoArgs",
			expected: "Error: provide version to install\n",
			flags:    "",
			useFlag:  false,
		},
		{
			name:     "Help",
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

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			godl := NewRootCmd()
			install := NewInstallCmd(testClient)
			godl.RegisterSubCommands([]*cobra.Command{install})

			var errOutput string
			if tC.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "install", tC.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl.CobraCmd, "install")
			}
			if errOutput != tC.expected {
				t.Errorf("godl install failed: expected: %s; got: %s", tC.expected, errOutput)
			}
		})
	}
}
