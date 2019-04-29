package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"unicode/utf8"
)

type inMemoryWriter struct {
	*bytes.Buffer
}

func (iM *inMemoryWriter) Write(p []byte) (int, error) {
	iM.Write(p)
	return iM.Len(), nil
}

func (iM *inMemoryWriter) Close() error { return nil }

type inMemoryFileCreator struct{}

func (inMem *inMemoryFileCreator) Create(name string) (io.WriteCloser, error) {
	w := &inMemoryWriter{new(bytes.Buffer)}
	return w, nil
}

func (inMem *inMemoryFileCreator) Rename(oldPath, newPath string) error {
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
	ic := &inMemoryFileCreator{}
	testClient := NewTestClient(func(req *http.Request) *http.Response {
		testData := bytes.NewBufferString("This is test data")
		contentLength := fmt.Sprintf("%v", utf8.RuneCount(testData.Bytes()))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(testData),
			Header: map[string][]string{
				"Content-Length": []string{contentLength},
			},
		}
	})
	testGoBinDownloader := &goBinaryDownloader{
		BaseURL: "",
		Client:  testClient,
		FC:      ic,
	}

	err := testGoBinDownloader.download("1.12", ".")
	if err != nil {
		t.Errorf("Error downloading go binary: %v", err)
	}
}
