package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/test"
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
			expected: "Error: accepts 1 arg(s), received 0\n",
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
			godl := newRootCmd()
			download := newDownloadCmd(testClient)
			registerSubCommands(godl, []*cobra.Command{download})

			var errOutput string
			if tC.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl, "download", tC.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl, "download")
			}

			require.Equal(t, tC.expected, errOutput)
		})
	}
}
