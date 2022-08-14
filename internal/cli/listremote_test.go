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

func TestListRemoteCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestListRemoteCmd in short mode.")
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

	lsRemote := newListRemoteCmd(testClient)

	godl := newRootCmd()
	registerSubCommands(godl, []*cobra.Command{lsRemote})

	_, errOutput := test.ExecuteCommand(t, false, godl, "list-remote")

	require.Equal(t, "", errOutput)
}
