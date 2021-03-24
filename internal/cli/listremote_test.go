package cli

import (
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
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

	lsRemote := NewListRemoteCmd(testClient)
	godl := NewRootCmd()
	godl.RegisterSubCommands([]*cobra.Command{lsRemote})

	_, errOutput := test.ExecuteCommand(t, false, godl.CobraCmd, "list-remote")
	expected := ""
	if errOutput != expected {
		t.Errorf("godl list failed: expected %s; got %s", expected, errOutput)
	}
}
