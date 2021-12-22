package download

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

func fakeHashVerifier(input io.Reader, hex string) error {
	return nil
}

func TestDownloadRelease(t *testing.T) {
	fakeRoundTripper := func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	}
	testClient := test.NewTestClient(test.RoundTripFunc(fakeRoundTripper))

	imFS := fsys.NewInMemFS(make(fstest.MapFS))
	dl := &downloader.Downloader{
		BaseURL:      "https://storage.googleapis.com/golang/",
		Client:       testClient,
		DownloadDir:  ".",
		FS:           imFS,
		Hasher:       hash.FakeHasher{},
		HashVerifier: fakeHashVerifier,
	}

	d := Download{dl, 5 * time.Second}
	err := d.Run(context.Background(), "1.12")
	if err != nil {
		t.Fatalf("Error downloading go binary: %v", err)
	}

	entries, err := imFS.ReadDir(".")
	if err != nil {
		t.Error(err)
	}

	expected := "go1.12.darwin-amd64.tar.gz"
	if entries[0].Name() != expected {
		t.Errorf("downloaded filename does not match. want %s; got: %s",
			expected, entries[0].Name())
	}
}
