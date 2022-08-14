// Package version provides functions for working with versions
package version

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	go_version "github.com/hashicorp/go-version"
)

const (
	// Asc sorts in the ascending direction
	SortAsc = "asc"
	// Desc sorts in the descending direction
	SortDesc = "desc"
)

func CompareVersions(left, right *go_version.Version, d string) bool {
	if d == SortAsc {
		return left.LessThan(right)
	}

	return left.GreaterThan(right)
}

// GetVersion returns the version from the given string.
// go1.11.4.darwin-amd64.tar.gz => 1.11.4
func GetVersion(s string) *go_version.Version {
	v := strings.Split(s, ".darwin-amd64")
	vv := go_version.Must(go_version.NewVersion(strings.TrimPrefix(v[0], "go")))
	return vv
}

// Exists checks if given archive version exists in the specified download directory.
func Exists(archiveVersion, downloadDir string) (bool, error) {
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
func Segments(v *go_version.Version) *go_version.Version {
	segStr := intSliceToString(v.Segments())
	return go_version.Must(go_version.NewSemver(segStr))
}

func intSliceToString(segments []int) string {
	b := make([]string, len(segments))
	for i, v := range segments {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, ".")
}
