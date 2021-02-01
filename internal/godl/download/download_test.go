package download

import (
	"bytes"
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
	testClient := test.NewTestClient(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	})

	imFS := inmem.NewFS(new(bytes.Buffer))
	dl := &downloader.Downloader{
		BaseURL:      "https://storage.googleapis.com/golang/",
		Client:       testClient,
		DownloadDir:  ".",
		Fsys:         imFS,
		Hasher:       hash.FakeHasher{},
		HashVerifier: fakeHashVerifier,
	}

	err := downloadRelease("1.12", dl)
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

	_, _, err := test.ExecuteCommand(godlCmd, "download")
	expected := "provide binary archive version to download"
	got := err.Error()
	if got != expected {
		t.Errorf("godl download Unknown error: %v", err)
	}
}

func TestDownloadCmdHelp(t *testing.T) {
	godlCmd := godl.New()
	download := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{download})

	_, _, err := test.ExecuteCommand(godlCmd, "download", "-h")
	if err != nil {
		t.Errorf("godl download failed: %v", err)
	}
}
