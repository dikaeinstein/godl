package install

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/pkg/fsys"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/dikaeinstein/godl/test"
)

type testGzUnArchiver struct{}

func (testGzUnArchiver) Unarchive(source, target string) error { return nil }

func fakeHashVerifier(input io.Reader, hex string) error {
	return nil
}

func TestInstallRelease(t *testing.T) {
	testClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	}))

	failingTestClient := test.NewTestClient(test.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}))

	tests := map[string]struct {
		c                 *http.Client
		downloadedVersion string
		installVersion    string
		errMsg            string
		goPathsD          string
	}{
		"installRelease downloads from remote when version not found locally": {
			testClient, "1.10.1", "1.11.7", "", "/usr/local/go/bin\n",
		},
		"installRelease installs local downloaded version": {
			testClient, "1.10.6", "1.10.6", "", "/usr/local/go/bin\n",
		},
		"installRelease handle error when fetching binary from remote": {
			failingTestClient, "1.10.1", "1.11.9",
			"error downloading 1.11.9: no binary release of 1.11.9", "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmpFile, _ := test.CreateTempGoBinaryArchive(t, tc.downloadedVersion)
			defer tmpFile.Close()

			imFS := fsys.NewInMemFS(fstest.MapFS{})
			dl := &downloader.Downloader{
				BaseURL:      "https://storage.googleapis.com/golang/",
				Client:       tc.c,
				DownloadDir:  ".",
				FS:           imFS,
				Hasher:       hash.FakeHasher{},
				HashVerifier: fakeHashVerifier,
			}
			install := Install{
				Archiver: testGzUnArchiver{},
				Dl:       dl,
				Timeout:  5 * time.Second,
			}
			err := install.Run(context.Background(), tc.installVersion)
			if err != nil && err.Error() != tc.errMsg {
				t.Error(err)
			}

			f, ok := imFS.MapFS["/etc/paths.d/go"]
			if ok && string(f.Data) != tc.goPathsD {
				t.Errorf("not matching want: %s, got: %s", tc.goPathsD, string(f.Data))
			}
		})
	}
}
