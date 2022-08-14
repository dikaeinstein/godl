// Package gv provides functions for working with versions
package gv

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	version "github.com/hashicorp/go-version"
)

const (
	// Asc sorts in the ascending direction
	SortAsc = "asc"
	// Desc sorts in the descending direction
	SortDesc = "desc"
)

func CompareVersions(left, right *version.Version, d string) bool {
	if d == SortAsc {
		return left.LessThan(right)
	}

	return left.GreaterThan(right)
}

// GetVersion returns the version from the given string.
// 		go1.11.4.darwin-amd64.tar.gz => 1.11.4
func GetVersion(s string) *version.Version {
	v := strings.Split(s, ".darwin-amd64")
	vv := version.Must(version.NewVersion(strings.TrimPrefix(v[0], "go")))
	return vv
}

// VersionExists checks if given archive version exists in the specified download directory.
func VersionExists(archiveVersion, downloadDir string) (bool, error) {
	const (
		archivePostfix = "darwin-amd64.tar.gz"
		archivePrefix  = "go"
	)

	archiveName := fmt.Sprintf("%s%s.%s", archivePrefix, archiveVersion, archivePostfix)
	downloadPath := filepath.Join(downloadDir, archiveName)

	if _, err := os.Stat(downloadPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return true, err // Assuming that other forms of errors are not due to non-existence
	}

	return true, nil
}

// Segments returns a new version which is just the numeric segments of v.
func Segments(v *version.Version) *version.Version {
	segStr := intSliceToString(v.Segments())
	return version.Must(version.NewSemver(segStr))
}

func intSliceToString(segments []int) string {
	b := make([]string, len(segments))
	for i, v := range segments {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, ".")
}
