package update

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/dikaeinstein/godl/test"
)

func TestCheckForUpdate(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "..", "test", "testdata", "releases.json"))
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
			StatusCode: http.StatusBadGateway,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}))

	testCases := map[string]struct {
		client         *http.Client
		currentVersion string
		want           bool
		err            error
	}{
		"Returns true if update is available": {
			client:         testClient,
			currentVersion: "0.11.5",
			want:           true,
			err:            nil,
		},
		"Returns false if no update": {
			client:         testClient,
			currentVersion: "0.11.6",
			want:           false,
			err:            nil,
		},
		"Returns error encountered while checking for update": {
			client:         failingTestClient,
			currentVersion: "0.11.5",
			want:           false,
			err:            errors.New("https://api.github.com/repos/dikaeinstein/godl/releases?per_page=10: 502 "),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testWriter := &bytes.Buffer{}

			u := Update{tc.client, testWriter}
			exists, _, err := u.CheckForUpdate(context.Background(), tc.currentVersion)
			if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
				t.Errorf("expected CheckForUpdate(ctx, %v) => %v, got %v", tc.currentVersion, tc.err, err)
			}
			if exists != tc.want {
				t.Errorf("expected CheckForUpdate(ctx, %v) => %v, got %v", tc.currentVersion, tc.want, exists)
			}
		})
	}
}

func TestRun(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		f, err := os.Open(path.Join("..", "..", "..", "test", "testdata", "releases.json"))
		if err != nil {
			panic(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}
	}))
	testCases := []struct {
		client         *http.Client
		currentVersion string
		desc           string
		err            error
		want           string
	}{
		{
			desc:           "Returns no error",
			client:         testClient,
			currentVersion: "0.11.6",
			err:            nil,
			want:           "No update available.\n",
		},
		{
			desc:           "Returns correct message when no update is available",
			client:         testClient,
			currentVersion: "0.11.6",
			err:            nil,
			want:           "No update available.\n",
		},
		{
			desc:           "Returns correct message when update is available",
			client:         testClient,
			currentVersion: "0.11.5",
			err:            nil,
			want: `Your version of Godl is out of date! The latest version
 is v0.11.6. You can update by downloading from https://github.com/dikaeinstein/godl/releases
`,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testWriter := &bytes.Buffer{}

			u := Update{tC.client, testWriter}
			err := u.Run(context.Background(), tC.currentVersion)
			if err != nil && tC.err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected Run(ctx, %v) => %v, got %v", tC.currentVersion, tC.err, err)
			}

			if testWriter.String() != tC.want {
				t.Errorf("Run(ctx, %v) wrong output: expected %v, got %v", tC.currentVersion, testWriter.String(), tC.want)
			}
		})
	}
}
