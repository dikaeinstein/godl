package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/test"
)

func TestVersionCmd(t *testing.T) {
	info := app.BuildInfo{
		BuildTime: "2021-03-14 00:28",
		GitHash:   "02cb593",
		GitTag:    "v0.11.6",
		GoVersion: "go1.16.2",
	}

	godl := newRootCmd()
	version := newVersionCmd(info)
	registerSubCommands(godl, []*cobra.Command{version})

	expectedOutput := fmt.Sprintf(`Version: %s
Go version: %s
Git hash: %s
Built: %s
`, info.GitTag, info.GoVersion, info.GitHash, info.BuildTime)

	output, errOutput := test.ExecuteCommand(t, false, godl, "version")
	require.Equal(t, "", errOutput)
	require.Equal(t, expectedOutput, output)
}
