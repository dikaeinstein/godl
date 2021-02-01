package godlutil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetDownloadDir returns the `godl` download directory
func GetDownloadDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("home directory cannot be detected: %v", err)
	}
	return path.Join(home, ".godl", "downloads"), nil
}

// Must enforces that error is non-nil else it panics
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// VersionExists checks if given archive version exists in the specified download directory.
func VersionExists(archiveVersion, downloadDir string) (bool, error) {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := filepath.Join(filepath.Join(downloadDir, archiveName))

	if _, err := os.Stat(downloadPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err // Assuming that other forms of errors are not due to non-existence
	}

	return true, nil
}

// VerifyHash reports whether the named file has contents with
// SHA-256 of the given hex value.
func VerifyHash(input io.Reader, hex string) error {
	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		return err
	}
	if hex != fmt.Sprintf("%x", hash.Sum(nil)) {
		return fmt.Errorf("%s corrupt? does not have expected SHA-256 of %v", input, hex)
	}

	return nil
}
