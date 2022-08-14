package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
)

// CreateTempGoBinaryArchive is test helper function used to create a fake golang binary archive.
func CreateTempGoBinaryArchive(t *testing.T, archiveVersion string) (tmpArchive *os.File, tmpDir string) {
	t.Helper()
	tmpDir = t.TempDir()

	archiveName := fmt.Sprintf("%s%s.%s", downloader.Prefix(), archiveVersion, downloader.Postfix())
	downloadPath := filepath.Join(tmpDir, archiveName)

	tmpArchive, err := os.Create(downloadPath)
	if err != nil {
		t.Fatalf("CreateTempGoBinaryArchive: failed to create temp binary archive: %v", err)
	}

	gzWriter := gzip.NewWriter(tmpArchive)
	_, err = gzWriter.Write([]byte("This is test data"))
	if err != nil {
		t.Fatalf("CreateTempGoBinaryArchive: failed to write content to temp binary archive: %v", err)
	}

	return tmpArchive, tmpDir
}

// ExecuteCommand is a test helper that executes the specified `godl` sub command
func ExecuteCommand(t *testing.T, ignoreCmdError bool, root *cobra.Command, args ...string) (output, errOutput string) {
	t.Helper()
	_, output, errOutput, err := executeCommandC(root, args)
	if err != nil && !ignoreCmdError {
		t.Errorf("godl %s failed: %v", args[0], err)
	}
	return output, errOutput
}

func executeCommandC(root *cobra.Command, args []string) (c *cobra.Command, output, errOutput string, err error) {
	outputBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	root.SetOut(outputBuf)
	root.SetErr(errBuf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, outputBuf.String(), errBuf.String(), err
}

// The RoundTripFunc type is an adapter to allow the use of
// ordinary functions as  net/http.RoundTripper. If f is a function
// with the appropriate signature, RoundTripFunc(f) is a
// RoundTripper that calls f.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with a Fake Transport
func NewTestClient(fn http.RoundTripper) *http.Client {
	return &http.Client{Transport: fn}
}
