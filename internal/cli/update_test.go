package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/test"
)

func TestUpdateCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestUpdateCmd in short mode.")
	}

	testCases := []struct {
		name        string
		errOutput   string
		godlVersion string
		output      string
	}{
		{
			name:        "Writes correct message to stdout when no update",
			errOutput:   "",
			godlVersion: "v0.11.6",
			output:      "No update available.\n",
		},
		{
			name:        "Writes correct message to stdout when update is available",
			errOutput:   "",
			godlVersion: "v0.11.5",
			output: heredoc.Doc(`
				Your version of Godl is out of date!

				The latest version is v0.11.6.
				You can update by downloading from https://github.com/dikaeinstein/godl/releases
			`),
		},
	}

	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "test", "testdata", "releases.json"))
		if err != nil {
			panic(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}
	}))

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			v := VersionOption{GodlVersion: tC.godlVersion}
			update := newUpdateCmd(testClient, v)
			godl := newRootCmd()
			registerSubCommands(godl, []*cobra.Command{update})

			output, errOutput := test.ExecuteCommand(t, false, godl, "update")

			require.Equal(t, tC.errOutput, errOutput)
			require.Equal(t, tC.output, output)
		})
	}
}
