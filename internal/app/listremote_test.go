package app

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/test"
	"github.com/google/go-cmp/cmp"
)

func TestListRemoteVersions(t *testing.T) {
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

	failingTestClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	}))

	testCases := []struct {
		name   string
		client *http.Client
		want   string
	}{
		{name: "getBinaryReleases succeeds", client: testClient},
		{
			name:   "handles getBinaryReleases error",
			client: failingTestClient,
			want:   "\nerror fetching list: https://storage.googleapis.com/golang/?prefix=go1: ",
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			lsRemote := ListRemote{tC.client, 2 * time.Second}
			err := lsRemote.Run(context.Background(), gv.Asc)
			if err != nil {
				diff := cmp.Diff(tC.want, err.Error())
				if diff != "" {
					t.Errorf(diff)
				}
			}
		})
	}
}
