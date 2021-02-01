package listremote

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/test"
	"github.com/google/go-cmp/cmp"
)

func TestListRemoteVersions(t *testing.T) {
	testClient := test.NewTestClient(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "..", "test", "testdata", "listbucketresult.xml"))
		if err != nil {
			panic(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}
	})

	failingTestClient := test.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	})

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
			err := listRemoteVersions(tc.input)
			if err != nil {
				diff := cmp.Diff(tc.want, err.Error())
				if diff != "" {
					t.Errorf(diff)
				}
			}
		})
	}
}
