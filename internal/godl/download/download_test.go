package download

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/internal/pkg/downloader"
	"github.com/dikaeinstein/godl/pkg/fs/inmem"
	"github.com/dikaeinstein/godl/pkg/hash"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
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

	imFS := inmem.NewFS(new(bytes.Buffer))
	dl := &downloader.Downloader{
		BaseURL:      "https://storage.googleapis.com/golang/",
		Client:       testClient,
		DownloadDir:  ".",
		Fsys:         imFS,
		Hasher:       hash.FakeHasher{},
		HashVerifier: fakeHashVerifier,
	}

	d := downloadCmd{dl}
	err := d.Run(context.Background(), "1.12")
	if err != nil {
		t.Fatalf("Error downloading go binary: %v", err)
	}

	if imFS.Content().String() != "This is test data" {
		t.Errorf("Data downloaded does not match data written to archive")
	}

	imFS.Content().Reset()
}

func TestDownloadCmdCalledWithNoArgs(t *testing.T) {
	godlCmd := godl.New()
	download := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{download})

	_, errOutput := test.ExecuteCommand(t, true, godlCmd, "download")
	expected := "Error: provide binary archive version to download\n"
	if errOutput != expected {
		t.Errorf("godl download failed: expected: %s; got: %s", expected, errOutput)
	}
}

func TestDownloadCmdHelp(t *testing.T) {
	godlCmd := godl.New()
	download := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{download})

	_, errOutput := test.ExecuteCommand(t, true, godlCmd, "download", "-h")
	if errOutput != "" {
		t.Errorf("godl download failed: %v", errOutput)
	}
}
