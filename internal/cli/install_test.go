package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/downloader/pkg/hash"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/archives"
	"github.com/dikaeinstein/godl/internal/downloader"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/test"
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
			expected: "Error: accepts 1 arg(s), received 0\n",
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
			godl := newRootCmd()

			dl, err := downloader.New(
				fsys.OsFS{},
				hash.NewRemoteHasher(http.DefaultClient),
				hash.Verifier{},
				testClient,
				distURL,
				t.TempDir(),
				false,
			)
			require.NoError(t, err)

			installer := app.Install{Archiver: archives.NewTarGZ()}

			install := newInstallCmd(testClient, dl, &installer)
			registerSubCommands(godl, []*cobra.Command{install})

			var errOutput string
			if tC.useFlag {
				_, errOutput = test.ExecuteCommand(t, true, godl, "install", tC.flags)
			} else {
				_, errOutput = test.ExecuteCommand(t, true, godl, "install")
			}

			require.Equal(t, tC.expected, errOutput)
		})
	}
}
