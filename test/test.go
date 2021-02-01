package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/spf13/cobra"
)

func CreateTempGoBinaryArchive(tmpDir, archiveVersion string) (*os.File, error) {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := filepath.Join(tmpDir, archiveName)

	file, err := os.Create(downloadPath)
	if err != nil {
		return nil, err
	}

	archive := gzip.NewWriter(file)
	_, err = archive.Write([]byte("This is test data"))

	return file, err
}

func ExecuteCommand(root *godl.GodlCmd, args ...string) (output, errOutput string, err error) {
	_, output, errOutput, err = executeCommandC(root, args)
	return output, errOutput, err
}

func executeCommandC(root *godl.GodlCmd, args []string) (c *cobra.Command, output, errOutput string, err error) {
	outputBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	// root.SetOut(outputBuf)
	// root.SetErr(errBuf)
	root.SetOutput(outputBuf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, outputBuf.String(), errBuf.String(), err
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
