package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/MakeNowJust/heredoc"

	"github.com/dikaeinstein/godl/test"
)

const (
	versionWithUpdate    = "0.11.5"
	versionWithoutUpdate = "0.11.6"
)

func TestCheckForUpdate(t *testing.T) {
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

	const errMsg = "Maximum number of login attempts exceeded. Please try again later."
	failingTestClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		errResp := ListReleasesErrorResp{
			Message:          errMsg,
			DocumentationURL: "https://docs.github.com/rest",
		}
		b, err := json.Marshal(errResp)
		if err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       io.NopCloser(bytes.NewReader(b)),
		}
	}))

	testCases := []struct {
		err            error
		client         *http.Client
		name           string
		currentVersion string
		want           bool
	}{
		{
			name:           "Returns true if update is available",
			client:         testClient,
			currentVersion: versionWithUpdate,
			want:           true,
			err:            nil,
		},
		{
			name:           "Returns false if no update",
			client:         testClient,
			currentVersion: versionWithoutUpdate,
			want:           false,
			err:            nil,
		},
		{
			name:           "Returns error encountered while checking for update",
			client:         failingTestClient,
			currentVersion: versionWithUpdate,
			want:           false,
			err:            fmt.Errorf("https://api.github.com/repos/dikaeinstein/godl/releases?per_page=10: 502 %s", errMsg),
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			testWriter := &bytes.Buffer{}

			u := Update{tC.client, testWriter}
			exists, _, err := u.CheckForUpdate(context.Background(), tC.currentVersion)
			if err != nil && tC.err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected CheckForUpdate(ctx, %v) => %v, got %v", tC.currentVersion, tC.err, err)
			}
			if exists != tC.want {
				t.Errorf("expected CheckForUpdate(ctx, %v) => %v, got %v", tC.currentVersion, tC.want, exists)
			}
		})
	}
}

func TestRun(t *testing.T) {
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

	testCases := []struct {
		client         *http.Client
		currentVersion string
		name           string
		err            error
		want           string
	}{
		{
			name:           "Returns no error",
			client:         testClient,
			currentVersion: versionWithoutUpdate,
			err:            nil,
			want:           "No update available.\n",
		},
		{
			name:           "Returns correct message when no update is available",
			client:         testClient,
			currentVersion: versionWithoutUpdate,
			err:            nil,
			want:           "No update available.\n",
		},
		{
			name:           "Returns correct message when update is available",
			client:         testClient,
			currentVersion: versionWithUpdate,
			err:            nil,
			want: heredoc.Doc(`
				Your version of Godl is out of date!

				The latest version is v0.11.6.
				You can update by downloading from https://github.com/dikaeinstein/godl/releases
			`),
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
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
