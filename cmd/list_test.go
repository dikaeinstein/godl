package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func createTempGodlDownloadDir() (string, error) {
	tmpDir, err := ioutil.TempDir(".", "_godlDownloadDir")
	if err != nil {
		return "", fmt.Errorf("create godlDownloadDir tempdir: %s", err)
	}
	return tmpDir, nil
}

func createTempGoBinaryArchive(tmpDir, archiveVersion string) (*os.File, error) {
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

func TestListDownloadedBinaryArchives(t *testing.T) {
	tmpDir, err := createTempGodlDownloadDir()
	if err != nil {
		t.Errorf("ListDownloadedBinaryArchives failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile, err := createTempGoBinaryArchive(tmpDir, "1.13")
	defer tmpFile.Close()

	if listDownloadedBinaryArchives(tmpDir); err != nil {
		t.Errorf("Error listing downloaded archive versions: %v", err)
	}
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestListCommand(t *testing.T) {
	_, err := executeCommand(rootCmd, "list")
	if err != nil {
		t.Errorf("godl list failed: %v", err)
	}
}
