package app

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"github.com/dikaeinstein/downloader/pkg/hash"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/internal/downloader"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/test"
)

type testGzUnArchiver struct{}

func (testGzUnArchiver) Unarchive(ctx context.Context, source, target string) error {
	return nil
}

type fakeHashVerifier struct{}

func (fakeHashVerifier) Verify(_ io.Reader, _ string) error {
	return nil
}

func TestInstallRelease(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	}))

	failingTestClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	}))

	testCases := []struct {
		name              string
		c                 *http.Client
		downloadedVersion string
		installVersion    string
		errMsg            string
		goPathsD          string
	}{
		{
			"RemoteVersion",
			testClient, "1.10.1", "1.11.7", "", "/usr/local/go/bin\n",
		},
		{
			"LocalVersion",
			testClient, "1.10.6", "1.10.6", "", "/usr/local/go/bin\n",
		},
		{
			"HandleErrorWhenFetchingBinaryFromRemote",
			failingTestClient, "1.10.1", "1.11.9",
			"error downloading 1.11.9: no binary release of https://dl.google.com/go/go1.11.9.darwin-arm64.tar.gz", "",
		},
	}

	for i := range testCases {
		tC := testCases[i]

		t.Run(tC.name, func(t *testing.T) {
			tmpFile, _ := test.CreateTempGoBinaryArchive(t, tC.downloadedVersion, "darwin", "arm64")
			t.Cleanup(func() {
				tmpFile.Close()
			})

			imFS := fsys.NewInMemFS(fstest.MapFS{})
			dl, err := downloader.New(
				imFS,
				hash.FakeHasher{},
				fakeHashVerifier{}, tC.c,
				"https://dl.google.com/go/",
				".",
				false,
			)
			require.NoError(t, err)

			install := Install{
				Archiver: testGzUnArchiver{},
				Dl:       dl,
				Timeout:  5 * time.Second,
				FS:       imFS,
			}
			err = install.Run(context.Background(), tC.installVersion, "darwin", "arm64", false)
			if err != nil && err.Error() != tC.errMsg {
				t.Error(err)
			}

			f, ok := imFS.MapFS["/etc/paths.d/go"]
			if ok && string(f.Data) != tC.goPathsD {
				t.Errorf("not matching want: %s, got: %s", tC.goPathsD, string(f.Data))
			}
		})
	}
}
