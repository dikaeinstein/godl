package cmd

import (
	"bytes"
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

func (iM *inMemoryWriter) Name() string { return "in-mem" }

type inMemoryFileCreatorRenamer struct {
	writer bytes.Buffer
}

func (inMemFCR *inMemoryFileCreatorRenamer) Create(name string) (writeCloseNamer, error) {
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

// getBinaryHash downloads the given URL and returns it as a string.
func genTestHash(url string) (string, error) {
	return "kdaksjkldfjkdjlks34u3u4iqkj", nil
}

func fakeVerifyHash(file, hex string) error {
	return nil
}

func TestDownloadGoBinary(t *testing.T) {
	ic := &inMemoryFileCreatorRenamer{}
	testClient := NewTestClient(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")

		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(testData),
			ContentLength: int64(len(testData.Bytes())),
		}
	})
	testGoBinDownloader := &goBinaryDownloader{
		BaseURL:    "https://dl.google.com/go/",
		Client:     testClient,
		fCR:        ic,
		genHash:    genTestHash,
		verifyHash: fakeVerifyHash,
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
