package godlutil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path"

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
