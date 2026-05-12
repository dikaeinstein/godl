package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/dikaeinstein/godl/internal/version"
	"github.com/dikaeinstein/godl/test"
)

func TestListRemoteVersions(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "test", "testdata", "go_releases.json"))
		if err != nil {
			panic(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}
	}))

	failingTestClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)),
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	}))

	testCases := []struct {
		name   string
		client *http.Client
		want   error
	}{
		{name: "getBinaryReleases succeeds", client: testClient},
		{
			name:   "handles getBinaryReleases error",
			client: failingTestClient,
			want: errors.New(
				"could not fetch Go releases: https://go.dev/dl/?mode=json&include=all: 404 Not Found",
			),
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			lsRemote := ListRemote{tC.client, 2 * time.Second}
			err := lsRemote.Run(context.Background(), version.SortAsc)
			if err != nil {
				if err.Error() != tC.want.Error() {
					t.Errorf("got: %v, want: %v", err, tC.want)
				}
			}
		})
	}
}
