package app

import (
	"fmt"
	"io"
	"runtime/debug"
)

type BuildInfo struct {
	BuildTime string
	GitHash   string
	GitTag    string
	GoVersion string
}

const (
	BuildTimeUnknown = "unknown"
	GitTagDefault    = "dev"
	goVersionDefault = "unknown"
	gitHashDefault   = "none"
)

// GoVersion returns the Go version used to build the binary.
func GoVersion(info *debug.BuildInfo) string {
	if info == nil {
		return goVersionDefault
	}

	return info.GoVersion
}

// GitHash returns the short git hash of the module.
func GitHash(info *debug.BuildInfo) string {
	if info == nil {
		return gitHashDefault
	}

	h, ok := readBuildInfoSetting(info, "vcs.revision")
	if !ok {
		return gitHashDefault
	}

	return h[:7]
}

// readBuildInfoSetting returns the value of the Go BuildInfo setting
// with the given key.
func readBuildInfoSetting(info *debug.BuildInfo, key string) (val string, ok bool) {
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value, true
		}
	}

	return "", false
}

// Version represents the version information of the binary.
type Version struct {
	buildInfo BuildInfo
}

func NewVersion(info BuildInfo) Version {
	debugBuildInfo, _ := debug.ReadBuildInfo()

	if info.GitHash == "" {
		info.GitHash = GitHash(debugBuildInfo)
	}
	if info.GoVersion == "" {
		info.GoVersion = GoVersion(debugBuildInfo)
	}
	if info.BuildTime == "" {
		info.BuildTime = BuildTimeUnknown
	}
	if info.GitTag == "" {
		info.GitTag = GitTagDefault
	}

	return Version{buildInfo: info}
}

// Run prints the version information to the given writer.
func (v Version) Run(out io.Writer) error {
	_, err := fmt.Fprintf(
		out,
		"Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s\n",
		v.buildInfo.GitTag,
		v.buildInfo.GoVersion,
		v.buildInfo.GitHash,
		v.buildInfo.BuildTime,
	)

	return err
}
