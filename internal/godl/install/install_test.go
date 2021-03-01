package install

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/pkg/fs/inmem"
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
		success           bool
		pathsD            string
	}{
		"installRelease downloads from remote when version not found locally": {
			testClient, "1.10.1", "1.11.7", true, "/usr/local/go/bin\n",
		},
		"installRelease installs local downloaded version": {
			testClient, "1.10.6", "1.10.6", true, "/usr/local/go/bin\n",
		},
		"installRelease handle error when fetching binary from remote": {
			failingTestClient, "1.10.1", "1.11.9", false, "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmpFile, _ := test.CreateTempGoBinaryArchive(t, tc.downloadedVersion)
			defer tmpFile.Close()

			storage := new(bytes.Buffer)
			dl := &downloader.Downloader{
				BaseURL:      "https://storage.googleapis.com/golang/",
				Client:       tc.c,
				DownloadDir:  ".",
				Fsys:         inmem.NewFS(storage),
				Hasher:       hash.FakeHasher{},
				HashVerifier: fakeHashVerifier,
			}
			install := Install{
				Archiver: testGzUnArchiver{},
				Dl:       dl,
				Timeout:  5 * time.Second,
			}
			err := install.Run(context.Background(), tc.installVersion)
			var got bool
			if err != nil {
				got = false
			} else {
				got = true
			}

			if storage.String() != tc.pathsD {
				t.Errorf("Error adding to $PATH")
			}

			if got != tc.success {
				t.Errorf("Error installing go binary: %v", err)
			}
		})
	}
}
