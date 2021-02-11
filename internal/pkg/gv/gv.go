// Package gv provides functions for working with versions
package gv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/hashicorp/go-version"
)

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

func CompareVersions(left, right *version.Version, d SortDirection) bool {
	if d == Asc {
		return left.LessThan(right)
	}

	return left.GreaterThan(right)
}

// GetVersion returns the version from the given string.
// 		go1.11.4.darwin-amd64.tar.gz => 1.11.4
func GetVersion(s string) *version.Version {
	v := strings.Split(s, ".darwin-amd64")
	vv, err := version.NewVersion(strings.TrimPrefix(v[0], "go"))
	godlutil.Must(err)
	return vv
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
