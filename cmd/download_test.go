package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

// inMemoryWriter is an in-memory writer that writes to
type inMemoryWriter struct {
	*bytes.Buffer
}

func (iM *inMemoryWriter) Write(p []byte) (int, error) {
	iM.Write(p)
	return iM.Len(), nil
}

func (iM *inMemoryWriter) Close() error { return nil }

type inMemoryFileCreatorRenamer struct {
	writer bytes.Buffer
}

func (inMemFCR *inMemoryFileCreatorRenamer) Create(name string) (io.WriteCloser, error) {
	w := &inMemoryWriter{&inMemFCR.writer}
	return w, nil
}

func (inMemFCR *inMemoryFileCreatorRenamer) Rename(oldPath, newPath string) error {
	return nil
}

// RoundTripFunc
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with a Fake Transport
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDownloadGoBinary(t *testing.T) {
	ic := &inMemoryFileCreatorRenamer{}
	testClient := NewTestClient(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")
		contentLength := fmt.Sprintf("%v", len(testData.Bytes()))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(testData),
			Header: map[string][]string{
				"Content-Length": []string{contentLength},
			},
		}
	})
	testGoBinDownloader := &goBinaryDownloader{
		BaseURL: "https://dl.google.com/go/",
		Client:  testClient,
		fCR:     ic,
	}

	err := testGoBinDownloader.download("1.12", ".", true)
	if err != nil {
		t.Errorf("Error downloading go binary: %v", err)
	}

	if ic.writer.String() != "This is test data" {
		t.Errorf("Data downloaded does not match data written to archive")
	}

	ic.writer.Reset()
}

func TestDownloadCmdCalledWithNoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "download")
	expected := "provide binary archive version to download"
	got := err.Error()
	if got != expected {
		t.Errorf("godl download Unknown error: %v", err)
	}
}

func TestDownloadCmdHelp(t *testing.T) {
	_, err := executeCommand(rootCmd, "download", "-h")
	if err != nil {
		t.Errorf("godl download failed: %v", err)
	}
}
