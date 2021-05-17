// Package gv provides functions for working with versions
package gv

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
)

type SortDirection string

const (
	// Asc sorts in the ascending direction
	Asc SortDirection = "asc"
	// Desc sorts in the descending direction
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
	downloadPath := filepath.Join(filepath.Join(downloadDir, archiveName))

	if _, err := os.Stat(downloadPath); err != nil {
		if os.IsNotExist(err) {
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
