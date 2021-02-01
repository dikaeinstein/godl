package install

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

type testGzUnArchiver struct{}

func (testGzUnArchiver) Unarchive(source, target string) error { return nil }

func fakeHashVerifier(input io.Reader, hex string) error {
	return nil
}

func TestInstallRelease(t *testing.T) {
	testClient := test.NewTestClient(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	})

	failingTestClient := test.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	})

	tests := map[string]struct {
		c                 *http.Client
		downloadedVersion string
		installVersion    string
		success           bool
	}{
		"installRelease downloads from remote when version not found locally": {
			testClient, "1.10.1", "1.11.7", true,
		},
		"installRelease installs local downloaded version": {testClient, "1.10.6", "1.10.6", true},
		"installRelease handle error when fetching binary from remote": {
			failingTestClient, "1.10.1", "1.11.9", false,
		},
	}

	tmpDir := t.TempDir()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmpFile, err := test.CreateTempGoBinaryArchive(tmpDir, tc.downloadedVersion)
			defer tmpFile.Close()

			dl := &downloader.Downloader{
				BaseURL:      "https://storage.googleapis.com/golang/",
				Client:       tc.c,
				DownloadDir:  ".",
				Fsys:         inmem.NewFS(new(bytes.Buffer)),
				Hasher:       hash.FakeHasher{},
				HashVerifier: fakeHashVerifier,
			}
			install := installCmd{
				archiver: testGzUnArchiver{},
				dl:       dl,
			}
			err = install.Run(tc.installVersion)
			var got bool
			if err != nil {
				got = false
			} else {
				got = true
			}

			if got != tc.success {
				t.Errorf("Error installing go binary: %v", err)
			}
		})
	}
}

func TestInstallCmdCalledWithNoArgs(t *testing.T) {
	godlCmd := godl.New()
	install := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{install})

	_, _, err := test.ExecuteCommand(godlCmd, "install")
	expected := "provide binary archive version to install"
	got := err.Error()
	if got != expected {
		t.Errorf("godl install Unknown error: %v", err)
	}
}

func TestInstallCommandHelp(t *testing.T) {
	godlCmd := godl.New()
	install := New()
	godlCmd.RegisterSubCommands([]*cobra.Command{install})

	_, _, err := test.ExecuteCommand(godlCmd, "install", "-h")
	if err != nil {
		t.Errorf("godl install failed: %v", err)
	}
}
