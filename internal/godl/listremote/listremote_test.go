package listremote

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/test"
	"github.com/google/go-cmp/cmp"
)

func TestListRemoteVersions(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "..", "test", "testdata", "listbucketresult.xml"))
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
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}))

	tests := map[string]struct {
		input *http.Client
		want  string
	}{
		"getBinaryReleases succeeds": {input: testClient},
		"handles getBinaryReleases error": {
			input: failingTestClient,
			want:  "\nerror fetching list: https://storage.googleapis.com/golang/?prefix=go1: ",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lsRemote := ListRemote{tc.input}
			err := lsRemote.Run(context.Background(), gv.Asc)
			if err != nil {
				diff := cmp.Diff(tc.want, err.Error())
				if diff != "" {
					t.Errorf(diff)
				}
			}
		})
	}
}
